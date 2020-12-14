// Calculations

package interestvalue

import (
	"fmt"
	"math"

	"github.com/pkg/errors"
)

// CalculateDelta calculations for all banks
func (cireq CreateInterestRequest) CalculateDelta() (CreateInterestResponse, error) {
	bks, delta, err := cireq.computeBanksDelta()
	if err != nil {
		return CreateInterestResponse{}, err
	}
	ciresp := CreateInterestResponse{
		Banks: bks,
		Delta: roundToNearest(delta),
	}
	return ciresp, nil
}

func (cireq CreateInterestRequest) computeBanksDelta() ([]Bank, float64, error) {
	var bks []Bank
	var delta float64
	for _, nb := range cireq.NewBanks {
		ds, bDelta, err := nb.computeBankDelta()
		if err != nil {
			return nil, 0, err
		}
		bk := Bank{
			Name:     nb.Name,
			Deposits: ds,
			Delta:    roundToNearest(bDelta),
		}
		bks = append(bks, bk)
		delta = delta + bk.Delta
	}
	return bks, delta, nil
}

func (nb NewBank) computeBankDelta() ([]Deposit, float64, error) {
	var ds []Deposit
	var bDelta float64
	for _, nd := range nb.NewDeposits {
		d := Deposit{
			Account:     nd.Account,
			AccountType: nd.AccountType,
			APY:         nd.APY,
			Years:       nd.Years,
			Amount:      nd.Amount,
		}
		err := d.computeDepositDelta()
		if err != nil {
			return nil, 0, err
		}
		ds = append(ds, d)
		bDelta = bDelta + d.Delta
	}
	return ds, bDelta, nil
}

// CalculateDelta calculates interest for 30 days for output/response Deposit
func (d *Deposit) computeDepositDelta() error {
	e := d.earned()
	e30Days, err := earned30days(e, d.Years)
	if err != nil {
		return errors.Wrapf(err, "calculation for Account: %s", d.Account)
	}
	d.Delta = roundToNearest(e30Days)
	return nil
}

func (d *Deposit) earned() float64 {
	switch d.AccountType {
	case Sa, CD:
		return compoundInterest(d.APY, d.Years, d.Amount)
	case Br:
		return simpleInterest(d.APY, d.Years, d.Amount)
	default:
		return 0.0
	}
}

func roundToNearest(n float64) float64 {
	return math.Round(n*100) / 100
}

func simpleInterest(apy float64, years float64, amount float64) float64 {
	rateInDecimal := apy / 100
	intEarned := amount * rateInDecimal * years
	return intEarned
}

func compoundInterest(apy float64, years float64, amount float64) float64 {
	rateInDecimal := apy / 100
	calInProcess := math.Pow(1+rateInDecimal, years)
	intEarned := amount*calInProcess - amount
	return intEarned
}

func earned30days(iEarned float64, years float64) (float64, error) {
	if years*365 < 30 {
		return 0, fmt.Errorf("NewDeposit period in years %v should not be less than 30 days", years)
	}
	i1Day := iEarned / (years * 365)
	i30 := i1Day * 30
	return math.Round(i30*100) / 100, nil
}
