#!/bin/sh

arg=$1
#Signup
if [ $arg = "benchmark" ]; then 
http POST localhost:9000/auth/signup username=Quavo password=stirfry role=mumble
http POST localhost:9000/auth/signup username=XXXtentacion password=makeouthill666 role=rapper
http POST localhost:9001/auth/signup username=alpha password=beta role=gamma
http POST localhost:9001/auth/signup username=Haskell password=curry role=language
http POST localhost:9000/auth/signup username=lilpump password=guccigang role=rapper
fi 
#signin
if [ $arg = "signin" ]; then 
http POST localhost:9000/auth/signin username=Quavo password=stirfry 
http POST localhost:9000/auth/signin username=XXXtentacion password=makeouthill666
http POST localhost:9001/auth/signin username=alpha password=beta
http POST localhost:9001/auth/signin username=Haskell password=curry
http POST localhost:9000/auth/signin username=lilpump password=guccigang 
signin_test_loc=$GOPATH/benchmark/signin_test.txt 
ab -p $signin_test_loc -T application/json -c 20 -n 700 localhost:9000/auth/signin 
fi 

#regional
if [ $arg = "regional" ]; then 
curl --header "AuthToken:$(http POST localhost:9000/auth/signin username=XXXtentacion password=makeouthill666)" --header "Accept-Language: en_US" http://localhost:9000/video/test.css  -v
fi 

#regular
if [ $arg = "regular" ]; then 
#curl --header "AuthToken:$(http POST localhost:6000/auth/signin username=XXXtentacion password=makeouthill666)" http://localhost:9000/video/test.jpg
ab -c 2000 -n 40000 -H "AuthToken: $(http POST localhost:6000/auth/signin username=XXXtentacion password=makeouthill666)" http://localhost:9000/video/test.css
fi 

