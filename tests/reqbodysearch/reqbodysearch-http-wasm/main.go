package main

import (
	"fmt"
	"strings"

	httpwasm "github.com/http-wasm/http-wasm-guest-tinygo/handler"
	"github.com/http-wasm/http-wasm-guest-tinygo/handler/api"
)

func main() {
	requiredFeatures := api.FeatureBufferRequest
	if want, have := requiredFeatures, httpwasm.Host.EnableFeatures(requiredFeatures); !have.IsEnabled(want) {
		httpwasm.Host.Log(api.LogLevelError, "Unexpected features, want: "+want.String()+", have: "+have.String())
	}
	httpwasm.HandleRequestFn = handleRequest
	httpwasm.Host.Log(api.LogLevelInfo, "Main called")
}

var headerRcvd bool

func handleRequest(req api.Request, _ api.Response) (next bool, reqCtx uint32) {
	httpwasm.Host.Log(api.LogLevelInfo, "handleRequest() called")
	next = true
	if headerRcvd == false {
		headerRcvd = true
		return
	}

	chunk := make([]byte, 128)
	originalBody := make([]byte, 0, 128)

	for {
		size, eof := req.Body().Read(chunk)
		httpwasm.Host.Log(api.LogLevelInfo, fmt.Sprintf("req.Body().Read(): size %d eof: %v", size, eof))
		originalBody = append(originalBody, chunk[:size]...)
		if eof || size == 0 {
			break
		}
	}

	if strings.Contains(string(originalBody), "payload") {
		httpwasm.Host.Log(api.LogLevelInfo, "pattern found in request body")
	} else {
		httpwasm.Host.Log(api.LogLevelInfo, "pattern not found in request body. Body:"+string(originalBody))
	}
	return
}
