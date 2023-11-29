//go:build mage

package main

var urlRequest = "http://localhost:8080/anything"

var testParamsMap = map[string][]string{
	"noop": {"run", "./k6/get-constant-arrival-rate.js", "--env", "URL=" + urlRequest}, // 10000 RPS - 20/40 concurrent requests
	// TODO: add the other tests and their params
}
