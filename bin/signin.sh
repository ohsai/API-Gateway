#!/bin/sh
http POST localhost:6000/auth/signin username=Quavo password=stirfry 
http POST localhost:6000/auth/signin username=XXXtentacion password=makeouthill666
http POST localhost:6000/auth/signin username=alpha password=beta
http POST localhost:6000/auth/signin username=Haskell password=curry
http POST localhost:6000/auth/signin username=lilpump password=guccigang

signin_test_loc=$GOPATH/signin_test.txt

ab -p $signin_test_loc -T application/json -c 20 -n 700 localhost:6000/auth/signin 


