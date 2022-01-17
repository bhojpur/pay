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
	"testing"
)

func TestHashPassword(t *testing.T) {
	_, err := HashPassword("myStrongPassword")
	if err != nil {
		t.Fatalf("Failed to hash password with error %s", err)
	}

	_, err = HashPassword("")
	if err == nil {
		t.Fatalf("Hashing an invalid password")
	}
}

func TestCheckPasswordHash(t *testing.T) {

	hashedPassword, err := HashPassword("myStrongPassword")
	if err != nil {
		t.Errorf("Strange, unable to hash password with error %s", err)
	}

	if equal := CheckPasswordHash("myStrongPassword", hashedPassword); !equal {
		t.Fatalf("Checkpassword is not retrieving current password")
	}
}

func TestGenerateAuthToken(t *testing.T) {

	_, _, err := GenerateAuthTokens("shashi.rai@bhojpur.net")
	if err != nil {
		t.Fatalf("Could not generate tokens with error, %s", err)
	}

	_, _, err = GenerateAuthTokens("")
	if err == nil {
		t.Fatalf("Generating token for an empty string as email")
	}
}

func TestCheckRefreshToken(t *testing.T) {

	_, refreshToken, err := GenerateAuthTokens("shashi.rai@bhojpur.net")
	if err != nil {
		t.Errorf("Could not generate token with error, %s", err)
	}

	ok, _, err := CheckRefreshToken(refreshToken)
	if !ok {
		t.Fatalf("Could not refresh token with error, %s", err)
	}

	ok, _, err = CheckRefreshToken("invalid token")
	if ok {
		t.Fatalf("Refreshing an invalid token, %s", err)
	}
}
