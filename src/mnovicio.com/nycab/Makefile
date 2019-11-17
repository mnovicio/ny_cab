NY_CAB_SERVER = ny_cab_server
NY_CAB_CLIENT_REST = ny_cab_client_rest
NY_CAB_CLIENT_GRPC = ny_cab_client_grpc

BIN = ./bin
DOCS = ./doc

head = ${firstword $(subst :, ,${GOPATH})}

all: server

clean:
	rm -rf $(BIN)

clean/proto:
	rm -rf ./protocol/objects/*.go
	rm -rf ./protocol/rpc/swagger
	rm -rf ./protocol/rpc/*.go

protocol: protocol/objects protocol/rpc

protocol/rpc:
	cd protocol/rpc && mkdir -p swagger && protoc -I. -I../../../ \
	-I${head}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
	-I${head}/src/github.com/grpc-ecosystem/grpc-gateway --go_out=plugins=grpc:../../../../ \
	--grpc-gateway_out=logtostderr=true:../../../../ \
	--swagger_out=logtostderr=true:./swagger \
	*.proto

protocol/objects:
	cd protocol/objects && protoc -I. --go_out=plugins=grpc:../../../../ *.proto

server:
	mkdir -p $(BIN)
	cd server/cmd && go build -o ../../$(BIN)/$(NY_CAB_SERVER) -v

client/rest:
	mkdir -p $(BIN)
	cd client/rest && go build -o ../../$(BIN)/$(NY_CAB_CLIENT_REST) -v

client/grpc:
	mkdir -p $(BIN)
	cd client/grpc && go build -o ../../$(BIN)/$(NY_CAB_CLIENT_GRPC) -v


.PHONY: clean protocol protocol/objects protocol/rpc server client/rest client/grpc