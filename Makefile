# Merely compiles everything
all: authServer httpServer proxyServer elbServer

# Merely Removes all executables, external libraries except sources and configs
clean:
	rm -rf ./src/github.com ./src/golang.org ./src/gopkg.in ./pkg ;\
	rm -rf ./bin/auth_server ./bin/http_server ./bin/elb ./bin/proxy 

#############Install###############

elbServer: dependency
	go install ./src/elb

proxyServer: dependency 
	go get github.com/go-redis/redis ;\
	go get gopkg.in/yaml.v2 ;\
	go install ./src/proxy

install : proxyServer elbServer 

dependency : 
	#sh ./benchmark/dependency.sh 

############BENCHMARKS##############
# regular request benchmark
regular: 
	sh ./benchmark/regquery.sh 

# authentication benchmark 
signin:
	sh ./benchmark/signin.sh

# Regional routing benchmark
region: 
	sh ./benchmark/region.sh 

# Setup several auth and http server for benchmark
# structure of auth and http server should comply with structure configs in src/resource
serverexec: server  
	sh ./benchmark/bootserver.sh 	

# Setup several proxy and elb servers for benchmark
# structure of proxy and elb should comply with structure configs in src/resource
proxyexec: install 
	sh ./benchmark/bootproxy.sh	

server: authServer httpServer 
authServer: 
	go get github.com/lib/pq ;\
	go get golang.org/x/crypto/bcrypt ;\
	go install ./src/auth_server

httpServer:  
	go install ./src/http_server

