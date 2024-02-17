package main

import (
	"fmt"

	"gorm.io/driver/postgres"
	_ "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB
)

func init() {
	// Connect to DB
	var err error
	DB, err = gorm.Open(postgres.Open("postgres://postgres:admin@localhost:5432/postgres?sslmode=disable"), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("failed to connect to DB: %v", err))
	}
}
