#!/bin/sh
goroot=$GOPATH
http_server_name=httpServer 
auth_server_name=authServer 
http=/bin/$http_server_name
Auth=/bin/$auth_server_name 

$goroot$http 5000 &
$goroot$http 5001 &
$goroot$http 5002 &
$goroot$http 5003 &
$goroot$http 5004 &
$goroot$http 5005 &
$goroot$http 5006 &
$goroot$http 5007 &
$goroot$http 5008 &
$goroot$http 5009 &
$goroot$http 5010 &
$goroot$http 5011 &
$goroot$http 5012 &
$goroot$Auth 8000 &
$goroot$Auth 8001 &
$goroot$Auth 8002 &
$goroot$Auth 8003 &
$goroot$Auth 8004 &
$goroot$Auth 8005 &
$goroot$Auth 8006 &
$goroot$Auth 8007 &

while true; do
        read -p "Kill servers?[y/n]" yn
        case $yn in
                [Yy]* ) killall $http_server_name & killall $auth_server_name &  break;;
                * ) echo "Keep Running Servers.";;
        esac
done
