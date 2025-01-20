package cpn

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/alxmsl/prmtvs/plexus"
	"github.com/alxmsl/prmtvs/skm"
	"golang.org/x/sync/errgroup"
)

type ProcessFunc func(context.Context, <-chan Token, chan<- Token) error

const (
	StateInactive = iota
	StateActive
)

type Net struct {
	sync.RWMutex

	ctx    context.Context
	cancel context.CancelFunc
	id     string
	state  int

	placements  map[string]*placement
	transitions map[string]*transition
}

func NewNet(id string) *Net {
	return &Net{
		id: id,

		placements:  make(map[string]*placement),
		transitions: make(map[string]*transition),
	}
}

func (n *Net) AddPlace(placeID string, processFn ProcessFunc) error {
	n.Lock()
	defer n.Unlock()
	if n.state == StateActive {
		return fmt.Errorf("can't add placement %s: %w", placeID, ErrorNetIsActive)
	}

	if _, ok := n.placements[placeID]; ok {
		return fmt.Errorf("can't add placement %s: %w", placeID, ErrorEntityAlreadyExists)
	}
	n.placements[placeID] = newPlacement(placeID, processFn)
	return nil
}

func (n *Net) AddTransition(transitionID string, fromPlaceIDs []string, toPlaceIDs []string) error {
	n.Lock()
	defer n.Unlock()
	if n.state == StateActive {
		return fmt.Errorf("can't add transition %s: %w", transitionID, ErrorNetIsActive)
	}

	if _, ok := n.transitions[transitionID]; ok {
		return fmt.Errorf("can't add transition %s: %w", transitionID, ErrorEntityAlreadyExists)
	}
	for _, placeID := range fromPlaceIDs {
		if _, ok := n.placements[placeID]; !ok {
			return fmt.Errorf("can't add transition from placement %s: %w", placeID, ErrorEntityNotFound)
		}
	}
	for _, placeID := range toPlaceIDs {
		if _, ok := n.placements[placeID]; !ok {
			return fmt.Errorf("can't add transition to placement %s: %w", placeID, ErrorEntityNotFound)
		}
	}

	var t = newTransition(transitionID, fromPlaceIDs, toPlaceIDs)
	n.transitions[transitionID] = t

	for _, placeID := range fromPlaceIDs {
		n.placements[placeID].addOuts(t)
	}
	for _, placeID := range toPlaceIDs {
		n.placements[placeID].addIns(t)
	}
	return nil
}

func (n *Net) PutToken(name string, token Token) error {
	if _, ok := n.placements[name]; !ok {
		return fmt.Errorf("can't put token in placement %s: %w", name, ErrorEntityNotFound)
	}
	return n.placements[name].putToken(token)
}

func (n *Net) Run(ctx context.Context) error {
	n.Lock()
	if n.state == StateActive {
		defer n.Unlock()
		return fmt.Errorf("can run %s: %w", n.id, ErrorNetIsActive)
	}

	n.ctx, n.cancel = context.WithCancel(ctx)
	defer func() {
		n.Lock()
		n.ctx, n.cancel = nil, nil
		n.Unlock()
	}()
	n.state = StateActive
	n.Unlock()

	var g = errgroup.Group{}
	for _, p := range n.placements {
		var p = p
		g.Go(func() error {
			return p.run(ctx)
		})
	}
	var err = g.Wait()
	return err
}

func (n *Net) Stats() *skm.SKM {
	var stats = skm.NewSortedKeyMap()
	for _, p := range n.placements {
		stats.Add(p.id, p.statsValue())
	}
	return stats
}

func (n *Net) Stop() error {
	n.Lock()
	defer n.Unlock()
	if n.state == StateInactive {
		return fmt.Errorf("can't stop %s: %w", n.id, ErrorNetIsInactive)
	}
	n.cancel()
	return nil
}

type Stats struct {
	Accepted   uint64
	Processing uint64
	Sent       uint64
}

func (s Stats) String() string {
	return fmt.Sprintf(`{"Accepted": %d, "Processing": %d, "Sent": %d}`, s.Accepted, s.Processing, s.Sent)
}

type Token struct {
	payload any
	word    checkpoints
}

func NewToken(payload any) *Token {
	return &Token{
		payload: payload,
	}
}

func (t Token) Payload() any {
	return t.payload
}

func (t Token) String() string {
	return fmt.Sprintf("%v: %s", t.payload, t.word)
}

func (a Token) Merge(b plexus.Mergeable) plexus.Mergeable {
	if _, ok := a.payload.(plexus.Mergeable); !ok {
		panic(fmt.Errorf("payload from a is not mergeable: %w", ErrorWrongTokenType))
	}
	if _, ok := b.(Token).payload.(plexus.Mergeable); !ok {
		panic(fmt.Errorf("payload from b is not mergeable: %w", ErrorWrongTokenType))
	}
	var (
		v1 = a.payload.(plexus.Mergeable)
		v2 = b.(Token).payload.(plexus.Mergeable)
	)
	a.payload = v1.Merge(v2)
	return a
}

type checkpoint struct {
	name string
	when time.Time
}

func (cp checkpoint) String() string {
	return fmt.Sprintf("%s[%s]", cp.name, cp.when.Format(time.RFC3339Nano))
}

type checkpoints []checkpoint

func (cps checkpoints) String() string {
	var result = make([]string, 0, len(cps))
	for _, cp := range cps {
		result = append(result, cp.String())
	}
	return strings.Join(result, ",")
}

func (cps checkpoints) Word() string {
	var ww = make([]string, 0, len(cps))
	for _, cp := range cps {
		ww = append(ww, cp.name)
	}
	return strings.Join(ww, ",")
}

type placement struct {
	id string

	ins  []*transition
	outs []*transition

	processFn    ProcessFunc
	processedCh  chan Token
	processingCh chan Token

	stats Stats
}

func newPlacement(id string, processFn ProcessFunc) *placement {
	return &placement{
		id: id,

		processFn:    processFn,
		processedCh:  make(chan Token),
		processingCh: make(chan Token),
	}
}

func (p *placement) addIns(tt ...*transition) *placement {
	p.ins = append(p.ins, tt...)
	return p
}

func (p *placement) addOuts(tt ...*transition) *placement {
	p.outs = append(p.outs, tt...)
	return p
}

func (p *placement) putToken(token Token) error {
	if p.processingCh == nil {
		return fmt.Errorf("can't put token in placement %s: processing channel is undefined", p.id)
	}
	p.processingCh <- token
	return nil
}

func (p *placement) run(ctx context.Context) error {
	// ins processor
	for _, in := range p.ins {
		go func(in *transition) {
			for {
				v, ok := in.plx.Recv(p.id)
				if !ok {
					panic("in is closed")
				}
				atomic.AddUint64(&p.stats.Accepted, 1)

				t := v.(Token)
				t.word = append(t.word, checkpoint{name: in.id, when: time.Now()})
				p.processingCh <- t
				atomic.AddUint64(&p.stats.Processing, 1)
			}
		}(in)
	}

	// outs processor
	for _, out := range p.outs {
		go func(out *transition) {
			for range out.plx.ReadySend(p.id) {
				out.plx.Send(p.id, <-p.processedCh)
				atomic.AddUint64(&p.stats.Sent, 1)
			}
		}(out)
	}

	return p.processFn(ctx, p.processingCh, p.processedCh)
}

func (p *placement) statsValue() Stats {
	var statsValue = Stats{
		Accepted:   atomic.LoadUint64(&p.stats.Accepted),
		Processing: atomic.LoadUint64(&p.stats.Processing),
		Sent:       atomic.LoadUint64(&p.stats.Sent),
	}
	return statsValue
}

type transition struct {
	id string

	plx *plexus.Plexus
}

func newTransition(id string, senders []string, receivers []string) *transition {
	return &transition{
		plx: plexus.NewPlexus(
			plexus.WithName(id),
			plexus.WithReceivers(receivers...),
			plexus.WithSenders(senders...),
			plexus.WithSelectableSenders(),
		),
		id: id,
	}
}
