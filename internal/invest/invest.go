// Package invest implements all business logic regarding interest and related types
package invest

import (
	"fmt"
	"math"

	"github.com/pkg/errors"
)

const (

	// Sa for aving type
	Sa = "Saving"
	// CD for cd type
	CD = "CD"
	// Ch gor checking type
	Ch = "Checking"
	// Br for Brokered type
	Br = "Brokered CD"
)

// NewInterestBanks is for input/request
type NewInterestBanks struct {
	NewBanks []NewBank `json:"banks"`
}

// InterestBanks is for outout/response
type InterestBanks struct {
	Banks []Bank `json:"banks"`

	Delta float64 `json:"30daysDelta"`
}

// NewBank is for input/request with Bank data and its deposits
type NewBank struct {
	Name        string       `json:"name"`
	NewDeposits []NewDeposit `json:"deposits"`
}

// Bank is for output/response with Bank data and its deposists
type Bank struct {
	Name     string    `json:"name"`
	Deposits []Deposit `json:"deposits"`

	Delta float64 `json:"30daysDelta"`
}

// NewDeposit is is for input/request with Bank data and its deposits
type NewDeposit struct {
	Account     string `json:"account"`
	AccountType string `json:"annualType"`

	APY    float64 `json:"annualRate%"`
	Years  float64 `json:"years"`
	Amount float64 `json:"amount"`
}

// Deposit is for output/reponse with Deposit data
type Deposit struct {
	Account     string `json:"account"`
	AccountType string `json:"annualType"`

	APY    float64 `json:"annualRate%"`
	Years  float64 `json:"years"`
	Amount float64 `json:"amount"`

	Delta float64 `json:"30daysDelta"`
}

// CalDelta calcuates deltacli - interest for 30 days for output/response Deposit
func (d Deposit) CalDelta() (float64, error) {
	e := earned(d)
	e30Days, err := earned30days(e, d.Years)
	if err != nil {
		return 0, errors.Wrapf(err, "deltacli calculation for Account: %s", d.Account)
	}
	return e30Days, nil
}

func earned30days(iEarned float64, years float64) (float64, error) {
	if years*365 < 30 {
		return 0, fmt.Errorf("NewDeposit period in years %v should not be less than 30 days", years)
	}
	i1Day := iEarned / (years * 365)
	i30 := i1Day * 30
	return math.Round(i30*100) / 100, nil
}
