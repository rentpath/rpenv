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
	"strings"
	"syscall"
)

const AppVersion = "3.0.2"
const ConfigPath = ".config/.rpenv"

func main() {
	cmdStatus := 1
	version := flag.Bool("v", false, "Prints rpenv version")
	longVersion := flag.Bool("version", false, "Prints rpenv version")
	flag.Parse()

	if (*version) || (*longVersion) {
		fmt.Println(AppVersion)
		os.Exit(0)
	}

	if flag.NArg() == 0 {
		fmt.Println("must provide an environment, e.g. 'ci', 'qa', or 'prod'...")
	} else {
		envVars := envVars(envUri(flag.Args()[0]))

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
	usr, _ := user.Current()
	conf := getConfig(usr.HomeDir + "/" + ConfigPath, env)

	return conf
}

func envVars(envUri string) []string {
	rawVars := strings.Split(httpRequestBodyAsString(envUri), "\n")
	envsMap := make(map[string]string)
	var keys []string
	formattedVars := make([]string, 0)

	for _, kvPair := range rawVars {
		if !strings.HasPrefix(kvPair, "#") && kvPair != "" {
			kvArray := strings.Split(strings.Replace(kvPair, "\"", "", -1), "=")
			envsMap[kvArray[0]] = kvArray[1]
		}
	}

	for _, kvPair := range os.Environ() {
		kvArray := strings.Split(kvPair, "=")
		envsMap[kvArray[0]] = kvArray[1]
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
		os.Exit(1)
	}

	if env == "production" {
		env = "prod"
	}

	uri := mymap[env]

	if uri == "" {
		fmt.Println("Provided environment must be one of 'ci', 'qa', 'prod', or 'production'.")
		os.Exit(1)
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
	configSlurp := strings.Trim(string(config), " \n")
	for _, line := range strings.Split(configSlurp, "\n") {
		fields := strings.Split(line, "=")
		mymap[fields[0]] = fields[1]
	}

	return mymap
}
