package interestcal

import (
	"io"
	"log"

	"github.com/pkg/errors"
	"github.com/rsachdeva/illuminatingdeposits-rest/interestcal/interestvalue"
	"github.com/rsachdeva/illuminatingdeposits-rest/jsonfmt"
)

// Service handler
type Service struct {
	Log *log.Logger
}

// CreateInterest investment calculates for all banks, sent to the desired writer in JSON format
func (ih Service) CreateInterest(w io.Writer, nibs interestvalue.CreateInterestRequest, executionTimes int) error {
	var ibs interestvalue.CreateInterestResponse
	var err error
	for j := 0; j < executionTimes; j++ {
		ibs, err = nibs.CalculateDelta()
	}
	if err != nil {
		return errors.Wrap(err, "create calculating for interestcal.CreateInterestRequest")
	}
	return jsonfmt.Output(w, ibs)
}
