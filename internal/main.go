package main

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
	"log"
	"net/http"

	"github.com/bhojpur/pay/pkg/account"
	"github.com/bhojpur/pay/pkg/config"
	"github.com/bhojpur/pay/pkg/transaction"
	"github.com/bhojpur/pay/pkg/utils"
)

func makemigrations() {
	config.CreateAccountTable(config.Db)
	config.CreateTransactionTable(config.Db)
	config.CreateBalanceTable(config.Db)
}

func main() {
	defer config.Db.Close()
	makemigrations()

	http.HandleFunc("/healthcheck", account.HealthCheck)
	http.HandleFunc("/register", account.RegisterUser)
	http.HandleFunc("/login", account.LoginUser)
	http.HandleFunc("/profile", account.UserProfile)
	http.HandleFunc("/refresh", account.RefreshToken)
	http.HandleFunc("/logout", account.Logout)
	http.HandleFunc("/fund", transaction.FundAccount)
	http.HandleFunc("/transactions", transaction.TransactionHistory)
	http.HandleFunc("/verify", transaction.VerifyTransaction)
	http.HandleFunc("/current_balance", transaction.GetBalance)
	http.HandleFunc("/transfer", transaction.TransferFunds)
	http.HandleFunc("/pay", transaction.PayForItem)
	if err := http.ListenAndServe(":"+utils.GetServerAddress(), nil); err != http.ErrServerClosed {
		log.Println(err)
	}
}
