package clihandler

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strings"

	"golang.org/x/term"
)

type CLIHandler struct {
	reader *bufio.Reader
}

func wrapError(err error) error {
	return fmt.Errorf("cliHandler: %w", err)
}

const inputPrompt = `

> `

func New() *CLIHandler {
	return &CLIHandler{
		reader: bufio.NewReader(os.Stdin),
	}
}

func (cli *CLIHandler) WriteOutput(output string) {
	fmt.Println(output)
}

func (cli *CLIHandler) GetUserInput(prompt string) (string, error) {
	fmt.Print(prompt, inputPrompt)

	// Read until new line
	text, err := cli.reader.ReadString('\n')
	if err != nil {
		return "", wrapError(err)
	}

	// Trim new line an carriage returns
	input := strings.ReplaceAll(strings.ReplaceAll(text, "\n", ""), "\r", "")

	return input, nil
}

// Get CLI input from the user without echoing out types characters.
// Used to get sensitive information, such as passwords.
func (cli *CLIHandler) GetSensitiveInput(prompt string) (string, error) {
	fmt.Print(prompt, inputPrompt)

	// Get the initial state of the terminal.
	initialTermState, e1 := term.GetState(int(os.Stdin.Fd()))
	if e1 != nil {
		panic(e1)
	}

	// Restore it in the event of an interrupt.
	// CITATION: Konstantin Shaposhnikov - https://groups.google.com/forum/#!topic/golang-nuts/kTVAbtee9UA
	interuptChan := make(chan os.Signal, 1)

	//nolint:staticcheck
	signal.Notify(interuptChan, os.Interrupt, os.Kill)

	go func() {
		<-interuptChan

		_ = term.Restore(int(os.Stdin.Fd()), initialTermState)

		os.Exit(1)
	}()

	// Read until new line
	sensitiveString, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return "", wrapError(err)
	}

	// Stop looking for ^C on the channel.
	signal.Stop(interuptChan)

	return string(sensitiveString), nil
}

func (cli *CLIHandler) ClearTerminal() {
	clearMap := make(map[string]func()) // Initialize it
	clearMap["linux"] = func() {
		cmd := exec.Command("clear") // Linux example, its tested
		cmd.Stdout = os.Stdout
		//nolint:errcheck
		cmd.Run()
	}
	clearMap["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls") // Windows example, its tested
		cmd.Stdout = os.Stdout
		//nolint:errcheck
		cmd.Run()
	}

	value, ok := clearMap[runtime.GOOS] // runtime.GOOS -> linux, windows, darwin etc.
	if ok {                             // if we defined a clear func for that platform:
		value() // we execute it
	} else { // unsupported platform
		panic("Your platform is unsupported! I can't clear terminal screen :(")
	}
}
