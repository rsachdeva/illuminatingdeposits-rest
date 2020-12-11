// Package interestvalue provides struct values and associated operations for interest calculations for
// each deposit in the bank, then for each bank and finally for all banks with 30day average earning
package interestvalue

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

// NewInterest is for input/request
type NewInterest struct {
	NewBanks []NewBank `json:"banks"`
}

// Interest is for output/response
type Interest struct {
	Banks []Bank `json:"banks"`

	Delta float64 `json:"30daysInterest"`
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

	Delta float64 `json:"30daysInterest"`
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

	Delta float64 `json:"30daysInterest"`
}
