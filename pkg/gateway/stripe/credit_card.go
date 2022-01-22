package stripe

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
	"fmt"

	"github.com/bhojpur/pay/pkg/merchant"
	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/card"
)

func (*Stripe) CreateCreditCard(creditCardParams merchant.CreateCreditCardParams) (merchant.CreditCardResponse, error) {
	var (
		expMonth = fmt.Sprint(creditCardParams.CreditCard.ExpMonth)
		expYear  = fmt.Sprint(creditCardParams.CreditCard.ExpYear)
	)

	c, err := card.New(&stripe.CardParams{
		Customer: &creditCardParams.CustomerID,
		Name:     &creditCardParams.CreditCard.Name,
		Number:   &creditCardParams.CreditCard.Number,
		ExpMonth: &expMonth,
		ExpYear:  &expYear,
		CVC:      &creditCardParams.CreditCard.CVC,
	})

	resp := merchant.CreditCardResponse{CreditCardID: c.ID}

	if c.Customer != nil {
		resp.CustomerID = c.Customer.ID
	}

	return resp, err
}

func (*Stripe) GetCreditCard(creditCardParams merchant.GetCreditCardParams) (merchant.GetCreditCardResponse, error) {
	c, err := card.Get(creditCardParams.CreditCardID, &stripe.CardParams{Customer: &creditCardParams.CustomerID})

	resp := merchant.GetCreditCardResponse{
		CreditCard: &merchant.CustomerCreditCard{
			CustomerName: c.Name,
			CreditCardID: c.ID,
			MaskedNumber: fmt.Sprint(c.Last4),
			ExpMonth:     uint(c.ExpMonth),
			ExpYear:      uint(c.ExpYear),
			Brand:        string(c.Brand),
		},
	}

	if c.Customer != nil {
		resp.CreditCard.CustomerID = c.Customer.ID
	}

	return resp, err
}

func (*Stripe) ListCreditCards(listCreditCardsParams merchant.ListCreditCardsParams) (merchant.ListCreditCardsResponse, error) {
	iter := card.List(&stripe.CardListParams{Customer: &listCreditCardsParams.CustomerID})
	resp := merchant.ListCreditCardsResponse{}
	for iter.Next() {
		c := iter.Card()
		customerCreditCard := &merchant.CustomerCreditCard{
			CustomerName: c.Name,
			CreditCardID: c.ID,
			MaskedNumber: fmt.Sprint(c.Last4),
			ExpMonth:     uint(c.ExpMonth),
			ExpYear:      uint(c.ExpYear),
			Brand:        string(c.Brand),
		}

		if c.Customer != nil {
			customerCreditCard.CustomerID = c.Customer.ID
		}

		resp.CreditCards = append(resp.CreditCards, customerCreditCard)
	}
	return resp, iter.Err()
}

func (*Stripe) DeleteCreditCard(deleteCreditCardParams merchant.DeleteCreditCardParams) (merchant.DeleteCreditCardResponse, error) {
	_, err := card.Del(deleteCreditCardParams.CreditCardID, &stripe.CardParams{Customer: &deleteCreditCardParams.CustomerID})
	return merchant.DeleteCreditCardResponse{}, err
}
