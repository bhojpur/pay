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

import "testing"

func TestCreditCardLuhnAlgorithm(t *testing.T) {
	validNumbers := []string{"4111111111111111", "5431111111111111", "341111111111111", "6011601160116611", "5105105105105100", "5555555555554444", "4222222222222", "378282246310005", "371449635398431", "378734493671000", "38520000023237", "30569309025904", "6011111111111117", "6011000990139424", "3530111333300000", "3566002020360505"}

	for _, number := range validNumbers {
		creditCard := CreditCard{Number: number}
		if !creditCard.ValidNumber() {
			t.Errorf("%v should be valid", number)
		}
	}
}
