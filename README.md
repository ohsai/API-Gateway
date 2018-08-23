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
  
###Architecture  
![apigw-architecture](https://github.com/ohsai/API-Gateway/blob/master/apigw_doc.png)  
  
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
  
#### Environment & Dependency
- Ubuntu 16.04 LTS
- make , git  
(packages installed during setup)
- golang-1.10
- redis-server
- postgresql
- httpie
- apache2-utils
  
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
Benchmark includes my own qualitative tests and apache benchmark load test. Following one command will install everything.  
```console
source ~/.path
make benchmark 
```
`make server` creates backend server structure independently, making it easy to monitor.  
`make proxy`  creates elb and proxy server structure independently, making it easy to monitor.  
  
Test commands below should be executed after benchmark proxies and backend servers are on.  
`make signin` benchmarks with signin authentication requests.  
`make regular` benchmarks with regular http requests with auth token header.  
`make region` benchmarks with requests for other regional proxy.  
`make concurrent` benchmarks with highly concurrent request load.  
`make image` benchmarks with requests for image file.  
`make large` benchmarks with large size pdf file.  

#### Server structure used for benchmark  
  
![benchmarkstructure](https://github.com/ohsai/API-Gateway/blob/master/apigw_doc2.png)  
  
#### Example benchmark result
Benchmark on AWS EC2 c5.2xlarge Ubuntu 16.04 LTS

```console
====== Signup beforehand ======
HTTP/1.1 200 OK
HTTP/1.1 200 OK
HTTP/1.1 200 OK
HTTP/1.1 200 OK
HTTP/1.1 200 OK

====== Signin request test ======
 * Right : Request for every signed up id/pw
HTTP/1.1 200 OK
HTTP/1.1 200 OK
HTTP/1.1 200 OK
HTTP/1.1 200 OK
HTTP/1.1 200 OK
 * Wrong : 
Wrong id and pw : HTTP/1.1 401 Unauthorized
Wrong pw : HTTP/1.1 401 Unauthorized
Wrong id : HTTP/1.1 401 Unauthorized
 * Load : 
 Complete requests:      6000
Failed requests:        0
Time per request:       58.353 [ms] (mean)
Time per request:       2.918 [ms] (mean, across all concurrent requests)
Transfer rate:          69.95 [Kbytes/sec] received
 100%    168 (longest request)

====== Regional request test ======
 * Right :
Request without region code : HTTP/1.1 200 OK
Request for kr region code and service only in kr at en region : HTTP/1.1 200 OK
Request for en region code and service only in en at kr region : HTTP/1.1 200 OK
Request for service only in kr : HTTP/1.1 200 OK
Request for service only in en : HTTP/1.1 200 OK
 * Wrong : 
Request for nonexistent region code : HTTP/1.1 404 Not Found
Request for service nonexistent in kr : HTTP/1.1 404 Not Found
Request for service nonexistent in en : HTTP/1.1 404 Not Found
 * Load : 
 Complete requests:      30000
Failed requests:        0
Time per request:       8.581 [ms] (mean)
Time per request:       0.429 [ms] (mean, across all concurrent requests)
Transfer rate:          1527.34 [Kbytes/sec] received
 100%     27 (longest request)

====== Regular request test =======
 * Right : Request for every service
HTTP/1.1 200 OK
HTTP/1.1 200 OK
HTTP/1.1 200 OK
Request with region code header : HTTP/1.1 200 OK
 * Wrong : 
Request for nonexistent resource : HTTP/1.1 404 Not Found
Request for nonexistent resource in other service : HTTP/1.1 404 Not Found
Request with wrong auth token : HTTP/1.1 400 Bad Request
Request without auth token header : HTTP/1.1 400 Bad Request
 * Load : 
 Complete requests:      100000
Failed requests:        0
Time per request:       8.692 [ms] (mean)
Time per request:       0.435 [ms] (mean, across all concurrent requests)
Transfer rate:          393.24 [Kbytes/sec] received
 100%     30 (longest request)

====== Concurrent request test ======
 * Load : 
 Complete requests:      30000
Failed requests:        0
Time per request:       1084.138 [ms] (mean)
Time per request:       0.542 [ms] (mean, across all concurrent requests)
Transfer rate:          315.27 [Kbytes/sec] received
 100%   8165 (longest request)

====== Image request test ======
 * Load : 
 Complete requests:      30000
Failed requests:        0
Time per request:       9.460 [ms] (mean)
Time per request:       0.473 [ms] (mean, across all concurrent requests)
Transfer rate:          86610.22 [Kbytes/sec] received
 100%     29 (longest request)

====== Large request test ======
 * Load : 
 Complete requests:      30000
Failed requests:        0
Time per request:       17.825 [ms] (mean)
Time per request:       0.891 [ms] (mean, across all concurrent requests)
Transfer rate:          1017095.25 [Kbytes/sec] received
 100%     57 (longest request)
```
