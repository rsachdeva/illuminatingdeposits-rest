package interestcal

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/rsachdeva/illuminatingdeposits-rest/interestcal/interestvalue"
	"github.com/rsachdeva/illuminatingdeposits-rest/muxhttp"
	"github.com/stretchr/testify/require"
)

func TestJsonDecodeCal(t *testing.T) {
	tt := []struct {
		name      string
		payload   *strings.Reader
		checkFunc func(err error, rec *httptest.ResponseRecorder)
	}{
		{
			name:    "FormatCorrect",
			payload: strings.NewReader(formatCorrectJson),
			checkFunc: func(err error, rec *httptest.ResponseRecorder) {
				require.Nil(t, err)

				var ciresp interestvalue.CreateInterestResponse
				decoder := json.NewDecoder(rec.Body)
				decoder.DisallowUnknownFields()
				err = decoder.Decode(&ciresp)

				require.Nil(t, err)
				require.Equal(t, 23.46, ciresp.Banks[0].Deposits[2].Delta, "delta for a deposit in a bank must match")
				require.Equal(t, 259.86, ciresp.Banks[0].Delta, "delta for a bank must match")
				require.Equal(t, 336.74, ciresp.Delta, "overall delta for all deposists in all banks must match")
			},
		},
		{
			name:    "FormatInCorrect",
			payload: strings.NewReader(formatIncorrectJson),
			checkFunc: func(err error, rec *httptest.ResponseRecorder) {
				require.NotNil(t, err)
				require.Regexp(t, regexp.MustCompile(`ecoding new interest calculation request with banks and deposits: json: unknown field "institutions"`), err)
			},
		},
		{
			name:    "DeltaCalErrNeedsMin30days",
			payload: strings.NewReader(err30DayMinJson),
			checkFunc: func(err error, rec *httptest.ResponseRecorder) {
				require.NotNil(t, err)
				require.Regexp(t, regexp.MustCompile("creating new interest calculations: calculation for Account: 1256: NewDeposit period in years 0.0001 should not be less than 30 days"), err)
			},
		},
	}
	for _, tc := range tt {
		tc := tc // capture range variable https://golang.org/pkg/testing/#hdr-Subtests_and_Sub_benchmarks
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			method := "POST"
			payload := tc.payload
			t.Parallel()

			url := fmt.Sprintf("%v/v1/interests", "https://localhost:2828")
			req, _ := http.NewRequest(method, url, payload)
			svc := service{
				log: log.New(os.Stderr, "", log.LstdFlags),
			}
			v := muxhttp.Values{
				TraceID: "100",
			}
			ctx := context.WithValue(context.Background(), muxhttp.KeyValues, &v)
			err := svc.CreateInterest(ctx, rec, req)

			tc.checkFunc(err, rec)
		})
	}
}

const (
	formatCorrectJson = `
             {
			   "banks": 
                 [
				   {
				     "name": "HAPPIEST",
				     "deposits": [
					{
					  "account": "1234",
					  "annualType": "Checking",
					  "annualRate%": 0,
					  "years": 1,
					  "amount": 100
					},
					{
					  "account": "1256",
					  "annualType": "CD",
					  "annualRate%": 24,
					  "years": 2,
					  "amount": 10700
					},
					{
					  "account": "1111",
					  "annualType": "CD",
					  "annualRate%": 1.01,
					  "years": 10,
					  "amount": 27000
					}
				  ]
				},
				{
				  "name": "NICE",
				  "deposits": [
					{
					  "account": "1234",
					  "annualType": "Brokered CD",
					  "annualRate%": 2.4,
					  "years": 7,
					  "amount": 10990
					}
				  ]
				},
				{
				  "name": "ANGRY",
				  "deposits": [
					{
					  "account": "1234",
					  "annualType": "Brokered CD",
					  "annualRate%": 5,
					  "years": 7,
					  "amount": 10990
					},
					{
					  "account": "9898",
					  "annualType": "CD",
					  "annualRate%": 2.22,
					  "years": 1,
					  "amount": 5500
					}
				  ]
				}
			  ]
             }`
	formatIncorrectJson = `{
			  "institutions": [
				{
				  "name": "HAPPIEST",
				  "collectedamount": [
					{
					  "account": "1234",
					  "annualType": "Checking",
					  "annualRate%": 0,
					  "years": 1,
					  "amount": 100
					},
					{
					  "account": "1256",
					  "annualType": "CD",
					  "annualRate%": 24,
					  "years": 2,
					  "amount": 10700
					},
					{
					  "account": "1111",
					  "annualType": "CD",
					  "annualRate%": 1.01,
					  "years": 10,
					  "amount": 27000
					}
				  ]
				}
			  ]
			}`
	err30DayMinJson = `{
			  "banks": [
				{
				  "name": "HAPPIEST",
				  "deposits": [
					{
					  "account": "1234",
					  "annualType": "Checking",
					  "annualRate%": 0,
					  "years": 1,
					  "amount": 100
					},
					{
					  "account": "1256",
					  "annualType": "CD",
					  "annualRate%": 24,
					  "years": 0.0001,
					  "amount": 10700
					},
					{
					  "account": "1111",
					  "annualType": "CD",
					  "annualRate%": 1.01,
					  "years": 10,
					  "amount": 27000
					}
				  ]
				}
			  ]
			}`
)
