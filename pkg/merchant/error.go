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

import "errors"

var (
	ErrInvalidNumber      = errors.New("merchant: the card number is not a valid credit card number.")
	ErrInvalidExpiryMonth = errors.New("merchant: the card's expiration month is invalid.")
	ErrInvalidExpiryYear  = errors.New("merchant: the card's expiration year is invalid.")
	ErrInvalidCVC         = errors.New("merchant: the card's security code is invalid.")
	ErrIncorrectNumber    = errors.New("merchant: the card number is incorrect.")
	ErrExpiredCard        = errors.New("merchant: the card has expired.")
	ErrIncorrectCVC       = errors.New("merchant: the card's security code is incorrect.")
	ErrIncorrectZip       = errors.New("merchant: the card's zip code failed validation.")
	ErrCardDeclined       = errors.New("merchant: the card was declined.")
	ErrMissing            = errors.New("merchant: there is no card on a customer that is being charged.")
	ErrProcessingError    = errors.New("merchant: an error occurred while processing the card.")
)
