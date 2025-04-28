MAKEFLAGS += --silent

# With Middleware

build:
	g++ -c -fPIC Persist-Manager/Persist.cpp Undo-Manager/Single_undo.cpp Undo-Manager/Multiple_undo.cpp -ldl -Wall 
	g++ -shared -o Library/libexample.so Persist.o Single_undo.o Multiple_undo.o -lrt -ldl
	go build -o Library/ Library/ast_rewrite.go
	Library/ast_rewrite < Library/pncounter.go > Library/pncounter_interceptor.go
	Library/ast_rewrite < Library/server.go > Library/ast_server.go
	go build Library/ast_server.go Library/client.go Library/pncounter_interceptor.go Library/file_reader.go Library/wrapper.go

clean:
	@rm -rf ast_server Library/ast_server.go Library/pncounter_interceptor.go Library/ast_rewrite *.json *.o Library/*.so *.so Library/DBs/


# Without Middleware

# build:
# 	go build Library/server.go Library/client.go Library/pncounter.go Library/file_reader.go
# clean:
# 	@rm -rf server