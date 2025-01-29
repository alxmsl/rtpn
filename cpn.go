package cpn

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/alxmsl/prmtvs/plexus"
	"github.com/alxmsl/prmtvs/skm"
	"golang.org/x/sync/errgroup"
)

// ProcessFunc defines a placement processor. Token comes from input channel.
// Once Token is processed it should be sent into the output channel.
type ProcessFunc func(ctx context.Context, input <-chan Token, output chan<- Token) error

const (
	// StateInactive defines that Net has not been run yet. The Net topology can be changed in the StateInactive.
	StateInactive = iota
	// StateActive defines that Net has been run already. The Net topology can not be changed.
	StateActive
)

// Net structure defines a Petri-Net as two sets: set of placements and set of transitions.
type Net struct {
	sync.RWMutex

	ctx    context.Context
	cancel context.CancelFunc
	id     string
	state  int

	placements  map[string]*placement
	transitions map[string]*transition
}

// NewNet creates a Net with a given identifier.
func NewNet(id string) *Net {
	return &Net{
		id: id,

		placements:  make(map[string]*placement),
		transitions: make(map[string]*transition),
	}
}

// AddPlace adds unique placement to the Net with a given processor.
// Net should not be in StateActive to add placement.
func (n *Net) AddPlace(id string, processFn ProcessFunc) error {
	n.Lock()
	defer n.Unlock()
	if n.state == StateActive {
		return fmt.Errorf("can't add placement %s: %w", id, ErrorNetIsActive)
	}

	if _, ok := n.placements[id]; ok {
		return fmt.Errorf("can't add placement %s: %w", id, ErrorEntityAlreadyExists)
	}
	n.placements[id] = newPlacement(id, processFn)
	return nil
}

// AddTransition creates a unique transition between given placements.
// Net should not be in StateActive to add transition.
func (n *Net) AddTransition(id string, fromIDs []string, toIDs []string) error {
	n.Lock()
	defer n.Unlock()
	if n.state == StateActive {
		return fmt.Errorf("can't add transition %s: %w", id, ErrorNetIsActive)
	}

	if _, ok := n.transitions[id]; ok {
		return fmt.Errorf("can't add transition %s: %w", id, ErrorEntityAlreadyExists)
	}
	for _, placeID := range fromIDs {
		if _, ok := n.placements[placeID]; !ok {
			return fmt.Errorf("can't add transition from placement %s: %w", placeID, ErrorEntityNotFound)
		}
	}
	for _, placeID := range toIDs {
		if _, ok := n.placements[placeID]; !ok {
			return fmt.Errorf("can't add transition to placement %s: %w", placeID, ErrorEntityNotFound)
		}
	}

	var t = newTransition(id, fromIDs, toIDs)
	n.transitions[id] = t

	for _, placeID := range fromIDs {
		n.placements[placeID].addOuts(t)
	}
	for _, placeID := range toIDs {
		n.placements[placeID].addIns(t)
	}
	return nil
}

// Markup sends a given Token into a given placement.
func (n *Net) Markup(placeID string, token Token) error {
	if _, ok := n.placements[placeID]; !ok {
		return fmt.Errorf("can't put token in placement %s: %w", placeID, ErrorEntityNotFound)
	}
	return n.placements[placeID].putToken(token)
}

// Run runs the Net to start Token processing.
func (n *Net) Run(ctx context.Context) error {
	n.Lock()
	if n.state == StateActive {
		defer n.Unlock()
		return fmt.Errorf("can't run %s: %w", n.id, ErrorNetIsActive)
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

// Stats returns sorted-key map with placements performance:
//   - Number of accepted Token by placement.
//   - Number of processing Token by placement.
//   - Number of sent Token by placement.
func (n *Net) Stats() *skm.SKM {
	var stats = skm.NewSortedKeyMap()
	for _, p := range n.placements {
		stats.Add(p.id, p.statsValue())
	}
	return stats
}

// Stop stops running of the Net.
func (n *Net) Stop() error {
	n.Lock()
	defer n.Unlock()
	if n.state == StateInactive {
		return fmt.Errorf("can't stop %s: %w", n.id, ErrorNetIsInactive)
	}
	n.cancel()
	n.state = StateInactive
	return nil
}

// Stats struct defines a placement performance data.
type Stats struct {
	Accepted   uint64 // Number of accepted Token.
	Processing uint64 // Number of processing Token.
	Sent       uint64 // Number of sent Token.
}

func (s Stats) String() string {
	return fmt.Sprintf(`{"Accepted": %d, "Processing": %d, "Sent": %d}`, s.Accepted, s.Processing, s.Sent)
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
