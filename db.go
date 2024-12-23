package main

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

type FirebaseConfig struct {
	APIKey            string `json:"apiKey"`
	AuthDomain        string `json:"authDomain"`
	DatabaseURL       string `json:"databaseURL"`
	ProjectID         string `json:"projectId"`
	StorageBucket     string `json:"storageBucket"`
	MessagingSenderID string `json:"messagingSenderId"`
	AppID             string `json:"appId"`
	MeasurementID     string `json:"measurementId"`
}

func InitializeFirebase(credentialsPath string) (*firestore.Client, error) {
	opt := option.WithCredentialsFile(credentialsPath)
	conf := &firebase.Config{
		ProjectID: "gologin-48e92",
	}

	app, err := firebase.NewApp(context.Background(), conf, opt)
	if err != nil {
		return nil, fmt.Errorf("error initializing Firebase: %v", err)
	}

	firestoreClient, err := app.Firestore(context.Background())
	if err != nil {
		return nil, fmt.Errorf("error getting Firestore client: %v", err)
	}

	return firestoreClient, nil
}
func checkForUserDups(users []User, user User) (bool, string) {
	for _, existingUser := range users {
		if user.Email == existingUser.Email {
			return true, "Email already exists"
		}
		if user.Username == existingUser.Username {
			return true, "Username already exists"
		}
	}
	return false, ""
}
