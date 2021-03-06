package main

import (
	"flag"
	"fmt"
	"github.com/kballard/dcpu16/dcpu"
	"github.com/kballard/dcpu16/dcpu/core"
	"github.com/kballard/termbox-go"
	"io/ioutil"
	"os"
)

var requestedRate dcpu.ClockRate = dcpu.DefaultClockRate
var printRate *bool = flag.Bool("printRate", false, "Print the effective clock rate at termination")

func main() {
	// command-line flags
	flag.Var(&requestedRate, "rate", "Clock rate to run the machine at")
	// update usage
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s [flags] program\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()
	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(2)
	}
	program := flag.Arg(0)
	data, err := ioutil.ReadFile(program)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	// Interpret the file as Words
	words := make([]core.Word, len(data)/2)
	for i := 0; i < len(data)/2; i++ {
		w := core.Word(data[i*2])<<8 + core.Word(data[i*2+1])
		words[i] = w
	}

	// Set up a machine
	machine := new(dcpu.Machine)
	if err := machine.State.LoadProgram(words, 0); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if err := machine.Start(requestedRate); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	var effectiveRate dcpu.ClockRate
	// now wait for the q key
	for {
		evt := termbox.PollEvent()
		if err := machine.HasError(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		if evt.Type == termbox.EventKey {
			if evt.Key == termbox.KeyCtrlC || (evt.Mod == 0 && evt.Ch == 'q') {
				effectiveRate = machine.EffectiveClockRate()
				if err := machine.Stop(); err != nil {
					fmt.Fprintln(os.Stderr, err)
					os.Exit(1)
				}
				break
			}
		}
	}
	if *printRate {
		fmt.Printf("Effective clock rate: %s\n", effectiveRate)
	}
}
