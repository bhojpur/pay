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
	"time"

	"github.com/bhojpur/pay/pkg/merchant"
	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/charge"
	"github.com/stripe/stripe-go/refund"
)

// Stripe implements gomerchant.PaymetGateway interface.
type Stripe struct {
	Config *Config
}

var _ merchant.PaymentGateway = &Stripe{}

// Config stripe config
type Config struct {
	Key string
}

// New creates Stripe struct.
func New(config *Config) *Stripe {
	stripe.Key = config.Key

	return &Stripe{
		Config: config,
	}
}

var capture bool = false

func (*Stripe) Authorize(amount uint64, params merchant.AuthorizeParams) (merchant.AuthorizeResponse, error) {
	int64Amount := int64(amount)
	chargeParams := &stripe.ChargeParams{
		Amount:      &int64Amount,
		Currency:    &params.Currency,
		Description: &params.Description,
		Capture:     &capture,
	}
	chargeParams.AddMetadata("order_id", params.OrderID)

	if params.PaymentMethod != nil {
		if params.PaymentMethod.CreditCard != nil {
			chargeParams.SetSource(toStripeCC(params.Customer, params.PaymentMethod.CreditCard, params.BillingAddress))
		}
		if params.PaymentMethod.SavedCreditCard != nil {
			if len(params.PaymentMethod.SavedCreditCard.CustomerID) > 0 {
				chargeParams.Customer = &params.PaymentMethod.SavedCreditCard.CustomerID
			}
			chargeParams.SetSource(params.PaymentMethod.SavedCreditCard.CreditCardID)
		}
	}

	charge, err := charge.New(chargeParams)
	if charge != nil {
		return merchant.AuthorizeResponse{TransactionID: charge.ID}, err
	}
	return merchant.AuthorizeResponse{}, err
}

func (*Stripe) CompleteAuthorize(paymentID string, params merchant.CompleteAuthorizeParams) (merchant.CompleteAuthorizeResponse, error) {
	return merchant.CompleteAuthorizeResponse{}, nil
}

func (*Stripe) Capture(transactionID string, params merchant.CaptureParams) (merchant.CaptureResponse, error) {
	_, err := charge.Capture(transactionID, nil)
	return merchant.CaptureResponse{TransactionID: transactionID}, err
}

func (s *Stripe) Refund(transactionID string, amount uint, params merchant.RefundParams) (merchant.RefundResponse, error) {
	transaction, err := s.Query(transactionID)

	if err == nil {
		if transaction.Captured {
			int64Amount := int64(amount)
			_, err = refund.New(&stripe.RefundParams{
				Charge: &transactionID,
				Amount: &int64Amount,
			})
		} else {
			int64Amount := int64(transaction.Amount - int(amount))
			_, err = charge.Capture(transactionID, &stripe.CaptureParams{
				Amount: &int64Amount,
			})
		}
	}

	return merchant.RefundResponse{TransactionID: transactionID}, err
}

func (*Stripe) Void(transactionID string, params merchant.VoidParams) (merchant.VoidResponse, error) {
	refundParams := &stripe.RefundParams{
		Charge: &transactionID,
	}
	_, err := refund.New(refundParams)
	return merchant.VoidResponse{TransactionID: transactionID}, err
}

func (*Stripe) Query(transactionID string) (merchant.Transaction, error) {
	c, err := charge.Get(transactionID, nil)
	created := time.Unix(c.Created, 0)
	transaction := merchant.Transaction{
		ID:        c.ID,
		Amount:    int(c.Amount - c.AmountRefunded),
		Currency:  string(c.Currency),
		Captured:  c.Captured,
		Paid:      c.Paid,
		Cancelled: c.Refunded,
		Status:    c.Status,
		CreatedAt: &created,
	}

	if transaction.Cancelled {
		transaction.Paid = false
		transaction.Captured = false
	}

	return transaction, err
}

func toStripeCC(customer string, cc *merchant.CreditCard, billingAddress *merchant.Address) *stripe.CardParams {
	var (
		expMonth = fmt.Sprint(cc.ExpMonth)
		expYear  = fmt.Sprint(cc.ExpYear)
	)
	cm := stripe.CardParams{
		Customer: &customer,
		Name:     &cc.Name,
		Number:   &cc.Number,
		ExpMonth: &expMonth,
		ExpYear:  &expYear,
		CVC:      &cc.CVC,
	}

	if billingAddress != nil {
		cm.AddressLine1 = &billingAddress.Address1
		cm.AddressLine1 = &billingAddress.Address2
		cm.AddressCity = &billingAddress.City
		cm.AddressState = &billingAddress.State
		cm.AddressZip = &billingAddress.ZIP
		cm.AddressCountry = &billingAddress.Country
	}

	return &cm
}
