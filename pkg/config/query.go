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
	"log"
)

const Payments_key = "sk_test_acca9e854e36539a5e88d3b2b635870e2cb9a3f1"

func CreateAccountTable(db *sql.DB) {

	query := `CREATE TABLE IF NOT EXISTS accounts (
		id serial PRIMARY KEY,
		username VARCHAR ( 50 ) UNIQUE NOT NULL,
		email VARCHAR ( 50 ) UNIQUE NOT NULL,
		password VARCHAR ( 200 ) NOT NULL,
		fullname VARCHAR ( 200 ),
		gender VARCHAR ( 200 ),
		activated BOOL,
		created_on TIMESTAMP,
		last_login TIMESTAMP
		);`

	_, err := db.Exec(query)
	if err != nil {
		log.Println(err)
	}
}

// Droptable which is used to drop the accounts table is currently not in use
func DropTable(db *sql.DB) {
	query := `DROP TABLE IF EXISTS accounts`

	_, err := db.Exec(query)
	if err != nil {
		log.Println(err)
	}

	// fmt.Println("\n\n Table dropped successfully")
}

func CreateTransactionTable(db *sql.DB) {

	query := `CREATE TABLE IF NOT EXISTS transactions (
		id serial PRIMARY KEY,
		email VARCHAR ( 50 ) NOT NULL,
		amount INT NOT NULL,
		payment_status BOOL,
		access_code VARCHAR ( 200 ),
		authorization_url VARCHAR ( 200 ),
		reference VARCHAR ( 200 ),
		payment_channel VARCHAR ( 200 ),
		transaction_date TIMESTAMP,
		verify_status BOOL,
		FOREIGN KEY(email)
			REFERENCES accounts(email)
			ON DELETE CASCADE
		);`

	_, err := db.Exec(query)
	if err != nil {
		log.Println(err)
	}
}

func CreateBalanceTable(db *sql.DB) {

	query := `CREATE TABLE IF NOT EXISTS balance (
		id serial PRIMARY KEY,
		email VARCHAR ( 50 ) UNIQUE NOT NULL,
		current_balance INT NOT NULL,
		last_update TIMESTAMP,
		FOREIGN KEY(email)
			REFERENCES accounts(email)
			ON DELETE CASCADE
	);`

	_, err := db.Exec(query)
	if err != nil {
		log.Println(err)
	}
}
