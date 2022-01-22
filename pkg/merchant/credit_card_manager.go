package merchant

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

// CreditCardManager interface
type CreditCardManager interface {
	CreateCreditCard(creditCardParams CreateCreditCardParams) (CreditCardResponse, error)
	GetCreditCard(creditCardParams GetCreditCardParams) (GetCreditCardResponse, error)
	ListCreditCards(listCreditCardsParams ListCreditCardsParams) (ListCreditCardsResponse, error)
	DeleteCreditCard(deleteCreditCardParams DeleteCreditCardParams) (DeleteCreditCardResponse, error)
}

// CreateCreditCard Params
type CreateCreditCardParams struct {
	CustomerID string
	CreditCard *CreditCard
}

type CreditCardResponse struct {
	CustomerID   string
	CreditCardID string
	Params
}

// Get Credit Cards Params
type GetCreditCardParams struct {
	CustomerID   string
	CreditCardID string
}

type GetCreditCardResponse struct {
	CreditCard *CustomerCreditCard
	Params
}

// Delete Credit Cards Params
type DeleteCreditCardParams struct {
	CustomerID   string
	CreditCardID string
}

type DeleteCreditCardResponse struct {
	Params
}

// List Credit Cards Params
type ListCreditCardsParams struct {
	CustomerID string
}

type ListCreditCardsResponse struct {
	CreditCards []*CustomerCreditCard
	Params
}

// CustomerCreditCard CustomerCard defination
type CustomerCreditCard struct {
	CustomerID   string
	CustomerName string
	CreditCardID string
	MaskedNumber string
	ExpMonth     uint
	ExpYear      uint
	Brand        string
	Params
}
