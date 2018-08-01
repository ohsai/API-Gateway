all: authServer httpServer proxyServer

server: authServer httpServer 
	./bin/bootserver.sh 	

proxy: proxyServer
	./bin/bootproxy.sh	

authServer:  
	go install ./src/Auth/authServer.go 

httpServer:  
	go install ./src/http_server/httpServer.go

proxyServer: 
	go install ./src/proxy/proxy.go 

