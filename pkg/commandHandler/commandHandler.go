package commandhandler

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
	clihandler "work-mini-project/pkg/cliHandler"
	"work-mini-project/pkg/configuration"
	crmhandler "work-mini-project/pkg/crmHandler"
	customerhandler "work-mini-project/pkg/customerHandler"
	transporthandler "work-mini-project/pkg/transportHandler"

	"github.com/jedib0t/go-pretty/v6/table"
)

type CommandHandler struct {
	config           *configuration.Config
	cliHandler       *clihandler.CLIHandler
	crmHandler       *crmhandler.CRMHandler
	customerHandler  *customerhandler.CustomerHandler
	transportHandler *transporthandler.TransportHandler
}

func errUnrecognisedCommand(command string) error {
	//nolint:err113 // This is a dynamic error, just wrapped in a function
	return fmt.Errorf("error, unrecognised command: %s", command)
}

var errInvalidSelection = errors.New("error, invalid selection, please try again")

var errNoSelfDelete = errors.New("error, cannot delete own user, please try again")

var errNoSelfRoleEdit = errors.New("error, cannot modify own users role, please try again")

var errKeywordEscape = errors.New("") // keyword escape, shouldn't show to user

func wrapError(err error) error {
	return fmt.Errorf("commandHandler: %w", err)
}

const AdminRole = "admin"

func New(
	config *configuration.Config,
	cliHandler *clihandler.CLIHandler,
	crmHandler *crmhandler.CRMHandler,
	customerHandler *customerhandler.CustomerHandler,
	transportHandler *transporthandler.TransportHandler,
) *CommandHandler {
	return &CommandHandler{
		config:           config,
		cliHandler:       cliHandler,
		crmHandler:       crmHandler,
		customerHandler:  customerHandler,
		transportHandler: transportHandler,
	}
}

func (ch *CommandHandler) Handle() error {
	return ch.initialCommands()
}

func (ch *CommandHandler) anyKeyToContinue() {
	//nolint:errcheck // checking this error would expand return type. Don't care about value in this case
	ch.cliHandler.GetUserInput("Press any key to continue...")
}

func (ch *CommandHandler) checkForKeywords(command string) bool {
	switch strings.ToLower(command) {
	case "exit":
		os.Exit(0)

		return true

	case "logout":
		ch.crmHandler.LoggedInUser = nil

		return true

	case "help":
		ch.cliHandler.WriteOutput(helpText)
		ch.anyKeyToContinue()

		return true

	case "cancel":
		return true

	default:
		return false
	}
}

func (ch *CommandHandler) initialCommands() error {
	// If not logged in, entry screen
	if ch.crmHandler.LoggedInUser == nil {
		selection, err := ch.cliHandler.GetUserInput(introText)
		if err != nil {
			return wrapError(err)
		}

		if ch.checkForKeywords(selection) {
			return nil
		}

		switch selection {
		case "1": // Login
			ch.cliHandler.ClearTerminal()
			ch.cliHandler.WriteOutput("Login")

			err = ch.crmHandler.Login()
			if err != nil {
				return wrapError(err)
			}

			ch.cliHandler.ClearTerminal()

			return ch.postLoginCommands()

		case "2": // Register
			ch.cliHandler.ClearTerminal()
			ch.cliHandler.WriteOutput("Register Account")

			return wrapError(ch.crmHandler.RegisterAccount())

		case "3": // Help
			ch.cliHandler.ClearTerminal()
			ch.cliHandler.WriteOutput(helpText)
			ch.anyKeyToContinue()

			return nil

		default:
			ch.cliHandler.ClearTerminal()

			return wrapError(errUnrecognisedCommand(selection))
		}
	}

	// If logged in, post log in commands
	return ch.postLoginCommands()
}

func (ch *CommandHandler) postLoginCommands() error {
	prompt := postLoginText
	if ch.crmHandler.LoggedInUser.Role == AdminRole {
		prompt += adminPostLoginText
	}

	selection, err := ch.cliHandler.GetUserInput(prompt)
	if err != nil {
		return wrapError(err)
	}

	if ch.checkForKeywords(selection) {
		return nil
	}

	switch selection {
	case "1": // Calculate Journey
		return ch.handleCalculateDelivery()

	case "2": // Manage Customers
		if ch.crmHandler.LoggedInUser.Role != AdminRole {
			return nil
		}

		return ch.handleManageCustomers()

	case "3": // Manage Users
		if ch.crmHandler.LoggedInUser.Role != AdminRole {
			return nil
		}

		return ch.handleManageUsers()

	default:
		ch.cliHandler.ClearTerminal()

		return errUnrecognisedCommand(selection)
	}
}

func (ch *CommandHandler) customerSelectMenu() (customerhandler.Customer, error) {
	ch.cliHandler.WriteOutput("Select Customer:\n")

	customerList := ""
	for i, customer := range ch.customerHandler.Customers {
		customerList += fmt.Sprintf("%d - %s\n", i+1, customer.Name)
	}

	selection, err := ch.cliHandler.GetUserInput(customerList)
	if err != nil {
		return customerhandler.Customer{}, wrapError(err)
	}

	if ch.checkForKeywords(selection) {
		return customerhandler.Customer{}, errKeywordEscape
	}

	index, err := strconv.ParseInt(selection, 10, 0)
	if err != nil {
		return customerhandler.Customer{}, wrapError(err)
	}

	if index > int64(len(ch.customerHandler.Customers)) {
		return customerhandler.Customer{}, errInvalidSelection
	}

	customer := ch.customerHandler.Customers[index-1]

	return customer, nil
}

func (ch *CommandHandler) userSelectMenu() (crmhandler.User, error) {
	ch.cliHandler.WriteOutput("Select User (username (role)):\n")

	userList := ""
	for i, user := range ch.crmHandler.Users {
		userList += fmt.Sprintf("%d - %s (%s)\n", i+1, user.Username, user.Role)
	}

	selection, err := ch.cliHandler.GetUserInput(userList)
	if err != nil {
		return crmhandler.User{}, wrapError(err)
	}

	if ch.checkForKeywords(selection) {
		return crmhandler.User{}, errKeywordEscape
	}

	index, err := strconv.ParseInt(selection, 10, 0)
	if err != nil {
		return crmhandler.User{}, wrapError(err)
	}

	if index > int64(len(ch.crmHandler.Users)) {
		return crmhandler.User{}, errInvalidSelection
	}

	user := ch.crmHandler.Users[index-1]

	return user, nil
}

func (ch *CommandHandler) roleSelectMenu() (crmhandler.AccountRole, error) {
	prompt := `Select New Role:
1 - User
2 - Admin`

	selection, err := ch.cliHandler.GetUserInput(prompt)
	if err != nil {
		return crmhandler.USER, wrapError(err)
	}

	if ch.checkForKeywords(selection) {
		return crmhandler.USER, errKeywordEscape
	}

	switch selection {
	case "1":
		return crmhandler.USER, nil

	case "2":
		return crmhandler.ADMIN, nil

	default:
		return crmhandler.USER, errUnrecognisedCommand(selection)
	}
}

func (ch *CommandHandler) handleCalculateDelivery() error {
	ch.cliHandler.ClearTerminal()

	customer, err := ch.customerSelectMenu()
	if err != nil {
		return err
	}

	trips := ch.transportHandler.CalculateCosts(customer)

	methodTable := table.NewWriter()
	methodTable.AppendHeader(table.Row{"Transport Method", "Time Taken", "Cost"})

	for _, trip := range trips {
		methodTable.AppendRow(table.Row{
			trip.Method,
			time.Unix(0, 0).UTC().Add(trip.Duration).Format("15:04:05"),
			fmt.Sprintf("£%.2f", trip.Cost),
		})
	}

	outputMessage := "Costs and durations for all available transport methods: \n\n"

	ch.cliHandler.WriteOutput(outputMessage)
	ch.cliHandler.WriteOutput(methodTable.Render())

	ch.anyKeyToContinue()

	ch.cliHandler.ClearTerminal()

	return nil
}

func (ch *CommandHandler) getCustomerName() (string, error) {
	prompt := "\nPlease provide a customer name:"

	for {
		customerName, err := ch.cliHandler.GetUserInput(prompt)
		if err != nil {
			return "", wrapError(err)
		}

		if ch.checkForKeywords(customerName) {
			return "", errKeywordEscape
		}

		_, err = ch.customerHandler.GetCustomer(customerName)

		// GetCustomer returning nil error means existing customer was found
		if err == nil {
			return customerName, nil
		}

		prompt = "\nCompany already exists, please try again:"
	}
}

func (ch *CommandHandler) getCustomerGridX() (int, error) {
	prompt := "\nPlease provide customers grid X coordinate:"

	for {
		customerGridXString, err := ch.cliHandler.GetUserInput(prompt)
		if err != nil {
			return -1, wrapError(err)
		}

		if ch.checkForKeywords(customerGridXString) {
			return -1, errKeywordEscape
		}

		customerGridX, err := strconv.ParseInt(customerGridXString, 10, 0)
		if err != nil {
			prompt = "\nError parsing value, please provide a single numerical value:"

			continue
		}

		if customerGridX < int64(ch.config.GridLimits.MinX) || customerGridX > int64(ch.config.GridLimits.MaxX) {
			prompt = fmt.Sprintf(
				"\nValue outside acceptable bounds, please provide a single numerical value between %d and %d:",
				ch.config.GridLimits.MinX,
				ch.config.GridLimits.MaxX,
			)

			continue
		}

		return int(customerGridX), nil
	}
}

func (ch *CommandHandler) getCustomerGridY() (int, error) {
	prompt := "\nPlease provide customers grid Y coordinate:"

	for {
		customerGridYString, err := ch.cliHandler.GetUserInput(prompt)
		if err != nil {
			return -1, wrapError(err)
		}

		if ch.checkForKeywords(customerGridYString) {
			return -1, errKeywordEscape
		}

		customerGridY, err := strconv.ParseInt(customerGridYString, 10, 0)
		if err != nil {
			prompt = "\nError parsing value, please provide a single numerical value:"

			continue
		}

		if customerGridY < int64(ch.config.GridLimits.MinY) || customerGridY > int64(ch.config.GridLimits.MaxY) {
			prompt = fmt.Sprintf(
				"\nValue outside acceptable bounds, please provide a single numerical value between %d and %d:",
				ch.config.GridLimits.MinY,
				ch.config.GridLimits.MaxY,
			)

			continue
		}

		return int(customerGridY), nil
	}
}

func (ch *CommandHandler) getCustomerInputs() (customerhandler.Customer, error) {
	customerName, err := ch.getCustomerName()
	if err != nil {
		return customerhandler.Customer{}, err
	}

	customerGridX, err := ch.getCustomerGridX()
	if err != nil {
		return customerhandler.Customer{}, err
	}

	customerGridY, err := ch.getCustomerGridY()
	if err != nil {
		return customerhandler.Customer{}, err
	}

	return customerhandler.Customer{
		Name:  customerName,
		GridX: customerGridX,
		GridY: customerGridY,
	}, nil
}

func (ch *CommandHandler) handleManageCustomers() error {
	ch.cliHandler.ClearTerminal()

	selection, err := ch.cliHandler.GetUserInput(adminCustomerMenu)
	if err != nil {
		return wrapError(err)
	}

	if ch.checkForKeywords(selection) {
		return nil
	}

	switch selection {
	case "1": // Add Customer
		newCustomer, err := ch.getCustomerInputs()
		if err != nil {
			return err
		}

		err = ch.customerHandler.AddCustomer(newCustomer)
		if err != nil {
			return wrapError(err)
		}

		return nil

	case "2": // Remove Customer
		customer, err := ch.customerSelectMenu()
		if err != nil {
			return err
		}

		err = ch.customerHandler.RemoveCustomer(customer)
		if err != nil {
			return wrapError(err)
		}

		return nil

	default:
		return nil
	}
}

//nolint:cyclop // function is still readable
func (ch *CommandHandler) handleManageUsers() error {
	ch.cliHandler.ClearTerminal()

	selection, err := ch.cliHandler.GetUserInput(adminUserMenu)
	if err != nil {
		return wrapError(err)
	}

	if ch.checkForKeywords(selection) {
		return nil
	}

	switch selection {
	case "1": // Remove User
		user, err := ch.userSelectMenu()
		if err != nil {
			return err
		}

		if user.Username == ch.crmHandler.LoggedInUser.Username {
			return wrapError(errNoSelfDelete)
		}

		err = ch.crmHandler.RemoveUser(user)
		if err != nil {
			return wrapError(err)
		}

		return nil

	case "2": // Change User Type
		user, err := ch.userSelectMenu()
		if err != nil {
			return err
		}

		if user.Username == ch.crmHandler.LoggedInUser.Username {
			return wrapError(errNoSelfRoleEdit)
		}

		role, err := ch.roleSelectMenu()
		if err != nil {
			return err
		}

		err = ch.crmHandler.SetUserRole(user, role)
		if err != nil {
			return wrapError(err)
		}

		return nil

	default:
		return nil
	}
}
