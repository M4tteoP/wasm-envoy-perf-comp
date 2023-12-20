package main

import (
	"strings"

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
	return &pluginContext{}
}

type pluginContext struct {
	types.DefaultPluginContext
}

type httpContext struct {
	types.DefaultHttpContext
	totalRequestBodySize int
}

// Override types.DefaultPluginContext.
func (ctx *pluginContext) OnPluginStart(pluginConfigurationSize int) types.OnPluginStartStatus {
	proxywasm.LogInfo("OnPluginStart called")
	return types.OnPluginStartStatusOK
}

func (*pluginContext) NewHttpContext(uint32) types.HttpContext { return &httpContext{} }

func (ctx *httpContext) OnHttpRequestHeaders(int, bool) types.Action {
	return types.ActionContinue
}

// Override types.DefaultHttpContext.
func (ctx *httpContext) OnHttpRequestBody(bodySize int, endOfStream bool) types.Action {
	ctx.totalRequestBodySize += bodySize
	if !endOfStream {
		// Wait until we see the entire body
		return types.ActionPause
	}
	originalBody, err := proxywasm.GetHttpRequestBody(0, ctx.totalRequestBodySize)
	if err != nil {
		proxywasm.LogErrorf("failed to get request body: %v", err)
		return types.ActionPause
	}
	if strings.Contains(string(originalBody), "payload") {
		proxywasm.LogInfo("pattern found in request body")
	} else {
		proxywasm.LogError("pattern not found in request body:")
	}
	return types.ActionContinue
}
