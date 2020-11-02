// Calculations

package invest

import (
	"fmt"
	"math"
)

// Delta calculations for all banks
func Delta(ni NewInterest) (Interest, error) {
	bks, delta, err := banksWithDelta(ni)
	if err != nil {
		return Interest{}, err
	}
	intBanks := Interest{
		Banks: bks,
		Delta: roundToNearest(delta),
	}
	return intBanks, nil
}

func banksWithDelta(ni NewInterest) ([]Bank, float64, error) {
	var bks []Bank
	var delta float64
	for _, nb := range ni.NewBanks {
		ds, bDelta, err := depositsWithDelta(nb)
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

func depositsWithDelta(nb NewBank) ([]Deposit, float64, error) {
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
		err := d.CalDelta()
		if err != nil {
			return nil, 0, err
		}
		ds = append(ds, d)
		bDelta = bDelta + d.Delta
	}
	return ds, bDelta, nil
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
