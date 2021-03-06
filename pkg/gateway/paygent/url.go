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

var TelegramServiceSandboxDomain = "https://sandbox.paygent.co.jp"
var TelegramServiceDomain = "https://service.paygent.co.jp"

var TelegramServiceURLs = map[string]string{
	// ###ATM決済URL###
	"01": "/n/atm/request",
	// ###クレジットカード決済URL1###
	"02": "/n/card/request",
	// ###クレジットカード決済URL2###
	"11": "/n/card/request",
	// ###コンビニ番号方式決済URL###
	"03": "/n/conveni/request",
	// ###コンビニ帳票方式決済URL###
	"04": "/n/conveni/request_print",
	// ###銀行ネット決済URL###
	"05": "/n/bank/request",
	// ###銀行ネット決済ASPURL###
	"06": "/n/bank/requestasp",
	// ###仮想口座決済URL###
	"07": "/n/virtualaccount/request",
	// ###決済情報照会URL###
	"09": "/n/ref/request",
	// ###決済情報差分照会URL###
	"091": "/n/ref/paynotice",
	// ###キャリア継続課金差分照会URL###
	"093": "/n/ref/runnotice",
	"094": "/n/ref/paymentref",
	// ###携帯キャリア決済URL###
	"10": "/n/c/request",
	// ###携帯キャリア決済URL（継続課金用）###
	"12": "/n/c/request",
	// ###ファイル決済URL###
	"20": "/n/o/requestdata",
	// ###PayPal決済URL###
	"13": "/n/paypal/request",
	// ###電子マネー決済URL###
	"15": "/n/emoney/request",
}
