package main

import (
	_ "embed"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"os/exec"
	"os/signal"
	"strings"
)

var (
	//go:embed .semver.yaml
	version string

	flagVerbose = false
)

var rootCmd = &cobra.Command{
	Use:          "readenv <.env file> <your command>",
	Short:        "Read file as dot env file and execute command with this env.",
	Args:         cobra.MinimumNArgs(2),
	Version:      getVersion(),
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
		var envTrimmed []string

		for _, s := range env {
			// only add the env variable if it not start with # or if it is not an empty line
			if strings.Trim(s, " ") != "" && !strings.HasPrefix(s, "#") {
				envTrimmed = append(envTrimmed, s)
			}
		}

		shell := os.Getenv("SHELL")
		if shell == "" {
			return errors.New("SHELL environment variable is not set")
		}

		if flagVerbose {
			fmt.Println("readenv:", shell, "; ", dotEnvFile, "; ", strings.Join(args[1:], " "))
		}

		c := exec.Command(shell, "-c", strings.Join(args[1:], " "))

		c.Env = append(c.Env, os.Environ()...)
		c.Env = append(c.Env, envTrimmed...)

		c.Stderr = os.Stderr
		c.Stdout = os.Stdout
		c.Stdin = os.Stdin

		if err := c.Run(); err != nil {
			fmt.Println(err)
		}

		return nil
	},
}

func getVersion() string {
	return strings.Split(strings.Split(version, "\n")[3], ":")[1][1:]
}

func main() {
	checkVersion()

	exitCode := make(chan int)
	rootCmd.Flags().BoolVarP(&flagVerbose, "verbose", "v", false, "add verbosity for debugging")

	go func() {
		err := rootCmd.Execute()
		if err != nil {
			exitCode <- 1
		}
		exitCode <- 0
	}()

	end := make(chan os.Signal, 1)
	signal.Notify(end, os.Interrupt, os.Kill)
	select {
	case <-end:
		os.Exit(<-exitCode)
	case code := <-exitCode:
		os.Exit(code)
	}
}
