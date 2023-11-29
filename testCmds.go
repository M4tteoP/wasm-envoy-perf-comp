//go:build mage

package main

var urlRequest = "http://localhost:8080/anything"

var testParamsMap = map[string][]string{
	"nowasm": {"run", "./k6/get-constant-arrival-rate.js", "--env", "URL=" + urlRequest, "--env", "RATE=20000"}, // 20000 RPS
	"noop":   {"run", "./k6/get-constant-arrival-rate.js", "--env", "URL=" + urlRequest, "--env", "RATE=20000"}, // 20000 RPS
	// TODO: add the other tests and their params
}
