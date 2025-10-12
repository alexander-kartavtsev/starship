module github.com/alexander-kartavtsev/starship/payment

go 1.24.6

require (
	github.com/alexander-kartavtsev/starship/shared v0.0.0-20251010165449-ece84198bbf7
	github.com/google/uuid v1.6.0
	google.golang.org/grpc v1.76.0
)

replace github.com/alexander-kartavtsev/starship/shared => ../shared

require (
	go.opentelemetry.io/otel/sdk/metric v1.38.0 // indirect
	golang.org/x/net v0.46.1-0.20251009175946-9f2f0b95b65d // indirect
	golang.org/x/sys v0.37.0 // indirect
	golang.org/x/text v0.30.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250804133106-a7a43d27e69b // indirect
	google.golang.org/protobuf v1.36.10 // indirect
)
