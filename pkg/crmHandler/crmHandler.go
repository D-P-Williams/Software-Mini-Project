package crmhandler

import (
	"errors"
	"fmt"
	"slices"
	clihandler "work-mini-project/pkg/cliHandler"
	"work-mini-project/pkg/configuration"
	filehandler "work-mini-project/pkg/fileHandler"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Username     string `json:"username"`
	PasswordHash string `json:"passwordHash"`
	Role         string `json:"role"`
}

type UsersList struct {
	Users []User `json:"users"`
}

type CRMHandler struct {
	config       *configuration.Config
	Users        []User
	cliHandler   *clihandler.CLIHandler
	LoggedInUser *User
}

func wrapError(err error) error {
	return fmt.Errorf("crmHandler: %w", err)
}

const (
	registrationUsernamePrompt = `
Please provide a username:`

	//nolint:gosec // password as string, not potential password
	registrationPasswordPrompt = `
Please provide a password (1-50 characters):`

	loginUsernamePrompt = `
Enter username:`

	loginPasswordPrompt = `
Enter password:`
)

var errUserNotFound = errors.New("specified user was not found")

var errUserAlreadyExists = errors.New("a user with that username already exists")

var errIncorrectCredentials = errors.New("error, username and password did not match. please try again")

type AccountRole string

const (
	USER  AccountRole = "user"
	ADMIN AccountRole = "admin"
)

func New(config *configuration.Config, cliHandler *clihandler.CLIHandler) (*CRMHandler, error) {
	// Parse customers on initialisation
	users, err := filehandler.ReadFile[UsersList](config.Users.FilePath)
	if err != nil {
		return nil, wrapError(err)
	}

	return &CRMHandler{
		config:       config,
		Users:        users.Users,
		cliHandler:   cliHandler,
		LoggedInUser: nil,
	}, nil
}

func (crm *CRMHandler) Login() error {
	username, err := crm.cliHandler.GetUserInput(loginUsernamePrompt)
	if err != nil {
		return wrapError(err)
	}

	password, err := crm.cliHandler.GetSensitiveInput(loginPasswordPrompt)
	if err != nil {
		return wrapError(err)
	}

	user, err := crm.GetUser(username)
	if err != nil {
		// Local call so don't need to re-wrap
		return err
	}

	passwordValid := verifyPassword(password, user.PasswordHash)
	if !passwordValid {
		return errIncorrectCredentials
	}

	crm.LoggedInUser = &user

	crm.cliHandler.WriteOutput("Successfully logged in!")

	return nil
}

func (crm *CRMHandler) RegisterAccount() error {
	username, err := crm.cliHandler.GetUserInput(registrationUsernamePrompt)
	if err != nil {
		return wrapError(err)
	}

	password, err := crm.cliHandler.GetSensitiveInput(registrationPasswordPrompt)
	if err != nil {
		return wrapError(err)
	}

	passwordHash, err := hashPassword(password)
	if err != nil {
		return wrapError(err)
	}

	err = crm.AddUser(User{Username: username, PasswordHash: passwordHash, Role: "user"})
	if err != nil {
		return wrapError(err)
	}

	crm.cliHandler.WriteOutput("Successfully created new account, please login to continue")

	//nolint:err113 // Return empty error to restart app to prompt for log in
	return errors.New("")
}

func (crm *CRMHandler) GetUser(username string) (User, error) {
	userIdx := slices.IndexFunc(crm.Users, func(E User) bool {
		return E.Username == username
	})

	if userIdx == -1 {
		return User{}, wrapError(errUserNotFound)
	}

	return crm.Users[userIdx], nil
}

func (crm *CRMHandler) AddUser(user User) error {
	// Check username is unique
	usernameMatchIdx := slices.IndexFunc(crm.Users, func(E User) bool {
		return E.Username == user.Username
	})
	if usernameMatchIdx != -1 {
		return wrapError(errUserAlreadyExists)
	}

	// Update stored users list
	crm.Users = append(crm.Users, user)

	// Update persistent customer store
	err := filehandler.WriteFile(crm.config.Users.FilePath, UsersList{Users: crm.Users})
	if err != nil {
		return wrapError(err)
	}

	return nil
}

func (crm *CRMHandler) RemoveUser(user User) error {
	// Find index of user in stored list
	index := slices.Index(crm.Users, user)
	if index == -1 {
		return wrapError(errUserNotFound)
	}

	// Crop the user out of the stored users list
	crm.Users = append(crm.Users[:index], crm.Users[index+1:]...)

	// Update persistent users store
	err := filehandler.WriteFile(crm.config.Users.FilePath, UsersList{Users: crm.Users})
	if err != nil {
		return wrapError(err)
	}

	return nil
}

func (crm *CRMHandler) SetUserRole(user User, role AccountRole) error {
	// Find index of user in stored list
	index := slices.Index(crm.Users, user)
	if index == -1 {
		return wrapError(errUserNotFound)
	}

	// Update stored users list
	crm.Users[index].Role = string(role)

	// Update persistent users store
	err := filehandler.WriteFile(crm.config.Users.FilePath, UsersList{Users: crm.Users})
	if err != nil {
		return wrapError(err)
	}

	return nil
}

// HashPassword generates a bcrypt hash for the given password.
func hashPassword(password string) (string, error) {
	// If changed, passwords will invalidate
	passwordCost := 4

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), passwordCost)

	return string(bytes), err
}

// VerifyPassword verifies if the given password matches the stored hash.
func verifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	return err == nil
}
