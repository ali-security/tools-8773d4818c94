// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package static

import (
	"bytes"
	"io/ioutil"
	"runtime"
	"strconv"
	"testing"
	"unicode"
)

func TestStaticIsUpToDate(t *testing.T) {
	if runtime.GOOS == "android" {
		t.Skip("files not available on android")
	}
	oldBuf, err := ioutil.ReadFile("static.go")
	if err != nil {
		t.Errorf("error while reading static.go: %v\n", err)
	}

	newBuf, err := Generate()
	if err != nil {
		t.Errorf("error while generating static.go: %v\n", err)
	}

	// gen.go stamps time.Now().UTC().Year() into the license header of the
	// generated file, so the first line (the "// Copyright" line) differs in
	// any year other than the one static.go was generated in. Compare the
	// rest of the file strictly.
	if bytes.Compare(skipFirstLine(oldBuf), skipFirstLine(newBuf)) != 0 {
		t.Error(`static.go is stale.  Run:
  $ go generate golang.org/x/tools/godoc/static
  $ git diff
to see the differences.`)

	}
}

// skipFirstLine returns b without its first line.
func skipFirstLine(b []byte) []byte {
	if i := bytes.IndexByte(b, '\n'); i >= 0 {
		return b[i+1:]
	}
	return nil
}

// TestAppendQuote ensures that AppendQuote produces a valid literal.
func TestAppendQuote(t *testing.T) {
	var in, out bytes.Buffer
	for r := rune(0); r < unicode.MaxRune; r++ {
		in.WriteRune(r)
	}
	appendQuote(&out, in.Bytes())
	in2, err := strconv.Unquote(out.String())
	if err != nil {
		t.Fatalf("AppendQuote produced invalid string literal: %v", err)
	}
	if got, want := in2, in.String(); got != want {
		t.Fatal("AppendQuote modified string") // no point printing got/want: huge
	}
}
