package main

import (
	httpwasm "github.com/http-wasm/http-wasm-guest-tinygo/handler"
	"github.com/http-wasm/http-wasm-guest-tinygo/handler/api"
)

func main() {
	httpwasm.HandleRequestFn = handleRequest
	httpwasm.Host.Log(api.LogLevelInfo, "Main called")
}

func handleRequest(req api.Request, _ api.Response) (next bool, reqCtx uint32) {
	next = true
	return
}
