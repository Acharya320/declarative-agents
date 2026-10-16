module github.com/Nokia-Bell-Labs/declarative-agents/magefiles

go 1.26.3

require github.com/magefile/mage v1.17.2

require (
	github.com/Nokia-Bell-Labs/declarative-agents/agent-core v0.0.0
	github.com/santhosh-tekuri/jsonschema/v6 v6.0.2
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/yuin/goldmark v1.8.6 // indirect
	golang.org/x/text v0.37.0 // indirect
)

replace github.com/Nokia-Bell-Labs/declarative-agents/agent-core v0.0.0 => ../agent-core
