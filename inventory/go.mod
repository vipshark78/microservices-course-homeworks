module github.com/vipshark78/microservices-course-homeworks/inventory

replace github.com/vipshark78/microservices-course-homeworks/shared => ../shared

go 1.25.2

require (
	github.com/brianvoe/gofakeit/v7 v7.8.0
	github.com/samber/lo v1.52.0
	github.com/stretchr/testify v1.11.1
	github.com/vipshark78/microservices-course-homeworks/shared v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.76.0
	google.golang.org/protobuf v1.36.10
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.27.3 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rogpeppe/go-internal v1.14.1 // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	golang.org/x/net v0.44.0 // indirect
	golang.org/x/sys v0.36.0 // indirect
	golang.org/x/text v0.29.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20251007200510-49b9836ed3ff // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251002232023-7c0ddcbb5797 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
