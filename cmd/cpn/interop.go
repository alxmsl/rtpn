package main

import (
	"context"
	"fmt"
	"plugin"
	"reflect"

	"github.com/alxmsl/cpn"
)

const defaultLookupName = "Net"

type PetriNet interface {
	Markup(placeID string, token cpn.Token) error
	Run(ctx context.Context) error
	Stop() error
}

func lookup(filename string) (PetriNet, error) {
	var pluginCPN, err = plugin.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("open plugin: %w", err)
	}
	sym, err := pluginCPN.Lookup(defaultLookupName)
	if err != nil {
		return nil, fmt.Errorf("lookup %s: %w", defaultLookupName, err)
	}
	net, ok := sym.(PetriNet)
	if !ok {
		return nil, fmt.Errorf("cast failed because %s is %s", defaultLookupName, reflect.TypeOf(sym).String())
	}
	return net, nil
}
