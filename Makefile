##
## * General Targets

## make; - show help
help:
	@make --no-print-directory show-help
show-help:
	@./scripts/fmtMakefileHelp Makefile | column -s ';' -t
	@echo

## make run; - start local backend server
run:
	go run main.go

## make install-deps; - install proto tooling
install-deps:
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install github.com/vadimi/grpc-client-cli/cmd/grpc-client-cli@latest
	go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest

## make protoc; - generate src from proto
protoc:
	protoc --proto_path=pb --go_out=pb --go_opt=paths=source_relative --go-grpc_out=pb --go-grpc_opt=paths=source_relative pb/helloworld.proto
	protoc --proto_path=pb --swagger_out=logtostderr=true:pb pb/helloworld.proto
	protoc -I pb --grpc-gateway_out pb --grpc-gateway_opt logtostderr=true --grpc-gateway_opt paths=source_relative --grpc-gateway_opt generate_unbound_methods=true pb/helloworld.proto

##
## * API Targets

## make grpc-client-cli; - start the interactive grpc-client-cli
grpc-client-cli:
	grpc-client-cli localhost:8080

## make grpc-say-hello; - send grpc request to say hello endpoint
grpc-say-hello:
	grpcurl -d '{"name": "hello"}' -plaintext localhost:8080 helloworld.Greeter/SayHello
## make http-say-hello; - send http request to say hello endpoint
http-say-hello:
	curl -s localhost:8080/say/strval/aeou | jq .
