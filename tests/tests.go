package tests

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
	"testing"
	"time"

	"github.com/bhojpur/pay/pkg/merchant"
)

type TestSuite struct {
	CreditCardManager   merchant.CreditCardManager
	Gateway             merchant.PaymentGateway
	GetRandomCustomerID func() string
}

func (testSuite TestSuite) TestAll(t *testing.T) {
	testSuite.TestCreateCreditCard(t)
	testSuite.TestAuthorizeAndCapture(t)
	testSuite.TestAuthorizeAndCaptureWithSavedCreditCard(t)
	testSuite.TestRefund(t)
	testSuite.TestVoid(t)

	testSuite.TestListCreditCards(t)
	testSuite.TestListCreditCardsWithNoResult(t)
	testSuite.TestGetCreditCard(t)
	testSuite.TestDeleteCreditCard(t)
}

func (testSuite TestSuite) createSavedCreditCard() (merchant.CreditCardResponse, error) {
	return testSuite.CreditCardManager.CreateCreditCard(merchant.CreateCreditCardParams{
		CustomerID: testSuite.GetRandomCustomerID(),
		CreditCard: &merchant.CreditCard{
			Name:     "VISA",
			Number:   "4242424242424242",
			ExpMonth: 1,
			ExpYear:  uint(time.Now().Year() + 1),
			CVC:      "1234",
		},
	})
}

func (testSuite TestSuite) TestCreateCreditCard(t *testing.T) {
	if result, err := testSuite.createSavedCreditCard(); err != nil || result.CreditCardID == "" {
		t.Error(err, result)
	}
}

func (testSuite TestSuite) TestAuthorizeAndCapture(t *testing.T) {
	authorizeResult, err := testSuite.Gateway.Authorize(100, merchant.AuthorizeParams{
		Currency: "JPY",
		OrderID:  fmt.Sprint(time.Now().Unix()),
		PaymentMethod: &merchant.PaymentMethod{
			CreditCard: &merchant.CreditCard{
				Name:     "VISA",
				Number:   "4242424242424242",
				ExpMonth: 1,
				ExpYear:  uint(time.Now().Year() + 1),
				CVC:      "1234",
			},
		},
	})

	authorizeResult, err = testSuite.Gateway.Authorize(100, merchant.AuthorizeParams{
		Currency: "JPY",
		OrderID:  fmt.Sprint(time.Now().Unix()),
		PaymentMethod: &merchant.PaymentMethod{
			CreditCard: &merchant.CreditCard{
				Name:     "VISA",
				Number:   "4242424242424242",
				ExpMonth: 1,
				ExpYear:  uint(time.Now().Year() + 1),
				CVC:      "1234",
			},
		},
	})

	if err != nil || authorizeResult.TransactionID == "" {
		t.Error(err, authorizeResult)
	}

	captureResult, err := testSuite.Gateway.Capture(authorizeResult.TransactionID, merchant.CaptureParams{})

	if err != nil || captureResult.TransactionID == "" {
		t.Error(err, captureResult)
	}
}

func (testSuite TestSuite) TestAuthorizeAndCaptureWithSavedCreditCard(t *testing.T) {
	if savedCreditCard, err := testSuite.createSavedCreditCard(); err == nil {
		authorizeResult, err := testSuite.Gateway.Authorize(100, merchant.AuthorizeParams{
			Currency: "JPY",
			OrderID:  fmt.Sprint(time.Now().Unix()),
			PaymentMethod: &merchant.PaymentMethod{
				SavedCreditCard: &merchant.SavedCreditCard{
					CustomerID:   savedCreditCard.CustomerID,
					CreditCardID: savedCreditCard.CreditCardID,
					CVC:          "1234",
				},
			},
		})

		if err != nil || authorizeResult.TransactionID == "" {
			t.Error(err, authorizeResult)
		}

		captureResult, err := testSuite.Gateway.Capture(authorizeResult.TransactionID, merchant.CaptureParams{})

		if err != nil || captureResult.TransactionID == "" {
			t.Error(err, captureResult)
		}
	}
}

func (testSuite TestSuite) createAuth() merchant.AuthorizeResponse {
	authorizeResponse, _ := testSuite.Gateway.Authorize(1000, merchant.AuthorizeParams{
		Currency: "JPY",
		OrderID:  fmt.Sprint(time.Now().Unix()),
		PaymentMethod: &merchant.PaymentMethod{
			CreditCard: &merchant.CreditCard{
				Name:     "VISA",
				Number:   "4242424242424242",
				ExpMonth: 1,
				ExpYear:  uint(time.Now().Year() + 1),
				CVC:      "1235",
			},
		},
	})

	return authorizeResponse
}

func (testSuite TestSuite) TestRefund(t *testing.T) {
	// refund authorized transaction
	authorizeResponse := testSuite.createAuth()
	if refundResponse, err := testSuite.Gateway.Refund(authorizeResponse.TransactionID, 100, merchant.RefundParams{}); err == nil {
		if transaction, err := testSuite.Gateway.Query(refundResponse.TransactionID); err == nil {
			if !(transaction.Amount == 900 && transaction.Paid == true && transaction.Cancelled == false && transaction.CreatedAt != nil) { // &&transaction.Captured == false) {
				t.Errorf("transaction after refund authorized transaction is not correct, but got %#v", transaction)
			}
		} else {
			t.Errorf("no error should happen when query transaction, but got %v, %#v", err, transaction)
		}
	} else {
		t.Errorf("no error should happen when refund transaction, but got %v", err)
	}

	// refund authorized transaction, and capture it
	authorizeResponse = testSuite.createAuth()
	if refundResponse, err := testSuite.Gateway.Refund(authorizeResponse.TransactionID, 150, merchant.RefundParams{Captured: true}); err == nil {
		if transaction, err := testSuite.Gateway.Query(refundResponse.TransactionID); err == nil {
			if !(transaction.Amount == 850 && transaction.Paid == true && transaction.Captured == true && transaction.Cancelled == false && transaction.CreatedAt != nil) {
				t.Errorf("transaction after refund authorized transaction is not correct, but got %#v", transaction)
			}
		} else {
			t.Errorf("no error should happen when query transaction, but got %v, %#v", err, transaction)
		}
	} else {
		t.Errorf("no error should happen when refund transaction, but got %v", err)
	}

	// refund captured transaction
	authorizeResponse = testSuite.createAuth()
	captureResponse, _ := testSuite.Gateway.Capture(authorizeResponse.TransactionID, merchant.CaptureParams{})
	if refundResponse, err := testSuite.Gateway.Refund(captureResponse.TransactionID, 200, merchant.RefundParams{Captured: true}); err == nil {
		if transaction, err := testSuite.Gateway.Query(refundResponse.TransactionID); err == nil {
			if !(transaction.Amount == 800 && transaction.Paid == true && transaction.Captured == true && transaction.Cancelled == false && transaction.CreatedAt != nil) {
				t.Errorf("transaction after refund captured transaction is not correct, but got %#v", transaction)
			}
		} else {
			t.Errorf("no error should happen when query transaction, but got %v, %#v", err, transaction)
		}
	} else {
		t.Errorf("no error should happen when refund transaction, but got %v", err)
	}
}

func (testSuite TestSuite) TestVoid(t *testing.T) {
	// void authorized transaction
	authorizeResponse := testSuite.createAuth()
	if refundResponse, err := testSuite.Gateway.Void(authorizeResponse.TransactionID, merchant.VoidParams{}); err == nil {
		if transaction, err := testSuite.Gateway.Query(refundResponse.TransactionID); err == nil {
			if !(transaction.Paid == false && transaction.Captured == false && transaction.Cancelled == true && transaction.CreatedAt != nil) { // && transaction.Amount == 1000) {
				t.Errorf("transaction after refund auth is not correct, but got %#v", transaction)
			}
		} else {
			t.Errorf("no error should happen when query transaction, but got %v, %#v", err, transaction)
		}
	} else {
		t.Errorf("no error should happen when refund transaction, but got %v", err)
	}

	// void captured transaction
	authorizeResponse = testSuite.createAuth()
	captureResponse, _ := testSuite.Gateway.Capture(authorizeResponse.TransactionID, merchant.CaptureParams{})
	if refundResponse, err := testSuite.Gateway.Void(captureResponse.TransactionID, merchant.VoidParams{Captured: true}); err == nil {
		if transaction, err := testSuite.Gateway.Query(refundResponse.TransactionID); err == nil {
			if !(transaction.Paid == false && transaction.Captured == false && transaction.Cancelled == true && transaction.CreatedAt != nil) { // && transaction.Amount == 1000) {
				t.Errorf("transaction after refund captured is not correct, but got %#v", transaction)
			}
		} else {
			t.Errorf("no error should happen when query transaction, but got %v, %#v", err, transaction)
		}
	} else {
		t.Errorf("no error should happen when refund transaction, but got %v", err)
	}
}

func (testSuite TestSuite) TestListCreditCards(t *testing.T) {
	if response, err := testSuite.createSavedCreditCard(); err == nil {
		// create another credit card
		_, err = testSuite.CreditCardManager.CreateCreditCard(merchant.CreateCreditCardParams{
			CustomerID: response.CustomerID,
			CreditCard: &merchant.CreditCard{
				Name:     "VISA",
				Number:   "4242424242424241",
				ExpMonth: 1,
				ExpYear:  uint(time.Now().Year() + 1),
				CVC:      "5678",
			},
		})
		if err != nil {
			t.Errorf("should not get err, but got %v ", err)
		}

		if response, err := testSuite.CreditCardManager.ListCreditCards(merchant.ListCreditCardsParams{CustomerID: response.CustomerID}); err == nil {
			if len(response.CreditCards) != 2 {
				t.Errorf("Should found two saved credit cards, but got %v", response.CreditCards)
				for _, c := range response.CreditCards {
					fmt.Printf("%+v\n", *c)
				}
			}

			for _, creditCard := range response.CreditCards {
				if creditCard.MaskedNumber == "" || creditCard.ExpYear == 0 || creditCard.ExpMonth == 0 || creditCard.CustomerID == "" || creditCard.CreditCardID == "" {
					t.Errorf("Credit card's information should be correct, but got %v", creditCard)
				}
			}
		} else {
			t.Errorf("no error should happen when query saved credit cards, but got %v", err)
		}
	}
}

func (testSuite TestSuite) TestListCreditCardsWithNoResult(t *testing.T) {
	if response, err := testSuite.CreditCardManager.ListCreditCards(merchant.ListCreditCardsParams{CustomerID: testSuite.GetRandomCustomerID()}); err != nil {
		t.Errorf("should not return error, but got %v", err)
	} else if len(response.CreditCards) != 0 {
		t.Errorf("credit card's count should be zero")
	}
}

func (testSuite TestSuite) TestGetCreditCard(t *testing.T) {
	if response, err := testSuite.createSavedCreditCard(); err == nil {
		if response, err := testSuite.CreditCardManager.GetCreditCard(merchant.GetCreditCardParams{CustomerID: response.CustomerID, CreditCardID: response.CreditCardID}); err == nil {
			creditCard := response.CreditCard
			if creditCard == nil {
				t.Errorf("Should found saved credit cards, but got %v", response)
			} else if creditCard.Brand == "" || creditCard.MaskedNumber == "" || creditCard.ExpYear == 0 || creditCard.ExpMonth == 0 || creditCard.CustomerID == "" || creditCard.CustomerName == "" || creditCard.CreditCardID == "" {
				t.Errorf("Credit card's information should be correct, but got %v", creditCard)
			}
		} else {
			t.Errorf("no error should happen when query saved credit card, but got %v", err)
		}
	}
}

func (testSuite TestSuite) TestDeleteCreditCard(t *testing.T) {
	if response, err := testSuite.createSavedCreditCard(); err == nil {
		if _, err := testSuite.CreditCardManager.DeleteCreditCard(merchant.DeleteCreditCardParams{CustomerID: response.CustomerID, CreditCardID: response.CreditCardID}); err == nil {
			if response, err := testSuite.CreditCardManager.GetCreditCard(merchant.GetCreditCardParams{CustomerID: response.CustomerID, CreditCardID: response.CreditCardID}); err == nil {
				t.Errorf("Should failed to get credit card, but got %v", response)
			}
		} else {
			t.Errorf("no error should happen when delete saved credit card, but got %v", err)
		}
	}
}
