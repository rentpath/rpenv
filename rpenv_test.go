package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"testing"
)

var (
	system = flag.Bool("system", false, "Run system tests that talk over the network")
)

func captureStd() (*os.File, *os.File, *os.File, *os.File) {
	oldStdout := os.Stdout
	oldStderr := os.Stderr
	_, writerStdout, _ := os.Pipe()
	_, writerStderr, _ := os.Pipe()
	os.Stdout = writerStdout
	os.Stderr = writerStderr
	return writerStdout, writerStderr, oldStdout, oldStderr
}

func restoreStd(writerStdout *os.File, writerStderr *os.File, stdout *os.File, stderr *os.File) {
	writerStdout.Close()
	writerStderr.Close()
	os.Stdout = stdout
	os.Stderr = stderr
}

func TestMain(m *testing.M) {
	flag.Parse()
	result := m.Run()
	os.Exit(result)
}

func TestEnvUri(t *testing.T) {
	for _, env := range []string{"ci", "qa"} {
		conf := envUri(env)

		if (!strings.HasPrefix(conf, "http") || !strings.Contains(conf, env)) {
			t.Errorf("Flag %q does not produce expected ci or qa uri", conf)
		}
	}

	for _, env := range []string{"prod", "production"} {
		conf := envUri(env)

		if !strings.HasPrefix(conf, "http") {
			t.Errorf("Flag %q does not produce expected production uri", conf)
		}
	}
}

func TestEnvVars(t *testing.T) {
	if !*system {
		t.Skip()
	} else {
		results := envVars(envUri("ci"), false)
		passes := 3
		for _, keyValue := range results {
			splitKeyValue := strings.Split(keyValue, "=")
			key := splitKeyValue[0]
			if key == "HOME" {
				passes--
			}
			if key == "APPLICATIONS_ROOT" {
				passes--
			}
			if key == "PATH" {
				passes--
			}
		}
		if passes > 0 {
			fmt.Fprintln(os.Stderr, "Expected envVars with no skip local to return environment variables, but couldn't find 'HOME', 'APPLICATIONS_ROOT', or 'PATH' vars which should be present.")
			t.Errorf("Environment variables found: %q", results)
		}

		results = envVars(envUri("ci"), true)
		fails := 0
		for _, keyValue := range results {
			splitKeyValue := strings.Split(keyValue, "=")
			key := splitKeyValue[0]
			if key == "PATH" {
				fails++
			}
		}
		if fails > 0 {
			fmt.Fprintln(os.Stderr, "Expected envVars with skip local to return environment variables, but those shouldn't include the 'PATH' var.")
			t.Errorf("Environment variables found: %q", results)
		}
	}
}

func TestExecuteCommand(t *testing.T) {
	writerStdout, writerStderr, stdout, stderr := captureStd()

	exitCode := executeCommand([]string{"ci", "ls"}, []string{})
	if exitCode != 0 {
		t.Errorf("Command 'ls' did not return a 0 exit code. Are you on a *nix machine?")
	}
	cmd := []string{"ci", "ls", "not-here-1234"}
	exitCode = executeCommand(cmd, []string{})
	if exitCode != 1 {
		t.Errorf("Command '%q' returned %q exit code, but expected %q.", cmd, exitCode, 0)
	}

	restoreStd(writerStdout, writerStderr, stdout, stderr)
}
