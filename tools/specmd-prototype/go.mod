module github.com/Nokia-Bell-Labs/declarative-agents/tools/specmd-prototype

go 1.26.3

require (
	github.com/Nokia-Bell-Labs/declarative-agents/agent-core v0.0.0
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/dominikbraun/graph v0.23.0 // indirect
	github.com/yuin/goldmark v1.8.6 // indirect
	go.opentelemetry.io/otel v1.44.0 // indirect
	go.opentelemetry.io/otel/metric v1.44.0 // indirect
	go.opentelemetry.io/otel/trace v1.44.0 // indirect
)

replace (
	github.com/Nokia-Bell-Labs/declarative-agents/agent-core v0.0.0 => ../../agent-core
	github.com/Nokia-Bell-Labs/declarative-agents/magefiles v0.0.0-20260917222628-140b8b5c3b12 => ../../magefiles
)
