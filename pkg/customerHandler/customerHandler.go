package customerhandler

import (
	"errors"
	"fmt"
	"slices"
	"work-mini-project/pkg/configuration"
	filehandler "work-mini-project/pkg/fileHandler"
)

type Customer struct {
	Name  string `json:"name"`
	GridX int    `json:"gridX"`
	GridY int    `json:"gridY"`
}

type CustomerList struct {
	Customers []Customer `json:"customers"`
}

type CustomerHandler struct {
	config    *configuration.Config
	Customers []Customer
}

func wrapError(err error) error {
	return fmt.Errorf("customerHandler: %w", err)
}

var errCustomerNotFound = errors.New("specified customer was not found")

var errCustomerAlreadyExists = errors.New("a customer with that username already exists")

func New(config *configuration.Config) (*CustomerHandler, error) {
	// Parse customers on initialisation
	customers, err := filehandler.ReadFile[CustomerList](config.Customers.FilePath)
	if err != nil {
		return nil, wrapError(err)
	}

	return &CustomerHandler{
		config:    config,
		Customers: customers.Customers,
	}, nil
}

func (ch *CustomerHandler) GetCustomer(name string) (*Customer, error) {
	customerIdx := slices.IndexFunc(ch.Customers, func(E Customer) bool {
		return E.Name == name
	})
	if customerIdx == -1 {
		return nil, wrapError(errCustomerNotFound)
	}

	return &ch.Customers[customerIdx], nil
}

func (ch *CustomerHandler) AddCustomer(customer Customer) error {
	// Check customer name is unique
	customerIdx := slices.IndexFunc(ch.Customers, func(E Customer) bool {
		return E.Name == customer.Name
	})
	if customerIdx != -1 {
		return wrapError(errCustomerAlreadyExists)
	}

	// Update stored customer list
	ch.Customers = append(ch.Customers, customer)

	// Update persistent customer store
	err := filehandler.WriteFile(ch.config.Customers.FilePath, CustomerList{Customers: ch.Customers})
	if err != nil {
		return wrapError(err)
	}

	return nil
}

func (ch *CustomerHandler) RemoveCustomer(customer Customer) error {
	// Find index of customer in stored list
	index := slices.Index(ch.Customers, customer)
	if index == -1 {
		return wrapError(errCustomerNotFound)
	}

	// Crop the customer out of the stored customer list
	ch.Customers = append(ch.Customers[:index], ch.Customers[index+1:]...)

	// Update persistent customer store
	err := filehandler.WriteFile(ch.config.Customers.FilePath, CustomerList{Customers: ch.Customers})
	if err != nil {
		return wrapError(err)
	}

	return nil
}
