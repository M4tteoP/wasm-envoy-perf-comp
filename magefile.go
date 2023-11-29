//go:build mage

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

var Default = Build.All

var testsDir = "./tests"
var proxyWasmSuffix = "proxy-wasm"
var httpWasmSuffix = "http-wasm"

var proxyWasmBinaryPath = "./envoybins/envoy-proxy-wasm"
var httpWasmBinaryPath = "./envoybins/envoy-http-wasm"

type Build mg.Namespace

// All builds all the wasm tests
func (b Build) All() error {
	if err := b.Proxywasm(); err != nil {
		return err
	}
	return b.Httpwasm()
}

// Proxywasm builds all the proxy-wasm tests
func (Build) Proxywasm() error {
	command := exec.Command("tinygo", "build", "-o", "main.wasm", "-scheduler=none", "--no-debug", "-target=wasi")
	err := walkAndBuild(proxyWasmSuffix, command)
	if err != nil {
		fmt.Printf("Error walking directory: %v\n", err)
	}
	return nil
}

// Httpwasm builds all the http-wasm tests
func (Build) Httpwasm() error {
	command := exec.Command("tinygo", "build", "-o", "main.wasm", "-scheduler=none", "--no-debug", "-target=wasi")
	err := walkAndBuild(httpWasmSuffix, command)
	if err != nil {
		fmt.Printf("Error walking directory: %v\n", err)
	}
	return nil
}

func walkAndBuild(suffix string, buildCommand *exec.Cmd) error {
	err := filepath.Walk(testsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("Error accessing path %s: %v\n", path, err)
		}
		if info.IsDir() && strings.HasSuffix(info.Name(), suffix) {
			mainGoPath := filepath.Join(path, "main.go")
			if _, err := os.Stat(mainGoPath); err == nil {
				cmd := exec.Command(buildCommand.Path, buildCommand.Args[1:]...)
				cmd.Dir = filepath.Dir(mainGoPath)

				fmt.Printf("Building %s...\n", mainGoPath)
				err := cmd.Run()
				if err != nil {
					return fmt.Errorf("Error building %s: %v\n", mainGoPath, err)
				}
			}
		}
		return nil
	})
	return err
}

type Test struct {
	name   string
	suffix string
}

// Run runs a test - Usage: run <testName>
func Run(testName string) error {

	// Split the testName into name and suffix
	split := strings.SplitN(testName, "-", 2)
	if len(split) != 2 {
		return fmt.Errorf("invalid test name: %q. Expected name: <testname>[-proxy-wasm|-http-wasm]", testName)
	}
	test := Test{split[0], split[1]}

	// Depending on the suffix, we select the correct envoy binary
	var binaryPath string
	if test.suffix == proxyWasmSuffix {
		binaryPath = proxyWasmBinaryPath
	} else if test.suffix == httpWasmSuffix {
		binaryPath = httpWasmBinaryPath
	} else {
		return fmt.Errorf("invalid test name: %q. Expected name: <testname>[-proxy-wasm|-http-wasm]", testName)
	}

	// Spin up envoy with the correct config file and binary
	configPath := filepath.Join(testsDir, test.name, testName, "envoy.yaml")

	cmd := exec.Command(binaryPath, "-c", configPath, "--service-cluster", "envoy"+test.suffix, "--service-node", "envoy"+test.suffix)
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start envoy: %v", err)
	}

	// Run the requests against Envoy

	if _, ok := testParamsMap[test.name]; !ok {
		return fmt.Errorf("test name %s not found in testParamsMap", test.name)
	}

	testParams := append([]string{"run", "go.k6.io/k6@v0.47.0"}, testParamsMap[test.name]...)
	err := sh.RunV("go", testParams...)
	if err != nil {
		return fmt.Errorf("error running command %q: %v", testParams, err)
	}

	// Kill envoy
	if err := cmd.Process.Kill(); err != nil {
		return fmt.Errorf("failed to kill envoy: %v", err)
	}

	return nil
}

// RunObservability spins up Grafana and Prometheus, access grafana at http://localhost:3000. Requires docker-compose.
func RunObservability() error {
	return sh.RunV("docker-compose", "--file", "./observability/docker-compose.yml", "up", "-d")
}

// TeardownObservability tears down Grafana and Prometheus. Requires docker-compose.
func TeardownObservability() error {
	return sh.RunV("docker-compose", "--file", "./observability/docker-compose.yml", "down")
}
