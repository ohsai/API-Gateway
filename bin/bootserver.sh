#!/bin/sh
goroot=$GOPATH
http=/bin/httpServer 
Auth=/bin/authServer

$goroot$http 5000 &
$goroot$http 5001 &
$goroot$http 5002 &
$goroot$http 5003 &
$goroot$http 5004 &
$goroot$Auth 8000 
