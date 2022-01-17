package transaction

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
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/bhojpur/pay/pkg/auth"
	"github.com/bhojpur/pay/pkg/config"
)

func FundAccount(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	authorized, email, err := auth.CheckAuth(req)
	if !authorized {
		auth.UnAuthorizedResponse(w, err)
		return
	}

	switch req.Method {
	case http.MethodPost:

		var wg sync.WaitGroup
		wg.Add(1)
		// Get email and amount from post request, make sure value is positive
		var pl GeneratePayment
		_ = json.NewDecoder(req.Body).Decode(&pl)
		if pl.Amount <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, `{"Message" : "Positive Amount value is needed"}`)
			return
		}
		pl.Email = fmt.Sprint(email)
		// Hit Payments to return link
		result, err := HitPayments(pl.Email, pl.Amount*100)
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			fmt.Fprint(w, `{"Message" : "Something went wrong"}`)
			return
		}

		// Send details to db
		go func() {
			if !AddTransaction(config.Db, pl, result) {
				log.Println("Could not add transaction to database")
			}
			wg.Done()
		}()

		// Send the link to the user
		jsonresp, _ := json.Marshal(result.Data)
		fmt.Fprint(w, string(jsonresp))
		wg.Wait()
		return

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprint(w, `{"Message" : "Method not allowed"}`)
		return
	}
}

func TransactionHistory(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	authorized, email, err := auth.CheckAuth(req)
	if !authorized {
		auth.UnAuthorizedResponse(w, err)
		return
	}

	switch req.Method {
	case http.MethodGet:
		transactions := GetTransactions(config.Db, fmt.Sprint(email))
		if transactions == nil {
			fmt.Fprint(w, `[]`)
			return
		}
		jsonresp, _ := json.Marshal(transactions)
		fmt.Fprint(w, string(jsonresp))
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprint(w, `{"Message" : "Method not allowed"}`)
		return
	}
}

func VerifyTransaction(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	authorized, email, err := auth.CheckAuth(req)
	if !authorized {
		auth.UnAuthorizedResponse(w, err)
		return
	}

	switch req.Method {
	case http.MethodPost:
		var tranx Transactions

		_ = json.NewDecoder(req.Body).Decode(&tranx)
		if tranx.Reference == "" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, `{"Message" : "Reference is needed"}`)
			return
		}

		result, attempted := PaymentsVerify(tranx.Reference)
		result.Email = fmt.Sprint(email)

		if !attempted {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{"Message" : "Transaction was abandoned"}`)
		}

		// Write to db, update transaction table and balance table
		addedToDB, alreadyverified := UpdateTransaction(config.Db, result)
		if alreadyverified {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, `{"Message" : "Transaction has already been verified"}`)
			return
		}

		if !addedToDB {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, `{"Message" : "Something went wrong, please try again"}`)
			return
		}

		w.WriteHeader(http.StatusOK)
		switch addedToDB {
		case result.Payment_status:
			fmt.Fprint(w, `{"Message" : "Transaction was successful and has been updated accordingly"}`)
			return
		default:
			fmt.Fprint(w, `{"Message" : "Transaction failed"}`)
			return
		}

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprint(w, `{"Message" : "Method not allowed"}`)
		return
	}
}

func GetBalance(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	authorized, email, err := auth.CheckAuth(req)
	if !authorized {
		auth.UnAuthorizedResponse(w, err)
		return
	}

	switch req.Method {
	case http.MethodGet:
		userBal, err := GetCurrentBalance(config.Db, fmt.Sprint(email))
		if err != nil {
			log.Println("Error occured in getting balance ", err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, `{"mesage":"Something went wrong"}`)
			return
		}

		jsonResp, err := json.Marshal(userBal)
		if err != nil {
			log.Println("Error marshalling response")
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, string(jsonResp))
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprint(w, `{"Message" : "Method not allowed"}`)
		return
	}
}

func TransferFunds(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	authorized, email, err := auth.CheckAuth(req)
	if !authorized {
		auth.UnAuthorizedResponse(w, err)
		return
	}

	switch req.Method {
	case http.MethodPost:
		var pl recipient
		_ = json.NewDecoder(req.Body).Decode(&pl)

		if pl.Amount <= 0 || pl.Email == "" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, `{"Message" : "Amount/Email necessary"}`)
			return
		}

		sent, insufficientBal := sendMoney(config.Db, pl, fmt.Sprint(email))
		if insufficientBal {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, `{"Message" : "Insufficient Balance"}`)
			return
		}

		if sent {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{"Message" : "Funds has been sent"}`)
			return
		} else {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, `{"Message" : "User does not exist"}`)
			return
		}

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprint(w, `{"Message" : "Method not allowed"}`)
		return
	}
}

func PayForItem(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	authorized, email, err := auth.CheckAuth(req)
	if !authorized {
		auth.UnAuthorizedResponse(w, err)
		return
	}

	switch req.Method {
	case http.MethodPost:

		var amt GeneratePayment
		_ = json.NewDecoder(req.Body).Decode(&amt)

		if amt.Amount <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, `{"Message" : "Amount/Email necessary with positive value"}`)
			return
		}

		debited, insufficientBal := deductBalance(config.Db, amt.Amount, fmt.Sprint(email))
		if insufficientBal {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, `{"Message" : "Insufficient balance"}`)
			return
		}
		if debited {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{"Message" : "Paid successfully"}`)
			return
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, `{"Message" : "Something went wrong"}`)
			return
		}

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprint(w, `{"Message" : "Method not allowed"}`)
		return
	}
}
