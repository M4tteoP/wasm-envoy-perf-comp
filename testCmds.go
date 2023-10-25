//go:build mage

package main

var urlRequest = "http://localhost:8080/anything"

var testParamsMap = map[string][]string{
	"noop": {"-X", "GET", urlRequest, "-c20", "-n1000000"}, // 20 concurrent requests, 1kk total requests
	// TODO: add the other tests and their params
}
