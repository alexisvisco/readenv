package main

import (
	_ "embed"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

var (
	flagVerbose = false
)

var (
	//go:embed .semver.yaml
	version string
)

var rootCmd = &cobra.Command{
	Use:          "readenv <.env file> <your command>",
	Short:        "Read file as dot env file and execute command with this env.",
	Args:         cobra.MinimumNArgs(2),
	Version:      strings.Split(strings.Split(version, "\n")[3], ":")[1][1:],
	SilenceUsage: true,
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

		env := strings.Split(string(file), "\n")

		shell := os.Getenv("SHELL")
		if shell == "" {
			return errors.New("SHELL environment variable is not set")
		}

		if flagVerbose {
			fmt.Println(shell, "; ", file, "; ", strings.Join(args[1:], " "))
		}

		c := exec.Command(shell, "-c", strings.Join(args[1:], " "))
		c.Env = append(c.Env, env...)

		output, err := c.CombinedOutput()
		if err != nil {
			return fmt.Errorf("%v: command fail", err)
		}

		fmt.Println(string(output))
		return nil
	},
}

func main() {
	rootCmd.Flags().BoolVarP(&flagVerbose, "verbose", "v", false, "add verbosity for debugging")
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
