all: authServer httpServer proxyServer

server: authServer httpServer 
	./bin/bootserver.sh 	

proxy: proxyServer
	./bin/bootproxy.sh	

regular: 
	./bin/regquery.sh 

signin:
	./bin/signin.sh

authServer:  
	go install ./src/Auth/authServer.go 

httpServer:  
	go install ./src/http_server/httpServer.go

proxy_src = ./src/proxy
proxyServer: 
	go install $(proxy_src)

