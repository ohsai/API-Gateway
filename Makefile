# Merely compiles everything
all: authServer httpServer proxyServer elbServer

# Merely Removes all executables, external libraries except sources and configs
clean:
	rm -rf ./src/github.com ./src/golang.org ./src/gopkg.in ./pkg ;\
	rm -rf ./bin/auth_server ./bin/http_server ./bin/elb ./bin/proxy 


############BENCHMARKS##############
# NEED TO EXECUTE 'make server', 'make proxy' BEFORE BENCHMARKING
# Regular request benchmark
# First requests for authentication, and regular requests using the auth token 
regular: 
	./benchmark/regquery.sh 

# authentication benchmark 
signin:
	./benchmark/signin.sh

# Regional routing benchmark
region: 
	./benchmark/region.sh 

# Setup several auth and http server for benchmark
# structure of auth and http server should comply with structure configs in src/resource
server: authServer httpServer 
	./benchmark/bootserver.sh 	

# Setup several proxy and elb servers for benchmark
# structure of proxy and elb should comply with structure configs in src/resource
proxy: proxyServer elbServer
	./benchmark/bootproxy.sh	
#############COMPILE###############
# Examples of auth server and http server for benchmarking
authServer: 
	go get github.com/lib/pq ;\
	go get golang.org/x/crypto/bcrypt ;\
	go install ./src/auth_server

httpServer:  
	go install ./src/http_server


elbServer:  
	go install ./src/elb

proxyServer:
	go get github.com/go-redis/redis ;\
	go get gopkg.in/yaml.v2 ;\
	go install ./src/proxy

