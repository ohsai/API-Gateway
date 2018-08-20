#!/bin/sh
goroot=$GOPATH
proxy_name=proxy 
proxy=/bin/$proxy_name
resource=/src/resource  
elb_path=/bin/elb 
elb_name=elb 


$goroot$elb_path 9000 $goroot$resource/r_kr.yaml &
$goroot$elb_path 9001 $goroot$resource/r_en.yaml &
$goroot$proxy 6000 $goroot$resource/msa.yaml $goroot$resource/rs_kr.yaml $goroot$resource/config.json &
$goroot$proxy 6001 $goroot$resource/msa.yaml $goroot$resource/rs_kr.yaml $goroot$resource/config.json &
$goroot$proxy 6002 $goroot$resource/msa.yaml $goroot$resource/rs_en.yaml $goroot$resource/config.json &
$goroot$proxy 6003 $goroot$resource/msa.yaml $goroot$resource/rs_en.yaml $goroot$resource/config.json &
ulimit -n 65536
for elbpid in $(pgrep elb) ; do 
        prlimit --nofile=10000 --pid=$elbpid &&
        prlimit --nofile --output RESOURCE,SOFT,HARD --pid $elbpid 

done 
for pxypid in $(pgrep proxy) ; do 
        prlimit --nofile=9000 --pid=$pxypid 
        prlimit --nofile --output RESOURCE,SOFT,HARD --pid $pxypid  
done 

while true; do
        read -p "Kill proxy?[y/n]" yn
        case $yn in
                [Yy]* ) killall $proxy_name & killall $elb_name &  break;;
                * ) echo "Keep Running Proxy.";;
        esac
done
