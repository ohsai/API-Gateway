#!/bin/sh
goroot=$GOPATH
proxy=/bin/proxy
resource=/src/resource 

$goroot$proxy 6000 $goroot$resource/msa.yaml 

