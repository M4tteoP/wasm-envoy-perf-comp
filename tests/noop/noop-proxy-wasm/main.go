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

func (*vmContext) NewPluginContext(uint32) types.PluginContext {
	return &noop{}
}

type noop struct {
	types.DefaultPluginContext
}

type httpContext struct {
	types.DefaultHttpContext
}

// Override types.DefaultPluginContext.
func (ctx *noop) OnPluginStart(pluginConfigurationSize int) types.OnPluginStartStatus {
	proxywasm.LogInfo("OnPluginStart called")
	return types.OnPluginStartStatusOK
}

func (*noop) NewHttpContext(uint32) types.HttpContext { return &httpContext{} }

func (ctx *httpContext) OnHttpRequestHeaders(int, bool) types.Action {
	return types.ActionContinue
}
