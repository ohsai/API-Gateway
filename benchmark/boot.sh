#!/bin/sh

arg=$1

goroot=$GOPATH

http_server_name=http_server 
http=/bin/$http_server_name

auth_server_name=auth_server 
Auth=/bin/$auth_server_name 

proxy_name=proxy 
proxy=/bin/$proxy_name

resource=/benchmark/resource  

elb_path=/bin/elb 
elb_name=elb 

log_path=/benchmark/log
proxy_log_path=$goroot$log_path/proxy.log 
elb_log_path=$goroot$log_path/elb.log 
server_log_path=$goroot$log_path/server.log

echo $arg 

if [ $arg = "server" ] || [ $arg = "benchmark" ] ; then 
        {
                $goroot$http 5000 &
        $goroot$http 5001 &
        $goroot$http 5002 &
        $goroot$http 5003 &
        $goroot$http 5004 &
        $goroot$http 5005 &
        $goroot$http 5006 &
        $goroot$Auth 8000 &
        $goroot$Auth 8001 &
        $goroot$Auth 8002 &
        $goroot$Auth 8003 &
        } > $server_log_path 2>&1
fi 
if [ $arg = 'proxy' ] || [ $arg = 'benchmark' ] ; then
        ulimit -n 65536
        {
                $goroot$elb_path 9000 $goroot$resource/r_kr.yaml &
                $goroot$elb_path 9001 $goroot$resource/r_en.yaml &
        } > $elb_log_path 2>&1
        echo "elb fd limit modify"
        for elbpid in $(pgrep elb) ; do 
                prlimit --nofile=10000 --pid=$elbpid &&
                prlimit --nofile --output RESOURCE,SOFT,HARD --pid $elbpid 
        done 
        { 
                $goroot$proxy 6000 $goroot$resource/msa.yaml $goroot$resource/rs_kr.yaml $goroot$resource/config.json &
        $goroot$proxy 6001 $goroot$resource/msa.yaml $goroot$resource/rs_kr.yaml $goroot$resource/config.json &
        $goroot$proxy 6002 $goroot$resource/msa.yaml $goroot$resource/rs_en.yaml $goroot$resource/config.json &
        $goroot$proxy 6003 $goroot$resource/msa.yaml $goroot$resource/rs_en.yaml $goroot$resource/config.json & 
} > $proxy_log_path 2>&1
        echo "proxy fd limit modify"
        for pxypid in $(pgrep proxy) ; do 
                prlimit --nofile=9000 --pid=$pxypid 
                prlimit --nofile --output RESOURCE,SOFT,HARD --pid $pxypid  
        done 
fi 

if [ $arg = 'server' ] ; then 
while true; do
        read -p "Kill servers?[y/n]" yn
        case $yn in
                [Yy]* ) killall $http_server_name & killall $auth_server_name &  break;;
                * ) echo "Keep Running Servers.";;
        esac
done
fi 

if [ $arg = 'proxy' ]; then 
while true; do
        read -p "Kill proxy?[y/n]" yn
        case $yn in
                [Yy]* ) killall $proxy_name & killall $elb_name &  break;;
                * ) echo "Keep Running Proxy.";;
        esac
done
fi 

if [ $arg = "after_bench" ]; then
while true; do
        read -p "Kill all?[y/n]" yn
        case $yn in
                [Yy]* ) killall $proxy_name & killall $elb_name & killall $http_server_name & killall $auth_server_name & break;;
                * ) echo "Keep Running all.";;
        esac
done

fi



