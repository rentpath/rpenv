package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"sort"
	"strconv"
	"strings"
	"syscall"
)

const AppVersion = "3.1.1"
const ConfigPath = ".config/.rpenv"

func main() {
	cmdStatus := 1
	version := flag.Bool("v", false, "Prints rpenv version")
	longVersion := flag.Bool("version", false, "Prints rpenv version")
	var skipLocalEnvs bool
	flag.BoolVar(&skipLocalEnvs, "skip-local", false, "Skips local env vars from output")
	flag.Parse()

	if (*version) || (*longVersion) {
		fmt.Println(AppVersion)
		os.Exit(0)
	}

	if flag.NArg() == 0 {
		fmt.Println("must provide an environment, e.g. 'ci', 'qa', or 'prod'...")
	} else {
		envVars := envVars(envUri(flag.Args()[0]), skipLocalEnvs)

		if flag.NArg() == 1 {
			fmt.Println(strings.Join(envVars, "\n"))
			cmdStatus = 0
		} else {
			cmdStatus = executeCommand(flag.Args(), envVars)
		}
	}

	os.Exit(cmdStatus)
}

func envUri(env string) string {
	usr, err := user.LookupId(strconv.Itoa(os.Getuid()))
	if err != nil {
		fmt.Printf("rpenv: %s\n", err)
		os.Exit(6)
	}

	conf := getConfig(usr.HomeDir + "/" + ConfigPath, env)

	return conf
}

func envVars(envUri string, skipLocalEnvs bool) []string {
	rawVars := strings.Split(httpRequestBodyAsString(envUri), "\n")
	envsMap := make(map[string]string)
	var keys []string
	formattedVars := make([]string, 0)

	for _, kvPair := range rawVars {
		kvPair = strings.TrimSpace(kvPair)
		if !strings.HasPrefix(kvPair, "#") && strings.Contains(kvPair, "=") {
			key, value := splitSimple(strings.Replace(kvPair, "\"", "", -1), "=")
			envsMap[key] = value
		}
	}

	if !skipLocalEnvs {
		for _, kvPair := range os.Environ() {
			key, value := splitSimple(kvPair, "=")
			envsMap[key] = value
		}
	}

	for k := range envsMap {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, key := range keys {
		formattedVars = append(formattedVars, key+"="+envsMap[key])
	}

	return formattedVars
}

func executeCommand(c []string, envs []string) int {
	exitStatus := 0
	cmd := exec.Command(c[1], c[2:]...)
	cmd.Env = envs
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin

	err := cmd.Run()
	if cmd.Process == nil {
		fmt.Fprintf(os.Stderr, "rpenv: %s\n", err)
		exitStatus = 1
	}
	exitStatus = cmd.ProcessState.Sys().(syscall.WaitStatus).ExitStatus()

	return exitStatus
}

func httpRequestBodyAsString(uri string) string {
	resp, err := http.Get(uri)
	defer func() {
		if closeError := resp.Body.Close(); closeError != nil && err == nil {
			fmt.Println("Network error closing connection to URL", uri)
			fmt.Println("Check the URL manually and try again shortly.")
			os.Exit(2)
		}
	}()

	if (err == nil) && (resp.StatusCode == 200) {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Network error reading from URL", uri)
			fmt.Println("Try again shortly.")
			os.Exit(1)
		} else {
			return string(body)
		}
	} else {
		fmt.Println("Unable to reach URL", uri)
		fmt.Println("Please contact infra to remedy this")
		os.Exit(3)
	}
	return "none"
}

func getConfig(configFile string, env string) string {
	mymap := readConfig(configFile)

	if mymap["ci"] == "" || mymap["qa"] == "" || mymap["prod"] == "" {
		fmt.Println("You must have a ~/.config/.rpenv with ci, qa, and prod keys")
		os.Exit(4)
	}

	if env == "production" {
		env = "prod"
	}

	uri := mymap[env]

	if uri == "" {
		fmt.Println("Provided environment must be one of 'ci', 'qa', 'prod', or 'production'.")
		os.Exit(5)
	}

	return uri
}

func readConfig(configFile string) map[string]string {
	mymap := make(map[string]string)
	config, err := ioutil.ReadFile(configFile)
	if err != nil {
		fmt.Printf("You must have a %s file to continue\n", configFile)
		panic(err)
	}
	configSlurp := strings.TrimSpace(string(config))
	for _, line := range strings.Split(configSlurp, "\n") {
		key, value := splitSimple(line, "=")
		mymap[key] = value
	}

	return mymap
}

// only splits on first instance of chr
func splitSimple(str string, chr string) (string, string) {
	split_up := strings.Split(str, "=")
	key := split_up[0]
	value := strings.Join(split_up[1:], "=")
	return key, value
}
