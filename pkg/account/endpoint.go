package account

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
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/bhojpur/pay/pkg/auth"
	"github.com/bhojpur/pay/pkg/config"
	"github.com/bhojpur/pay/pkg/transaction"
)

type healthJSON struct {
	Name   string
	Active bool
}

type loginInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type accessToken struct {
	Access string `json:"access"`
}

type loginResponse struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Fullname  string    `json:"fullname"`
	Gender    string    `json:"gender"`
	Activated bool      `json:"activated"`
	CreatedOn time.Time `json:"created_on"`
	LastLogin time.Time `json:"last_login"`
	Token     string    `json:"token"`
}

func HealthCheck(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch req.Method {
	case http.MethodGet:
		w.WriteHeader(http.StatusOK)
		resp := &healthJSON{
			Name:   "Bhojpur Wallet API, up and running",
			Active: true,
		}
		jsonResp, _ := json.Marshal(resp)
		fmt.Fprint(w, string(jsonResp))
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprint(w, `{"message" : "Method not allowed"}`)
		return
	}

}

func RegisterUser(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch req.Method {
	case http.MethodPost:
		var (
			user Accounts
			err  error
		)

		_ = json.NewDecoder(req.Body).Decode(&user)
		if user.Username == "" || user.Password == "" || user.Email == "" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, `{"Message":"Username, Email and Password is required"}`)
			return
		}

		user.CreatedOn = time.Now()
		user.LastLogin = time.Now()
		user.Password, err = auth.HashPassword(user.Password)

		if err != nil {
			log.Println("Error occurred in hashing password, ", err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, `{"message" : "Something went wrong"}`)
			return
		}

		if created := AddRecordToAccounts(config.Db, user); created {
			w.WriteHeader(http.StatusCreated)
			fmt.Fprint(w, `{"Message" : "Successfully Created"}`)
			if !transaction.InitializeBalance(config.Db, user.Email) {
				log.Println("Could not initialize balance for user", user.Email)
				return
			}
		} else {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, `{"Message" : "User already exists, please login"}`)
			return
		}

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprint(w, `{"Message" : "Method not allowed"}`)
		return
	}

}

func LoginUser(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	switch req.Method {
	case http.MethodPost:
		var loginDet loginInfo

		_ = json.NewDecoder(req.Body).Decode(&loginDet)
		if loginDet.Username == "" || loginDet.Password == "" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, `{"Message":"Username and Password is required"}`)
			return
		}

		user, err := GetUserLogin(config.Db, loginDet.Username)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, `{"Message":"User does not exist"}`)
			return
		}

		if auth.CheckPasswordHash(loginDet.Password, user.Password) {
			accessToken, refreshToken, err := auth.GenerateAuthTokens(user.Email)
			if err != nil {
				fmt.Println(err)
			}
			b := loginResponse{
				ID:        user.ID,
				Username:  user.Username,
				Email:     user.Email,
				Fullname:  user.Fullname,
				Gender:    user.Gender,
				Activated: user.Activated,
				CreatedOn: user.CreatedOn,
				LastLogin: user.LastLogin,
				Token:     accessToken,
			}
			jsonResp, err := json.Marshal(b)
			if err != nil {
				log.Println(err)
			}

			expire := time.Now().Add(time.Hour * 6)
			cookie := http.Cookie{Name: "Refreshtoken", Value: refreshToken, Path: "/",
				Expires: expire, HttpOnly: true} // extra agruement, Secure : true, test this on deployment
			http.SetCookie(w, &cookie)
			fmt.Fprint(w, string(jsonResp))
			return
		} else {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, `{"Message":"Invalid credentials"}`)
			return
		}

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprint(w, `{"Message":"Method not allowed"}`)
		return
	}
}

func UserProfile(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	authorized, email, err := auth.CheckAuth(req)
	if !authorized {
		auth.UnAuthorizedResponse(w, err)
		return
	}

	switch req.Method {
	case http.MethodGet:

		user, err := GetUser(config.Db, fmt.Sprint(email))
		if err != nil {
			log.Println("Error occured in getting user ,", err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, `{"message" : "something went wrong"}`)
			return
		}

		b := loginResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			Fullname:  user.Fullname,
			Gender:    user.Gender,
			Activated: user.Activated,
			CreatedOn: user.CreatedOn,
			LastLogin: user.LastLogin,
		}

		jsonResp, err := json.Marshal(b)
		if err != nil {
			log.Println(err)
		}
		fmt.Fprint(w, string(jsonResp))
		return

	case http.MethodPut:
		var user Accounts

		_ = json.NewDecoder(req.Body).Decode(&user)
		if user.Fullname == "" && user.Gender == "" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, `{"Message":"Fullname or Gender is required"}`)
			return
		}

		user.Email = fmt.Sprint(email)
		err := EditUser(config.Db, &user)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, `{"Message" : "Something went wrong"}`)
			return
		}
		w.WriteHeader(http.StatusAccepted)
		fmt.Fprint(w, `{"Message" : "Successfully edited"}`)
		return

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprint(w, `{"Message" : "Method not allowed"}`)
		return
	}
}

func RefreshToken(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch req.Method {
	case http.MethodPost:
		token, err := req.Cookie("Refreshtoken")
		if err != nil {
			auth.UnAuthorizedResponse(w, errors.New(`{"Message" : "Credentials not sent"}`))
			return
		}

		if authorized, email, _ := auth.CheckRefreshToken(token.Value); authorized {
			accessString, err := auth.NewAccessToken(fmt.Sprint(email))
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprint(w, `{"Message" : "Could not generate accesstoken"}`)
				return
			}
			message := accessToken{accessString}
			jsonResp, err := json.Marshal(message)
			if err != nil {
				log.Println(err)
			}
			cookie := http.Cookie{Name: "Refreshtoken", Value: token.Value, Path: "/",
				HttpOnly: true} // extra agruement, Secure : true, test this on deployment
			http.SetCookie(w, &cookie)

			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, string(jsonResp))
			return
		}

		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, `{"Message" : "Please login"}`)
		return

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprint(w, `{"Message" : "Method not allowed"}`)
		return
	}

}

func Logout(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	authorized, _, err := auth.CheckAuth(req)
	if !authorized {
		auth.UnAuthorizedResponse(w, err)
		return
	}
	switch req.Method {
	case http.MethodPost:
		cookie := http.Cookie{Name: "Refreshtoken", Value: "", Path: "/",
			MaxAge: -1, HttpOnly: true}
		http.SetCookie(w, &cookie)
		w.WriteHeader(http.StatusNoContent)
		fmt.Fprint(w, `{"Message" : "Goodbye!"}`)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprint(w, `{"Message" : "Method not allowed"}`)
		return
	}

}
