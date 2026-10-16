module github.com/Nokia-Bell-Labs/declarative-agents/application-dsl-spec/magefiles

go 1.26.3

require gopkg.in/yaml.v3 v3.0.1

require github.com/yuin/goldmark v1.8.6

require github.com/Nokia-Bell-Labs/declarative-agents/agent-core v0.0.0 // indirect

replace github.com/Nokia-Bell-Labs/declarative-agents/agent-core v0.0.0 => ../../agent-core
