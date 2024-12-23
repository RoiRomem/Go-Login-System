package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

func main() {
	var input string
	var emailInput, nameInput string
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Welcome to the GoLogin system!")
	credentialsPath := "serviceAccount.json"

	firestoreClient, err := InitializeFirebase(credentialsPath)
	if err != nil {
		fmt.Printf("Firebase initialization error: %v\n", err)
		return
	}
	defer firestoreClient.Close()

	for {
		fmt.Print("Do you want to login or sign up (l-login, s-sign up): ")
		input, _ = reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if input == "l" || input == "s" {
			break
		}
		fmt.Println("Invalid input, try again")
	}

	if input == "s" {
		for {
			fmt.Print("Enter a username: ")
			nameInput, _ = reader.ReadString('\n')
			nameInput = strings.TrimSpace(nameInput)

			fmt.Print("Enter an email: ")
			emailInput, _ = reader.ReadString('\n')
			emailInput = strings.TrimSpace(emailInput)

			fmt.Print("Enter a password: ")
			passInput, _ := reader.ReadString('\n')
			passInput = strings.TrimSpace(passInput)

			if nameInput == "" || emailInput == "" || passInput == "" {
				fmt.Println("All fields are required")
				continue
			}

			if !isValidEmail(emailInput) {
				fmt.Println("Invalid email format")
				continue
			}

			encryptedPasscode, err := encryptPass(passInput)
			if err != nil {
				fmt.Printf("Password encryption failed: %v\n", err)
				return
			}

			newUser := &User{
				Username: nameInput,
				Email:    emailInput,
				Passcode: encryptedPasscode,
			}

			err = AddUser(firestoreClient, randomUUID(), *newUser)
			if err != nil {
				fmt.Printf("Failed to add user: %v\n", err)
				return
			}
			fmt.Println("User registered successfully!")
			break
		}
	} else {
		fmt.Print("Enter your email or username: ")
		emailInput, _ = reader.ReadString('\n')
		emailInput = strings.TrimSpace(emailInput)

		fmt.Print("Enter your password: ")
		passInput, _ := reader.ReadString('\n')
		passInput = strings.TrimSpace(passInput)

		userList, err := GetAllUsers(firestoreClient)
		if err != nil {
			fmt.Printf("Failed to get users: %v\n", err)
			return
		}

		found := false
		for _, user := range userList {
			if user.Username == emailInput || user.Email == emailInput {
				found = true
				if verify(user.Passcode, passInput) {
					fmt.Println("Login successful!")
					return
				}
				fmt.Println("Wrong password")
				return
			}
		}
		if !found {
			fmt.Println("User not found")
		}
	}
}

func isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
	return emailRegex.MatchString(email)
}
