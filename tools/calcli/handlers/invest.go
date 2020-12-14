package handlers

import (
	"io"
	"log"

	"github.com/pkg/errors"
	"github.com/rsachdeva/illuminatingdeposits-rest/interestcal/interestvalue"
	"github.com/rsachdeva/illuminatingdeposits-rest/jsonfmt"
)

//Interest handler
type Interest struct {
	Log *log.Logger
}

// ListCalculations investment calculates for all banks, sent to the desired writer in JSON format
func (ih Interest) ListCalculations(w io.Writer, nibs interestvalue.InterestRequest, executionTimes int) error {
	var ibs interestvalue.InterestResponse
	var err error
	for j := 0; j < executionTimes; j++ {
		ibs, err = nibs.CalculateDelta()
	}
	if err != nil {
		return errors.Wrap(err, "create calculating for interestcal.InterestRequest")
	}
	return jsonfmt.Output(w, ibs)
}
