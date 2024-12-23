package main

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	"github.com/google/uuid"
	"google.golang.org/api/iterator"
)

type User struct {
	Username string
	Email    string
	Passcode []byte
}

func randomUUID() string {
	return uuid.NewString()
}

func AddUser(client *firestore.Client, userID string, user User) error {
	userLists, err := GetAllUsers(client)
	if err != nil {
		return fmt.Errorf("failed to get userList: %v", err)
	}
	if temp, em := checkForUserDups(userLists, user); temp {
		return fmt.Errorf("%s", em)
	}
	_, err = client.Collection("users").Doc(userID).Set(context.Background(), user)
	if err != nil {
		return fmt.Errorf("failed to add user: %v", err)
	}
	return nil
}

func GetAllUsers(client *firestore.Client) ([]User, error) {
	var users []User
	iter := client.Collection("users").Documents(context.Background())
	defer iter.Stop()

	for {
		doc, err := iter.Next()
		if err != nil {
			if err == iterator.Done {
				return users, nil // Return empty slice if no documents
			}
			return nil, fmt.Errorf("error iterating documents: %v", err)
		}

		var user User
		err = doc.DataTo(&user)
		if err != nil {
			return nil, fmt.Errorf("error parsing user data: %v", err)
		}
		users = append(users, user)
	}
}
