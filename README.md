# Simple API Gateway  
Written in Golang  
samsung electronics 4 week intern project of Oh Hyun Seok  
  
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
- benchmark/log : save stdout logs of proxy/elb/backend during benchmark
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
# At root folder of this repo
source ./path
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
source ~/.path
make benchmark 
```
`make server` creates backend server structure independently, making it easy to monitor.  
`make proxy`  creates elb and proxy server structure independently, making it easy to monitor.  
`make signin` benchmarks with signin authentication requests.  
`make regular` benchmarks with regular http requests with auth token header.  
`make region` benchmarks with requests for other regional proxy.  
`make concurrent` benchmarks with highly concurrent request load.  
`make image` benchmarks with requests for image file.  
`make large` benchmarks with large size pdf file.  


