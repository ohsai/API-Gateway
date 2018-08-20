Done these in 3 weeks
## What I've implemented
* Transparent proxy for Micro Service Architecture 
* yaml format MSA structure
* Authentication / Authorization using HMAC, client token
* Load Balancing : Round robin, ip / url hash, random
* Regional Routing : yaml format MSA structure
* Static Response Cacheing : Redis
* Configuration reloading without restart : json format config
* Load Balancer in front of APIGW
* Error handling
## Where I've implemented
- src/proxy : apigw implementation  
- src/elb : front elb implementation  
- src/http_server/httpServer.go : simple http server  
- src/auth_server : simple authentication server  
******
## HOWTO compile & execute  
### install golang 1.10
```console
sudo add-apt-repository ppa:gophers/archive
sudo apt-get update
sudo apt-get install golang-1.10-go
export GOPATH=$(pwd)
export PATH="/usr/lib/go-1.10/bin:"$(PATH)
```
### clone repo
```console
git clone https://github.com/ohsai/PersonalUse.git
git checkout APIGW
```
### install dependent data store
```console
sudo apt-get install postgresql
netstat -ntlp | grep postgresql
sudo apt-get install redis-server
netstat -ntlp | grep redis
```
### benchmarking
```console
sudo apt-get install httpie
sudo apt-get install apache2-utils
sudo -u postgres createuser ubuntu
sudo -u postgres createdb users
sudo -u postgres psql -c "alter user ubuntu with encrypted password 'ubuntu'"
psql -c users -U ubuntu -c 'create table users( username text \
        primary key, password text,role text );' -h localhost
password : ubuntu # this could be changed
make signup
make server #These two should use one foreground each. How to deal with this?
make proxy 
make signin
make regular
make region
```

