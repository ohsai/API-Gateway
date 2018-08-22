#!/bin/sh

 
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
