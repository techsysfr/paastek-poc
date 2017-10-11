package main

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"

	"github.com/DATA-DOG/godog"
)

var output string

func oneEnvironmentVariablesPAASTEKCOREADDRThatPointsToAAddrport(envVar, addr string) error {
	v, ok := os.LookupEnv(envVar)
	if !ok {
		return errors.New("Environment variable undefined")
	}
	_, _, err := net.SplitHostPort(v)
	if err != nil {
		return err
	}
	return nil
}

func anotherVariablePAASTEKCORESCHEMEhttp(envVar, scheme string) error {
	v, ok := os.LookupEnv(envVar)
	if !ok {
		return errors.New("Environment variable undefined")
	}
	if v != scheme {
		return errors.New("Scheme is unexpected")
	}
	return nil
}

func aServiceListeningOnPAASTEKCOREADDR(envVar string) error {
	v, ok := os.LookupEnv(envVar)
	if !ok {
		return errors.New("Environment variable undefined")
	}
	_, err := net.Dial("tcp", v)
	return err
}

func iExecute(arg1 string) error {
	commands := strings.Split(os.ExpandEnv(arg1), " ")
	cmd := exec.Command(commands[0], commands[1:]...)
	cmd.Env = (os.Environ())
	cmd.Stdout = bytes.NewBufferString(output)
	cmd.Stderr = bytes.NewBufferString(output)
	return cmd.Run()
}

func itReturns(arg1 string) error {
	if arg1 != output {
		return fmt.Errorf("Unexpected output %v", output)
	}
	return nil
}

func FeatureContext(s *godog.Suite) {
	s.Step(`^one environment variables "([^"]*)" that matches "([^"]*)"$`, oneEnvironmentVariablesPAASTEKCOREADDRThatPointsToAAddrport)
	s.Step(`^another variable (.*)=(.*)$`, anotherVariablePAASTEKCORESCHEMEhttp)
	s.Step(`^a service listening on (.*)$`, aServiceListeningOnPAASTEKCOREADDR)
	s.Step(`^I execute "([^"]*)"$`, iExecute)
	s.Step(`^it returns "([^"]*)"$`, itReturns)
}
