src/http_server/httpServer.go : simple http server  
src/proxy : apigw implementation  
src/auth_server : simple authentication server  
sec/elb : front elb implementation  
******
# HOWTO Compile & Execute  
## Install golang 1.10
        sudo add-apt-repository ppa:gophers/archive
        sudo apt-get update
        sudo apt-get install golang-1.10-go
        export GOPATH=$(pwd)
        export PATH="/usr/lib/go-1.10/bin:"$(PATH)
## Clone repo
        git clone https://github.com/ohsai/PersonalUse.git
        git checkout APIGW
## Install depencent data store
        sudo apt-get install postgresql
        netstat -ntlp | grep postgresql
        sudo -u postgres createuser ubuntu
        sudo -u postgres createdb users
        sudo -u postgres psql -c "alter user ubuntu with encrypted password 'ubuntu'"
        psql -c users -U ubuntu -c 'create table users( username text primary key, password text,role text );' -h localhost
        password : ubuntu # this could be changed
        sudo apt-get install redis-server
        netstat -ntlp | grep redis
## benchmarking
        sudo apt-get install httpie
        sudo apt-get install apache2-utils
        make signup
        make server #These two should use one foreground each. How to deal with this?
        make proxy 
        make signin
        make regular
        make region

