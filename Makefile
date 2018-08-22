# Merely compiles everything
all: authServer httpServer proxyServer elbServer

# Merely Removes all executables, external libraries except sources and configs
clean:
	rm -rf ./src/github.com ./src/golang.org ./src/gopkg.in ./pkg ;\
	rm -rf ./bin/auth_server ./bin/http_server ./bin/elb ./bin/proxy ;\
	rm -rf ./benchmark/log/*
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
benchmark: install backend setup
	sh ./benchmark/boot.sh benchmark 
	sh ./benchmark/request.sh benchmark
	sh ./benchmark/boot.sh after_bench
concurrent: 	# highly concurrent request benchmark
	sh ./benchmark/request.sh concurrent
image: 		# image type load req benchmark
	sh ./benchmark/request.sh image
large: 		#large load req benchmark
	sh ./benchmark/request.sh large
regular: 	# regular request benchmark
	sh ./benchmark/request.sh regular 
signin: 	# authentication benchmark 
	sh ./benchmark/request.sh signin
region:  	# Regional routing benchmark
	sh ./benchmark/request.sh regional

# structure of auth and http server should comply with structure configs in src/resource
server:  backend
	sh ./benchmark/boot.sh server 	
# structure of proxy and elb should comply with structure configs in src/resource
proxy: install 
	sh ./benchmark/boot.sh proxy 
backend: authServer httpServer 
authServer: 
	go get github.com/lib/pq ;\
	go get golang.org/x/crypto/bcrypt ;\
	go install ./src/auth_server
httpServer:  
	go install ./src/http_server
setup:
	sh ./benchmark/setup.sh
