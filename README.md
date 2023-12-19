Latency, CPU and memory consumption comparions between proxy-wasm and http-wasm filters in Envoy proxy.

## Benchmarks
Following benchmarks have been run on a [TODO] with:
- Go: `1.20.7`
- Tinygo: `0.30.0`

## Getting started
`go run mage.go -l` lists all the available commands:
```bash
▶ mage -l
Targets:
  build:all*               builds all the wasm tests
  build:httpwasm           builds all the http-wasm tests
  build:proxywasm          builds all the proxy-wasm tests
  run                      runs a test - Usage: run <testName>
  runObservability         spins up Grafana and Prometheus, access grafana at http://localhost:3000.
  teardownObservability    tears down Grafana and Prometheus.

* default target
```

## Setting up the environment
1. Build the filters 
```bash
mage build:all
```
2. Provide Envoy binary under `/envoybins/` directory. Envoy binaries that ran these tests have been built following these steps:
```bash
# proxy-wasm Envoy (upstream)
gh repo clone envoyproxy/envoy
cd envoy
git checkout v1.27.0
./ci/run_envoy_docker.sh './ci/do_ci.sh release.server_only'

mv /tmp/envoy-docker-build/envoy/x64/source/exe/envoy/envoy ~/Repo/wasm-envoy-perf-comp/envoybins/envoy-proxy-wasm

# http-wasm Envoy
gh repo clone vikaschoudhary16/envoy
cd envoy
git checkout d869351dc9a6b4b1badcb986f4de46af734e8057 # For the latest commit checkout http-wasm branch
./ci/run_envoy_docker.sh './ci/do_ci.sh release.server_only'
mv /tmp/envoy-docker-build/envoy/x64/source/exe/envoy/envoy ~/Repo/wasm-envoy-perf-comp/envoybins/envoy-http-wasm # TODO CHECK
```

It is also possible to get them from the docker images:
- **envoy-proxy-wasm**: `envoyproxy/envoy:v1.27.0`
- **envoy-http-wasm**: ... # TODO build and push the image
The envoy binary is located at `/usr/local/bin/envoy` in the docker image.

3. Spin up observability tools. a Grafana dashboard is provided to visualize memory and cpu usage. It requires docker-compose
```bash
mage runObservability # Access the dashboard via: http://localhost:3000. Default login: admin/admin
# When done, tear down the observability tools:
mage teardownObservability
```
4. Run a backend service listening on port 8000
```bash
go run github.com/mccutchen/go-httpbin/v2/cmd/go-httpbin@v2.9.0 -port 8000
```

## Running a test
1. Run `mage run <testName>`.
2. Monitor  the output of the test for latency outputs.
3. Monitor the Grafana dashboard to see the memory and cpu usage of the Envoy process.

## TODOS:
- Add tests:
    - helloworld
    - header manipulation
    - big body not used by the filter
    - body manipulation (we need the whole body)
    - response
- Add CPU with prometheus/node_exporter
