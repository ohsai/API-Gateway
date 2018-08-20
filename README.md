Done these in 3 weeks
## What I've implemented
* Transparent proxy for Micro Service Architecture : yaml format MSA structure
* Authentication / Authorization : HMAC, client token
* Load Balancing : Round robin, ip / url hash, random
* Regional Routing : yaml format MSA structure
* Static Response Cacheing : Redis
* Configuration reloading without restart : json format config
* Load Balancer in front of APIGW
* Error handling
## Where I've implemented
- src/proxy : apigw 
- src/elb : front elb  
- src/http_server/httpServer.go : naive http server  
- src/auth_server : naive authentication server  
- src/resource : MSA / Regional elb / Regional proxy structure, Configuration data
- benchmark : benchmark script
******
## HOWTO Setup & execute  
### Setup
#### clone repo
```console
git clone https://github.com/ohsai/PersonalUse.git
git checkout APIGW 
```
#### install
```console
#benchmark/dependency.sh
make install
```
### execute
```console
elb [PORT] [PATH to Regional proxy structure] 
proxy [PORT] [PATH to MSA structure] [PATH to Regional elb structure] [PATH to Configuration]

```
### benchmarking
#### Setup
```console
make benchmark 
#benchmark/setup.sh
```
command `make serverexec` creates backend server structure independently, making it easy to monitor.  
command `make proxyexec`  creates elb and proxy server structure independently, making it easy to monitor.  
command `make signin` benchmarks with a load of authentication requests.  
command `make regular` benchmarks with a load of regular http requests with auth token header.  
command `make region` benchmarks with a load of requests for other regional proxy.  


