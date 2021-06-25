package main

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

var rootCmd = &cobra.Command{
	Use:   "readenv <.env file> <your comamnd>",
	Short: "readenv is a tool for reading .env files",
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		dotEnvFile := args[0]
		if f, err := os.Stat(dotEnvFile); os.IsNotExist(err) {
			return fmt.Errorf("%s does not exist", dotEnvFile)
		} else if f.IsDir() {
			return fmt.Errorf("%s must be a file", dotEnvFile)
		}

		file, err := ioutil.ReadFile(dotEnvFile)
		if err != nil {
			return fmt.Errorf("%v: unable to read dot env file", err)
		}

		command := strings.Split(string(file), "\n")[1:]

		shell := os.Getenv("SHELL")
		if shell == "" {
			return errors.New("SHELL environment variable is not set")
		}

		c := exec.Command(shell, "-c", strings.Join(command, " "))

		output, err := c.CombinedOutput()
		if err != nil {
			return fmt.Errorf("%v: command fail", err)
		}

		fmt.Println(string(output))

		return nil
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
