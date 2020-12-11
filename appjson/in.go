package appjson

import (
	"encoding/json"
	"io/ioutil"
)

// InputFile for an inout json file reading and creating go values from it
func InputFile(jsonFileName string, val interface{}) error {
	data, err := ioutil.ReadFile(jsonFileName)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, val)
	if err != nil {
		return err
	}
	return nil
}

// keeping test data for file path testing
// 	var data = `{
//   "banks": [
//     {
//       "name": "HAPPIESTHAPPY",
//       "deposits": [
//         {
//           "account": "1234",
//           "annualType": "Checking",
//           "annualRate%": 0,
//           "years": 1,
//           "amount": 1
//         },
//         {
//           "account": "1256",
//           "annualType": "CD",
//           "annualRate%": 24,
//           "years": 2,
//           "amount": 9900
//         },
//         {
//           "account": "1111",
//           "annualType": "CD",
//           "annualRate%": 1.01,
//           "years": 10,
//           "amount": 10000
//         }
//       ]
//     },
//     {
//       "name": "NICE",
//       "deposits": [
//         {
//           "account": "1234",
//           "annualType": "Brokered CD",
//           "annualRate%": 2.4,
//           "years": 7,
//           "amount": 10990
//         }
//       ]
//     },
//     {
//       "name": "ANGRY",
//       "deposits": [
//         {
//           "account": "1234",
//           "annualType": "Brokered CD",
//           "annualRate%": 2.4,
//           "years": 7,
//           "amount": 10990
//         },
//         {
//           "account": "9898",
//           "annualType": "CD",
//           "annualRate%": 2.22,
//           "years": 1,
//           "amount": 5500
//         },
//         {
//           "account": "HALF YEAR",
//           "annualType": "CD",
//           "annualRate%": 2.22,
//           "years": 0.5,
//           "amount": 5500
//         }
//       ]
//     }
//   ]
// }`
