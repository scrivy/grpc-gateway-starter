## grpc-gateway-starter

This boilerplate grpc service handles http grpc-gateway traffic on the same port.

```
$ make
=========================  
 * General Targets         
 make                       - show help
 make run                   - start local backend server
 make install-deps          - install proto tooling
 make protoc                - generate src from proto
=========================  
 * API Targets             
 make grpc-client-cli       - start the interactive grpc-client-cli
 make grpc-say-hello        - send grpc request to say hello endpoint
 make http-say-hello        - send http request to say hello endpoint
```
