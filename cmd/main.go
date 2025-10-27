package main

import (
	"fmt"
	"log"

	"github.com/google/uuid"
	"go.etcd.io/bbolt"
)

/*
*
THIS IS MY PESONAL PROEJCT, LEARNING GO LANGUAGE AND THE CORE CONCEPTS OF DATABASE
*/
func main() {
	// This is a placeholder for the main function.
	db, err := bbolt.Open(".db", 0o666, nil)
	if err != nil {
		log.Fatal(err)
	}

	// support int, string, []byte, float, ...
	// temp data
	user := map[string]string{
		"name":  "John Doe",
		"age":   "30",
		"email": "john.doe@example.com",
	}

	// tx is transaction
	db.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucket([]byte("users"))
		if err != nil {
			return err
		}

		id := uuid.New()

		for key, value := range user {
			if err = bucket.Put([]byte(key), []byte(value)); err != nil {
				return err
			}
		}

		if err := bucket.Put([]byte("id"), []byte(id.String())); err != nil {
			return err
		}

		return nil
	})

	userData := make(map[string]string)

	if err := db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte("users"))

		if bucket == nil {
			return fmt.Errorf("Bucket (%s) not found", "users")
		}

		bucket.ForEach(func(key, value []byte) error {
			userData[string(key)] = string(value)
			return nil
		})

		return nil
	}); err != nil {
		log.Fatal(err)
	}

	fmt.Println("User Information:", userData)
}
