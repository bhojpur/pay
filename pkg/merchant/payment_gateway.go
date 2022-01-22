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

import "net/http"

// PaymentGateway interface
type PaymentGateway interface {
	Authorize(amount uint64, params AuthorizeParams) (AuthorizeResponse, error)
	CompleteAuthorize(paymentID string, params CompleteAuthorizeParams) (CompleteAuthorizeResponse, error)
	Capture(transactionID string, params CaptureParams) (CaptureResponse, error)
	Refund(transactionID string, amount uint, params RefundParams) (RefundResponse, error)
	Void(transactionID string, params VoidParams) (VoidResponse, error)

	Query(transactionID string) (Transaction, error)
}

// AuthorizeParams authorize params
type AuthorizeParams struct {
	Amount          uint64
	Currency        string
	Customer        string
	Description     string
	OrderID         string
	BillingAddress  *Address
	ShippingAddress *Address
	PaymentMethod   *PaymentMethod
	Params
}

// AuthorizeResponse authorize response
type AuthorizeResponse struct {
	TransactionID  string
	HandleRequest  bool                                                   // need process request after authorize or not
	RequestHandler func(http.ResponseWriter, *http.Request, Params) error // process request
	Params
}

// CompleteAuthorizeParams complete authorize params
type CompleteAuthorizeParams struct {
	Params
}

// CompleteAuthorizeResponse complete authorize response
type CompleteAuthorizeResponse struct {
	Params
}

// CaptureParams capture params
type CaptureParams struct {
	Params
}

// CaptureResponse capture response
type CaptureResponse struct {
	TransactionID string
	Params
}

// RefundParams refund params
type RefundParams struct {
	Captured bool
	Params
}

// RefundResponse refund response
type RefundResponse struct {
	TransactionID string
	Params
}

// VoidParams void params
type VoidParams struct {
	Captured bool
	Params
}

// VoidResponse void response
type VoidResponse struct {
	TransactionID string
	Params
}
