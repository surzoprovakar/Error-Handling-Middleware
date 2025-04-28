build:
	export GOPATH=$HOME/go
	export PATH=$PATH:$GOPATH/bin
	protoc --go_out=. --go_opt=paths=source_relative SyncMsg/msg.proto
	go mod init map
	go mod tidy
	go build server.go client.go map.go file_reader.go

test:	build
	go test -v map.go map_test.go

clean:	
	@rm -rf server go.mod go.sum SyncMsg/msg.pb.go