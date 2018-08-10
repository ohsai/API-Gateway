#!/bin/sh
goroot=$GOPATH
proxy_name=proxy 
proxy=/bin/$proxy_name
resource=/src/resource 

$goroot$proxy 6000 $goroot$resource/msa.yaml &

while true; do
        read -p "Kill proxy?[y/n]" yn
        case $yn in
                [Yy]* ) killall $proxy_name & break;;
                * ) echo "Keep Running Proxy.";;
        esac
done
