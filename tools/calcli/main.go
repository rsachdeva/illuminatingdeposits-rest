package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/pprof"

	"github.com/pkg/errors"
	"github.com/rsachdeva/illuminatingdeposits-rest/interestcal/interestvalue"
	"github.com/rsachdeva/illuminatingdeposits-rest/jsonfmt"
	"github.com/rsachdeva/illuminatingdeposits-rest/tools/calcli/handlers"
)

func main() {
	if err := createInterest(); err != nil {
		log.Printf("error: quitting appserver: %+v", err)
		os.Exit(1)
	}

}

func createInterest() error {
	var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")
	var memprofile = flag.String("memprofile", "", "write memory profile to `file`")
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("could not createInterest CPU profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	log := log.New(os.Stdout, "DEPOSITS : ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

	hi := handlers.Interest{Log: log}
	var ni interestvalue.NewInterest

	fmt.Println("flag.Arg(1) is", flag.Arg(1))
	if err := jsonfmt.InputFile(flag.Arg(1), &ni); err != nil {
		return errors.Wrap(err, "parsing json file for interest")
	}
	executionTimes := 1
	if *memprofile != "" {
		executionTimes = 100000
	}
	fmt.Println("executionTimes is", executionTimes)
	if err := hi.Create(os.Stdout, ni, executionTimes); err != nil {
		return errors.Wrap(err, "printing all banks calculations")
	}

	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal("could not createInterest memory profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		runtime.GC()    // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
	}
	return nil
}
