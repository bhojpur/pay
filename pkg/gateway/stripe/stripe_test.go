package stripe_test

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

	cfgsvr "github.com/bhojpur/configure/pkg/markup"
	"github.com/bhojpur/pay/pkg/gateway/stripe"
	"github.com/bhojpur/pay/tests"
	"github.com/stripe/stripe-go/customer"
)

var Stripe *stripe.Stripe

type Config struct {
	Key string `required:"true"`
}

func init() {
	var config = &Config{}
	os.Setenv("CONFIGURE_ENV_PREFIX", "-")
	if err := cfgsvr.Load(config); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	Stripe = stripe.New(&stripe.Config{
		Key: config.Key,
	})
}

func TestTestSuite(t *testing.T) {
	tests.TestSuite{
		CreditCardManager: Stripe,
		Gateway:           Stripe,
		GetRandomCustomerID: func() string {
			Customer, err := customer.New(nil)
			if err != nil {
				fmt.Printf("Get error when create customer: %v", err)
			}
			return Customer.ID
		},
	}.TestAll(t)
}
