package handlers

import (
	"io"
	"log"

	"github.com/pkg/errors"
	"github.com/rsachdeva/illuminatingdeposits-rest/interestcal/interestvalue"
	"github.com/rsachdeva/illuminatingdeposits-rest/responder"
)

//Interest handler
type Interest struct {
	Log *log.Logger
}

// Create investment calculates for all banks, sent to the desired writer in JSON format
func (ih Interest) Create(w io.Writer, nibs interestvalue.NewInterest, executionTimes int) error {
	var ibs interestvalue.Interest
	var err error
	for j := 0; j < executionTimes; j++ {
		ibs, err = nibs.ComputeDelta()
	}
	if err != nil {
		return errors.Wrap(err, "create calculating for interestcal.NewInterest")
	}
	return responder.Output(w, ibs)
}
