package database

import (
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
)

func TestCreateDBConnection(t *testing.T) {

	//Load .env file
	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatalf("Error loading .env file")
	}

	var dbGlobal *Database

	t.Run("test create connection", func(t *testing.T) {
		db := Database{
			Host:     os.Getenv("DB_HOST"),
			Username: os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			DbName:   os.Getenv("DB_NAME"),
			DbPort:   os.Getenv("DB_PORT"),
		}

		err := db.Connect()
		assert.Nil(t, err)

		dbGlobal = &db

	})

	t.Run("test connection already established", func(t *testing.T) {
		err := dbGlobal.Connect()
		assert.Nil(t, err)
	})

	t.Run("test get active connection", func(t *testing.T) {
		assert.NotNil(t, dbGlobal.GetConnection())
	})

	t.Run("test can't establish connection using invalid host ip", func(t *testing.T) {
		db := Database{
			Host:     "0",
			Username: os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			DbName:   os.Getenv("DB_NAME"),
			DbPort:   os.Getenv("DB_PORT"),
		}

		err := db.Connect()
		assert.NotNil(t, err)
	})
}
