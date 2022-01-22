package paygent_test

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
	"os"
	"testing"
	"time"

	cfgsvr "github.com/bhojpur/configure/pkg/markup"

	"github.com/bhojpur/pay/pkg/gateway/paygent"
	"github.com/bhojpur/pay/pkg/merchant"
	"github.com/bhojpur/pay/tests"
)

var Paygent *paygent.Paygent

type Config struct {
	MerchantID      string `required:"true"`
	ConnectID       string `required:"true"`
	ConnectPassword string `required:"true"`
	TelegramVersion string `required:"true" default:"1.0"`

	ClientFilePath string `required:"true" default:"paygent.pem"`
	CertPassword   string `required:"true" default:"changeit"`
	CAFilePath     string `required:"true" default:"curl-ca-bundle.crt"`

	ProductionMode  bool
	SecurityCodeUse bool
}

func init() {
	var config = &Config{}
	if err := cfgsvr.New(&cfgsvr.Config{ENVPrefix: "PAYGENT_CONFIG"}).Load(config); err != nil {
		fmt.Println(config)
		os.Exit(1)
	}

	Paygent = paygent.New(&paygent.Config{
		MerchantID:      config.MerchantID,
		ConnectID:       config.ConnectID,
		ConnectPassword: config.ConnectPassword,
		ClientFilePath:  config.ClientFilePath,
		CertPassword:    config.CertPassword,
		CAFilePath:      config.CAFilePath,
		ProductionMode:  config.ProductionMode,
		SecurityCodeUse: config.SecurityCodeUse,
	})
}

func TestTestSuite(t *testing.T) {
	tests.TestSuite{
		CreditCardManager: Paygent,
		Gateway:           Paygent,
		GetRandomCustomerID: func() string {
			return fmt.Sprint(time.Now().Unix())
		},
	}.TestAll(t)
}

func Test3DAuthorizeAndCapture(t *testing.T) {
	cards := map[string]bool{
		"5123459358515820": true,
		"5123459358515821": false,
	}

	for card, is3D := range cards {
		authorizeResult, err := Paygent.SecureCodeAuthorize(100,
			paygent.SecureCodeParams{
				UserAgent: "User-Agent	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12) AppleWebKit/602.3.12 (KHTML, like Gecko) Version/10.0.2 Safari/602.3.12",
				TermURL:    "https://www.bhojpur-consulting.com/p/terms-of-service.html",
				HttpAccept: "http",
			},
			merchant.AuthorizeParams{
				Currency: "JPY",
				OrderID:  fmt.Sprint(time.Now().Unix()),
				PaymentMethod: &merchant.PaymentMethod{
					CreditCard: &merchant.CreditCard{
						Name:     "JCB Card",
						Number:   card,
						ExpMonth: 1,
						ExpYear:  uint(time.Now().Year() + 1),
						CVC:      "1234",
					},
				},
			})

		if err != nil || authorizeResult.TransactionID == "" {
			t.Error(err, authorizeResult)
		}

		if is3D != authorizeResult.HandleRequest {
			t.Errorf("HandleRequest for card %v should be %v, but got %v", card, is3D, authorizeResult.HandleRequest)
		}

		if is3D {
			if result, ok := authorizeResult.Get("out_acs_html"); !ok || result.(string) == "" {
				t.Errorf("should get HTML, but %v", authorizeResult)
			}
		}
	}
}
