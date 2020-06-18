package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/pprof"

	"github.com/pkg/errors"
	"github.com/rsachdeva/illuminatingdeposits/cmd/deltacli/internal/handlers"
	"github.com/rsachdeva/illuminatingdeposits/internal/invest"
	"github.com/rsachdeva/illuminatingdeposits/internal/platform/inout"
)

func main() {
	if err := run(); err != nil {
		log.Printf("error: quitting app: %+v", err)
		os.Exit(1)
	}

}

func run() error {
	var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")
	var memprofile = flag.String("memprofile", "", "write memory profile to `file`")
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	log := log.New(os.Stdout, "ILLUMINATINGDEPOSITS : ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

	hi := handlers.Interest{Log: log}
	var ni invest.NewInterestBanks

	fmt.Println("flag.Arg(0) is", flag.Arg(0))
	if err := inout.InputJSON(flag.Arg(0), &ni); err != nil {
		return errors.Wrap(err, "parsing json file for interest")
	}
	executionTimes := 1
	if *memprofile != "" {
		executionTimes = 100000
	}
	fmt.Println("executionTimes is", executionTimes)
	if err := hi.BatchGet(os.Stdout, ni, executionTimes); err != nil {
		return errors.Wrap(err, "printing all banks calculations")
	}

	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		runtime.GC()    // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
	}
	return nil
}
