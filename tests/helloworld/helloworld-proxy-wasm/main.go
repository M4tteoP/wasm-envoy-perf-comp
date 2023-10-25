package main

import (
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

func main() {
	proxywasm.SetVMContext(&vmContext{})
}

type vmContext struct {
	types.DefaultVMContext
}

func (*vmContext) NewPluginContext(contextID uint32) types.PluginContext {
	return &helloWorld{}
}

type helloWorld struct {
	types.DefaultPluginContext
}

// Override types.DefaultPluginContext.
func (ctx *helloWorld) OnPluginStart(pluginConfigurationSize int) types.OnPluginStartStatus {
	proxywasm.LogInfo("OnPluginStart called")
	return types.OnPluginStartStatusOK
}

func (*helloWorld) NewHttpContext(uint32) types.HttpContext { return &types.DefaultHttpContext{} }

func (ctx *helloWorld) OnHttpResponseHeaders(int, bool) types.Action {
	proxywasm.LogInfo("Adding helloworld header")
	if err := proxywasm.AddHttpResponseHeader("x-proxy-wasm-go-sdk-example", "http_headers"); err != nil {
		proxywasm.LogCriticalf("failed to set response constant header: %v", err)
	}
	return types.ActionContinue
}
