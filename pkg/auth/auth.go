package auth

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
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	if len(password) == 0 {
		return "", errors.New("can't has an empty string")
	}
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 4)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

var (
	signingKey        = []byte("3d67d77426d9878967a177437316554b0088fa88be95846252011528e8bad788")
	refreshSigningKey = []byte("b178604f6216f904f394641fd167078e426d5fe9ce20d4c07a65e8dd051a40d9")
)

// Returns accesstoken, refreshtoken and error
func GenerateAuthTokens(email string) (string, string, error) {

	if len(email) == 0 {
		return "", "", errors.New("can't generate tokens for empty email")
	}

	// Access token
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["client"] = email
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	accessString, err := token.SignedString(signingKey)

	if err != nil {
		return "", "", err
	}

	// Refresh token
	refreshToken := jwt.New(jwt.SigningMethodHS256)
	refreshClaims := refreshToken.Claims.(jwt.MapClaims)

	refreshClaims["authorized"] = true
	refreshClaims["client"] = email
	refreshClaims["exp"] = time.Now().Add(time.Hour * 6).Unix()

	refreshString, err := refreshToken.SignedString(refreshSigningKey)

	if err != nil {
		return "", "", err
	}

	return accessString, refreshString, nil
}

func CheckAuth(r *http.Request) (bool, interface{}, error) {
	if r.Header["Authorization"] != nil {
		if len(strings.Split(r.Header["Authorization"][0], " ")) < 2 {
			return false, "", errors.New("invalid credentials")
		}
		accessToken := strings.Split(r.Header["Authorization"][0], " ")[1]
		token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("an error occured")
			}
			return signingKey, nil
		})
		if err != nil {
			return false, "", err
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			return true, claims["client"], nil
		}
	}
	return false, "", errors.New("credentials not provided")
}

func UnAuthorizedResponse(w http.ResponseWriter, err error) {
	// cookie := http.Cookie{Name: "Refreshtoken", Value: "", Path: "/",
	// 	MaxAge: -1, HttpOnly: true}
	// http.SetCookie(w, &cookie)
	w.WriteHeader(http.StatusForbidden)
	fmt.Fprint(w, err.Error())
}

func CheckRefreshToken(refreshToken string) (bool, interface{}, error) {

	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("an error occured")
		}
		return refreshSigningKey, nil
	})
	if err != nil {
		return false, "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return true, claims["client"], nil
	}

	return false, "", errors.New("credentials not provided")
}

func NewAccessToken(email string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["client"] = email
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	accessString, err := token.SignedString(signingKey)

	if err != nil {
		return "", err
	}

	return accessString, nil
}
