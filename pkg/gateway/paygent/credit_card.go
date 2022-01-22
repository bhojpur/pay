package paygent

// Copyright (c) 2018 Bhojpur Consulting Private Limited, India. All rights reserved.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

import (
	"encoding/csv"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/bhojpur/pay/pkg/merchant"
)

func getValidTerm(creditCard *merchant.CreditCard) string {
	year := fmt.Sprint(creditCard.ExpYear)
	if len(year) >= 2 {
		year = year[len(year)-2:]
	} else {
		year = ""
	}
	return fmt.Sprintf("%02d", creditCard.ExpMonth) + year
}

var brandsMap = map[string]string{"visa": "V", "master": "M", "american_express": "X", "diners_club": "C", "jcb": "J"}

func (paygent *Paygent) CreateCreditCard(creditCardParams merchant.CreateCreditCardParams) (merchant.CreditCardResponse, error) {
	var (
		response   = merchant.CreditCardResponse{CustomerID: creditCardParams.CustomerID}
		creditCard = creditCardParams.CreditCard
		brand, _   = brandsMap[creditCard.Brand()]
	)

	params := merchant.Params{
		"customer_id":     creditCardParams.CustomerID,
		"card_number":     creditCard.Number,
		"card_valid_term": getValidTerm(creditCard),
		"cardholder_name": creditCard.Name,
		"card_brand":      brand,
	}.IgnoreBlankFields()
	if paygent.Config.SecurityCodeUse {
		params.Set("security_code_use", 1)
		params.Set("card_conf_number", creditCard.CVC)
	}
	results, err := paygent.Request("025", params)

	if err == nil {
		if customerCardID, ok := results.Get("customer_card_id"); ok {
			response.CreditCardID = fmt.Sprint(customerCardID)
		}
	}
	response.Params = results.Params

	return response, err
}

func (paygent *Paygent) GetCreditCard(getCreditCardParams merchant.GetCreditCardParams) (merchant.GetCreditCardResponse, error) {
	var response merchant.GetCreditCardResponse
	results, err := paygent.Request("027", merchant.Params{"customer_id": getCreditCardParams.CustomerID, "credit_card_id": getCreditCardParams.CreditCardID})

	if err == nil {
		cards, err := parseListCreditCardsResponse(&results)
		if len(cards) > 0 {
			response.CreditCard = cards[0]
		} else {
			err = errors.New("credit card not found")
		}
		return response, err
	}

	return response, err
}

func (paygent *Paygent) DeleteCreditCard(deleteCreditCardParams merchant.DeleteCreditCardParams) (merchant.DeleteCreditCardResponse, error) {
	var response = merchant.DeleteCreditCardResponse{}

	results, err := paygent.Request("026", merchant.Params{"customer_id": deleteCreditCardParams.CustomerID, "customer_card_id": deleteCreditCardParams.CreditCardID}.IgnoreBlankFields())
	response.Params = results.Params
	return response, err
}

func (paygent *Paygent) ListCreditCards(listCreditCardsParams merchant.ListCreditCardsParams) (merchant.ListCreditCardsResponse, error) {
	var response = merchant.ListCreditCardsResponse{}

	results, err := paygent.Request("027", merchant.Params{"customer_id": listCreditCardsParams.CustomerID})

	if err == nil {
		response.CreditCards, err = parseListCreditCardsResponse(&results)
	}

	if results.ResponseCode == "P026" {
		err = nil
	}

	return response, err
}

func parseListCreditCardsResponse(response *Response) (cards []*merchant.CustomerCreditCard, err error) {
	var headers []string

	for _, str := range strings.Split(response.RawBody, "\r\n") {
		if str != "" {
			row := csv.NewReader(strings.NewReader(str))
			if record, err := row.Read(); err == nil {
				if len(record) == 0 {
					return nil, errors.New("wrong format")
				}

				switch record[0] {
				case "1":
					// response information
					if len(record) != 4 {
						return nil, errors.New("wrong format")
					}

					response.Result = record[1]
					response.ResponseCode = record[2]
					response.ResponseDetail = record[3]
				case "2":
					// card header
					headers = record
				case "3":
					// card information
					params := merchant.Params{}
					for idx, value := range record {
						params[headers[idx]] = value
					}
					customerCard := &merchant.CustomerCreditCard{Params: params}

					if v, ok := params.Get("customer_id"); ok {
						customerCard.CustomerID = fmt.Sprint(v)
					}

					if v, ok := params.Get("customer_card_id"); ok {
						customerCard.CreditCardID = fmt.Sprint(v)
					}

					if v, ok := params.Get("cardholder_name"); ok {
						customerCard.CustomerName = fmt.Sprint(v)
					}

					if v, ok := params.Get("card_number"); ok {
						customerCard.MaskedNumber = fmt.Sprint(v)
					}

					if v, ok := params.Get("card_brand"); ok {
						for key, value := range brandsMap {
							if fmt.Sprint(v) == value {
								customerCard.Brand = key
							}
						}

						if customerCard.Brand == "" {
							customerCard.Brand = fmt.Sprint(v)
						}
					}

					if v, ok := params.Get("card_valid_term"); ok {
						if u, err := strconv.Atoi(fmt.Sprint(v)[0:2]); err == nil {
							customerCard.ExpMonth = uint(u)
						}
						if u, err := strconv.Atoi(fmt.Sprint(v)[2:4]); err == nil {
							customerCard.ExpYear = uint(time.Now().Year()/100*100 + u)
						}
					}

					cards = append(cards, customerCard)
				case "4":
					// card numbers
				}
			} else {
				return nil, err
			}
		}
	}

	if response.Result == "1" {
		if response.ResponseDetail != "" {
			err = errors.New(response.ResponseDetail)
		} else {
			err = errors.New("failed to process this request")
		}
	}

	return cards, err
}
