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
	"errors"
	"net/http"

	"github.com/bhojpur/pay/pkg/merchant"
)

type SecureCodeParams struct {
	UserAgent  string
	TermURL    string
	HttpAccept string
}

func (paygent *Paygent) SecureCodeAuthorize(amount uint64, secureCodeParams SecureCodeParams, params merchant.AuthorizeParams) (merchant.AuthorizeResponse, error) {
	if params.Params == nil {
		params.Params = merchant.Params{}
	}
	params.Set("Paygent3DMode", true)
	params.Set("Paygent3DParams", secureCodeParams)

	return paygent.Authorize(amount, params)
}

func (paygent *Paygent) CompleteAuthorize(paymentID string, params merchant.CompleteAuthorizeParams) (merchant.CompleteAuthorizeResponse, error) {
	if req, ok := params.Get("request"); ok {
		if request, ok := req.(*http.Request); ok {
			request.ParseForm()
			response, err := paygent.Request("024", merchant.Params{"MD": request.Form.Get("MD"), "PaRes": request.Form.Get("PaRes")})
			return merchant.CompleteAuthorizeResponse{Params: response.Params}, err
		}
	}
	return merchant.CompleteAuthorizeResponse{}, errors.New("no valid request params found")
}
