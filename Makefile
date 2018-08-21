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
	sh ./benchmark/dependency.sh 

############BENCHMARKS##############
#benchmark: setup
	#sh ./benchmark/benchmark.sh
#concurrent:
	#sh ./benchmark/request.sh concurrent
#image:
	#sh ./benchmark/request.sh image
#large:
	#sh ./benchmark/request.sh large

# regular request benchmark
regular: 
	sh ./benchmark/request.sh regular 

# authentication benchmark 
signin:
	sh ./benchmark/request.sh signin

# Regional routing benchmark
region: 
	sh ./benchmark/request.sh regional

# Setup several auth and http server for benchmark
# structure of auth and http server should comply with structure configs in src/resource
server:  
	sh ./benchmark/boot.sh server 	

# Setup several proxy and elb servers for benchmark
# structure of proxy and elb should comply with structure configs in src/resource
proxy: install 
	sh ./benchmark/boot.sh proxy 

authServer: 
	go get github.com/lib/pq ;\
	go get golang.org/x/crypto/bcrypt ;\
	go install ./src/auth_server

httpServer:  
	go install ./src/http_server

setup:
	sh ./benchmark/setup.sh
