// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"context"
	"reflect"
	"strings"
	"testing"
	"time"
)

const (
	ExampleLine   = exampleLine
	ExampleColumn = exampleColumn
	ExampleOffset = exampleOffset
)

const implicitAllInterfaces = "implicitly binds all network interfaces"

// runServe runs the serve command with the given -listen address and returns
// its error. It fails the test if the command does not return promptly, which
// means it is serving on the address.
func runServe(t *testing.T, addr string) error {
	t.Helper()
	s := &Serve{Address: addr, app: New("gopls-test", "", nil, nil)}
	errc := make(chan error, 1)
	go func() { errc <- s.Run(context.Background()) }()
	select {
	case err := <-errc:
		return err
	case <-time.After(10 * time.Second):
		t.Fatalf("serve -listen=%s did not return: it is listening for connections", addr)
		return nil
	}
}

// TestServeListenImplicitAllInterfaces checks that -listen rejects an address
// with no host (e.g. ":12345"), which would implicitly bind all network
// interfaces.
func TestServeListenImplicitAllInterfaces(t *testing.T) {
	for _, addr := range []string{":0", ":12345"} {
		err := runServe(t, addr)
		if err == nil || !strings.Contains(err.Error(), implicitAllInterfaces) {
			t.Errorf("serve -listen=%s: got error %v, want %q", addr, err, implicitAllInterfaces)
		}
	}

	// An explicit host is accepted and passed on to the listener, which
	// fails here only because the port is out of range.
	for _, addr := range []string{"localhost:99999", "0.0.0.0:99999"} {
		err := runServe(t, addr)
		if err == nil {
			t.Errorf("serve -listen=%s: got no error, want invalid port error", addr)
		} else if strings.Contains(err.Error(), implicitAllInterfaces) {
			t.Errorf("serve -listen=%s: explicit host rejected: %v", addr, err)
		}
	}
}

// TestServeNoPortFlag checks that the -port debugging flag, which listened on
// all network interfaces, is no longer defined.
func TestServeNoPortFlag(t *testing.T) {
	typ := reflect.TypeOf(Serve{})
	for i := 0; i < typ.NumField(); i++ {
		if f := typ.Field(i); f.Tag.Get("flag") == "port" {
			t.Errorf("Serve.%s defines the -port flag", f.Name)
		}
	}
}
