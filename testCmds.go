//go:build mage

package main

var urlRequest = "http://localhost:8080/anything"

var testParamsMap = map[string][]string{
	"nowasm":        {"run", "./k6/constant-arrival-rate.js", "--env", "URL=" + urlRequest, "--env", "RATE=10000"},     // 10000 RPS
	"noop":          {"run", "./k6/constant-arrival-rate.js", "--env", "URL=" + urlRequest, "--env", "RATE=10000"},     // 10000 RPS
	"reqbodysearch": {"run", "./k6/body-constant-arrival-rate.js", "--env", "URL=" + urlRequest, "--env", "RATE=1000"}, // 1000 RPS
	// TODO: add the other tests and their params
}
