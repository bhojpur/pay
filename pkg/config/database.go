package config

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
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

const (
	POSTGRES_USER     = "bhojpur"
	POSTGRES_PASSWORD = "bhojpur"
	DB_NAME           = "wallet"
	POSTGRES_HOST     = "localhost"
	POSTGRES_PORT     = 5432
)

func databaseURL() string {

	dBUrl, present := os.LookupEnv("DATABASE_URL")
	if present {
		return dBUrl
	}
	postgres_conn := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable",
		POSTGRES_HOST, POSTGRES_PORT, POSTGRES_USER, POSTGRES_PASSWORD,
		DB_NAME)
	return postgres_conn
}

func ConnectDatabase() *sql.DB {

	db, err := sql.Open("postgres", databaseURL())
	if err != nil {
		log.Println(err)
		panic("Failed to connect to database")
	}

	bdErr := db.Ping()
	if bdErr != nil {
		panic(bdErr)
	}
	return db
}

var Db = ConnectDatabase()
