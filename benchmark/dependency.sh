#!/bin/sh

#if golang 1.10 not installed, install it

sudo add-apt-repository ppa:gophers/archive
sudo apt-get update
sudo apt-get install golang-1.10-go
export GOPATH=$(pwd)
export PATH="/usr/lib/go-1.10/bin:"$(PATH)

#if redis-server not installed, install it
if [ $(dpkg-query -W -f='${Status}' redis-server 2>/dev/null | grep -c "ok installed") -eq 0 ];
then
        apt-get install redis-server;
fi 

