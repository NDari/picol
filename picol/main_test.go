package main

import (
	"io/ioutil"
	"testing"
)

func BenchmarkLoopAndSquare(b *testing.B) {

	source, err := ioutil.ReadFile("../examples/t2.tcl")
	if err != nil {
		b.Fail()
	}

	interp := LoadedInterp()
	str_source := string(source)
	for i := 0; i < b.N; i++ {
		interp.Eval(str_source)
	}
}

func BenchmarkFib(b *testing.B) {
	source, err := ioutil.ReadFile("../examples/fib.tcl")
	if err != nil {
		b.Fail()
	}

	interp := LoadedInterp()
	str_source := string(source)

	for i := 0; i < b.N; i++ {
		interp.Eval(str_source)
	}
}
