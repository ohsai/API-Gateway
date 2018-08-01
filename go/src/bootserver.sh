#!/bin/sh
goroot=/home/ubuntu/Workspace/haskell/PersonalUse/go/src
http=/http_server/httpServer 
Auth=/Auth/authenticate

$goroot$http 5000 &
$goroot$http 5001 &
$goroot$http 5002 &
$goroot$http 5003 &
$goroot$http 5004 &
$goroot$Auth 8000 &
ps -al
