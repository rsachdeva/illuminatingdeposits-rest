// Calculations

package invest

import (
	"math"
)

// Delta calculations for all banks
func Delta(ni NewInterestBanks) (InterestBanks, error) {
	var i InterestBanks
	var ibs []Bank
	var iDelta float64
	for _, nb := range ni.NewBanks {
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
			delta, err := d.CalDelta()
			if err != nil {
				return i, err
			}
			d.Delta = roundToNearest(delta)
			ds = append(ds, d)
			bDelta = bDelta + roundToNearest(delta)
		}
		bk := Bank{
			Name:     nb.Name,
			Deposits: ds,
			Delta:    roundToNearest(bDelta),
		}
		ibs = append(ibs, bk)
		iDelta = iDelta + bk.Delta
	}
	i = InterestBanks{
		Banks: ibs,
		Delta: roundToNearest(iDelta),
	}
	return i, nil
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

func earned(d Deposit) float64 {
	switch d.AccountType {
	case Sa, CD:
		return compoundInterest(d.APY, d.Years, d.Amount)
	case Br:
		return simpleInterest(d.APY, d.Years, d.Amount)
	default:
		return 0.0
	}
}
