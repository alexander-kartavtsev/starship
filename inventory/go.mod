module github.com/alexander-kartavtsev/starship/inventory

go 1.24.6

require (
	github.com/alexander-kartavtsev/starship/shared v0.0.0-20251010025859-18c65edb0b36
	github.com/go-faster/errors v0.7.1
	github.com/samber/lo v1.52.0
	github.com/stretchr/testify v1.11.1
	google.golang.org/grpc v1.76.0
)

replace github.com/alexander-kartavtsev/starship/shared => ../shared

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rogpeppe/go-internal v1.13.1 // indirect
	go.opentelemetry.io/otel/sdk/metric v1.38.0 // indirect
	golang.org/x/net v0.46.1-0.20251009175946-9f2f0b95b65d // indirect
	golang.org/x/sys v0.37.0 // indirect
	golang.org/x/text v0.30.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250804133106-a7a43d27e69b // indirect
	google.golang.org/protobuf v1.36.10 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/alexander-kartavtsev/starship/payment => ../payment

replace github.com/alexander-kartavtsev/starship/order => ../order
