This package contains a GoLang implementation of the Software Discipline Mini Project.
It implements a simple CRM, FileHandler, CustomerHandler, CommandHandler and CLIHandler;
all of which work together to form the CLI application required to view, manage and use the tool to calculate delivery costs and times.

## Pre-requisites

The project requires the [Go 1.22 or higher runtime](https://go.dev/doc/install) to be installed.

## Run the app

```
go run main.go
```

## Test the app

TODO

## Using the app

The app utilises the CLI for user interaction;
populating the screen with information text and a series of numbered options, promoting the user for input.
Based on the inputted number, the associated action will be invoke, or menu opened.

The commands follow a hierarchical structure, which is laid out below.
Where commands are shown as strings and a short description is shown in (brackets), and role requirements are shown in [square brackets], if applicable.

```
Initial load
│
├─── Login (Prompt the user for their username and password)
│   │
│   ├─── Calculate Journey (Calculate the time and costs for a journey to a specific customer)
|   |
│   ├─── Manage Customers [Admin] (Provide customer management tools)
|   |   |
|   |   ├─── Add Customer [Admin] (Prompt the admin for new customer details)
|   |   |
|   |   └─── Remove Customer [Admin] (Remove the selected customer)
|   |
│   └─── Manage Users [Admin] (Provide user management tools)
|       |
|       ├─── Remove User [Admin] (Remove the selected user)
|       |
|       └─── Change User Type [Admin] (Change the type of the selected user; user vs admin)
│
├─── Register (Prompt for new user for a username and password)
│
└─── Help (Display some help text for using the app)
```

### Global commands

The following commands work globally on the majority of input prompts:

- `exit` - exit out of the application
- `logout` - Sign out of the current user, but do not exit the application
- `help` - Display some help text for how to use the application and interact with it
- `cancel` - Abort the current command and go back to a previous menu
