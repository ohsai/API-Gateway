#!/bin/sh

#if golang 1.10 not installed, install it
#if golang 1.10 installed but not on path? 
export GOPATH="$(pwd)"
#if go 1.10 does not exist in path
if [ $(echo "$PATH" |grep -q "/usr/lib/go-1.10/bin" | grep 0) ]; then  
        echo "Go 1.10 does not exist in path"
        export PATH="/usr/lib/go-1.10/bin:$(PATH)"
fi 
go_cur_ver="$(go version | head -n1 | cut -d" " -f3)"
go_req_ver="go1.10"
 if [ "$(printf '%s\n' "$go_req_ver" "$go_cur_ver" | sort -V | head -n1)" = "$go_req_ver" ];then 
        echo "Golang version greater than or equal to 1.10"
 else 
        echo "golang 1.10 install"
        sudo add-apt-repository ppa:gophers/archive
        sudo apt-get update
        sudo apt-get install -y golang-1.10-go
 fi 

#if redis-server not installed, install it

install_if_not()
{
        package_name=$1
if [ $(dpkg-query -W -f='${Status}' $package_name 2>/dev/null | grep -c "ok installed") -eq 0 ];
then
        echo "$package_name install"
        sudo apt-get install -y package_name;
else 
        echo "$package_name already installed"
fi
}

install_if_not redis-server
