#!/bin/sh

arg=$1
concurrency_normal=20
concurrency_test=2000
reqno_normal=5000
reqno_test=1000

get_essential()
{
result=$( \
        echo "$1" \
| grep "Complete\|Failed\|Time per\|Transfer\|longest" \
)
echo "$result"
}

#Signup
if [ $arg = "benchmark" ]; then 
printf "\nSignup beforehand\n"
http -h POST localhost:9000/auth/signup username=Quavo password=stirfry role=mumble | head -n 1
http -h POST localhost:9000/auth/signup username=XXXtentacion password=makeouthill666 role=rapper | head -n 1
http -h POST localhost:9001/auth/signup username=alpha password=beta role=gamma | head -n 1
http -h POST localhost:9001/auth/signup username=Haskell password=curry role=language |  head -n 1
http -h POST localhost:9000/auth/signup username=lilpump password=guccigang role=rapper | head -n 1
fi 
#signin
if [ $arg = "signin" ] || [ $arg = "benchmark" ] ; then 
printf "\nSignin request test\n"
printf " * Right : Call for every signed up id/pw\n"
http -h POST localhost:9000/auth/signin username=Quavo password=stirfry  | head -n 1
http -h POST localhost:9000/auth/signin username=XXXtentacion password=makeouthill666 | head -n 1
http -h POST localhost:9001/auth/signin username=alpha password=beta | head -n 1
http -h POST localhost:9001/auth/signin username=Haskell password=curry | head -n 1
http -h POST localhost:9000/auth/signin username=lilpump password=guccigang | head -n 1
printf " * Wrong : \n"
#wrong id /pw
printf "Wrong id and pw : "
http -h POST localhost:9001/auth/signin username=samsung password=electronics | head -n 1
#wrong pw
printf "Wrong pw : "
http -h POST localhost:9000/auth/signin username=Haskell password=Carey | head -n 1
#wrong id
printf "Wrong id : "
http -h POST localhost:9001/auth/signin username=beta password=beta | head -n 1
printf " * Load : \n "
signin_test_loc=$GOPATH/benchmark/signin_test.txt 
result=$(ab -q -p $signin_test_loc -T application/json -c $concurrency_normal -n $(expr $reqno_normal / 5) localhost:9000/auth/signin)
get_essential "$result"
fi 

#regional
if [ $arg = "regional" ] || [ $arg = "benchmark" ] ; then 
printf "\nRegional request test\n"
printf " * Right :\n"
#without region -> just on here
printf "Call without region code : "
curl -si --header "AuthToken:$(http POST localhost:9000/auth/signin username=XXXtentacion password=makeouthill666)" http://localhost:9000/video/test.css \
       | head -n 1
#kr region call on en
printf "Call for en region code at kr region : " 
curl -si  --header "AuthToken:$(http POST localhost:9001/auth/signin username=XXXtentacion password=makeouthill666)" --header "Accept-Language: ko_KR" http://localhost:9001/video/test.css \
       | head -n 1
#en_us region call on kr
printf "Call for kr region code at en region : " 
curl  -si --header "AuthToken:$(http POST localhost:9000/auth/signin username=XXXtentacion password=makeouthill666)" --header "Accept-Language: en_US" http://localhost:9000/video/test.css \
       | head -n 1
printf "Call for service only in kr : " 
curl  -si --header "AuthToken:$(http POST localhost:9000/auth/signin username=XXXtentacion password=makeouthill666)" --header "Accept-Language: ko_KR" http://localhost:9000/image/test.css \
       | head -n 1
printf "Call for service only in en : " 
curl  -si --header "AuthToken:$(http POST localhost:9001/auth/signin username=XXXtentacion password=makeouthill666)" --header "Accept-Language: en_US" http://localhost:9001/sound/test.css \
       | head -n 1
printf " * Wrong : \n"
#nonexistent region
printf "Call for nonexistent region code : "
curl -si --header "AuthToken:$(http POST localhost:9000/auth/signin username=XXXtentacion password=makeouthill666)" --header "Accept-Language: de_EU" http://localhost:9000/video/test.css \
       | head -n 1
printf "Call for service nonexistent in kr : " 
curl -si --header "AuthToken:$(http POST localhost:9000/auth/signin username=XXXtentacion password=makeouthill666)" --header "Accept-Language: ko_KR" http://localhost:9000/sound/test.css \
       | head -n 1
printf "Call for service nonexistent in en : " 
curl -si --header "AuthToken:$(http POST localhost:9001/auth/signin username=XXXtentacion password=makeouthill666)" --header "Accept-Language: en_US" http://localhost:9001/image/test.css \
       | head -n 1
printf " * Load : \n "
result=$(ab -q -c $concurrency_normal -n $reqno_normal -H "AuthToken: $(http POST localhost:9000/auth/signin username=XXXtentacion password=makeouthill666)" -H "Accept-Language: en_US" http://localhost:9000/video/index.html)
get_essential "$result"
fi 

#regular
if [ $arg = "regular" ] || [ $arg = "benchmark" ]; then 
printf "\nRegular request test\n"
printf " * Right : Call for every service\n"
curl -si --header "AuthToken:$(http POST localhost:9000/auth/signin username=XXXtentacion password=makeouthill666)" http://localhost:9000/video/test.css \
        | head -n 1
curl -si --header "AuthToken:$(http POST localhost:9000/auth/signin username=XXXtentacion password=makeouthill666)" http://localhost:9000/image/index.html \
        | head -n 1
curl -si --header "AuthToken:$(http POST localhost:9001/auth/signin username=XXXtentacion password=makeouthill666)" http://localhost:9001/sound/index.html \
        | head -n 1
#with header
printf "Call with region code header : "
curl -si --header "AuthToken:$(http POST localhost:9000/auth/signin username=XXXtentacion password=makeouthill666)" --header "Accept-Language: ko_KR" http://localhost:9000/video/test.css \
        | head -n 1
printf " * Wrong : \n"
#no resource
printf "Call for nonexistent resource : "
curl -si --header "AuthToken:$(http POST localhost:9000/auth/signin username=XXXtentacion password=makeouthill666)" --header "Accept-Language: ko_KR" http://localhost:9000/image/test.hs \
        | head -n 1
#no resource on other service
printf "Call for nonexistent resource in other service : "
curl -si --header "AuthToken:$(http POST localhost:9000/auth/signin username=XXXtentacion password=makeouthill666)" --header "Accept-Language: ko_KR" http://localhost:9000/video/test.hs \
        | head -n 1
#wrong auth token
printf "Call with wrong auth token : "
curl -si --header "AuthToken:{nice:2222}" --header "Accept-Language: ko_KR" http://localhost:9000/video/test.hs \
        | head -n 1
#no auth header
printf "Call without auth token header : "
curl -si http://localhost:9000/video/test.hs \
        | head -n 1
printf " * Load : \n "
result=$(ab -q -c $concurrency_normal -n $reqno_test -H "AuthToken: $(http POST localhost:9000/auth/signin username=XXXtentacion password=makeouthill666)" http://localhost:9000/video/test.css)
get_essential "$result" #must pass on string form
fi 

#concurrent
if [ $arg = "concurrent" ] || [ $arg = "benchmark" ] ; then 
printf "\nConcurrent request test\n"
ulimit -n 9000
printf " * Load : \n "
result=$(ab -q -c $concurrency_test -n $reqno_normal -H "AuthToken: $(http POST localhost:9000/auth/signin username=XXXtentacion password=makeouthill666)" http://localhost:9000/video/test.css)
get_essential "$result"
fi 

#image
if [ $arg = "image" ] || [ $arg = "benchmark" ] ; then 
printf "\nImage request test\n"
printf " * Load : \n "
result=$(ab -q -c $concurrency_normal -n $reqno_normal -H "AuthToken: $(http POST localhost:6000/auth/signin username=XXXtentacion password=makeouthill666)" http://localhost:9000/video/test.jpg)
get_essential "$result"
fi 

#large
if [ $arg = "large" ] || [ $arg = "benchmark" ]; then
printf "\nLarge request test\n" 
printf " * Load : \n "
result=$(ab -q -c $concurrency_normal -n $reqno_normal -H "AuthToken: $(http POST localhost:9000/auth/signin username=XXXtentacion password=makeouthill666)" http://localhost:9000/video/test.pdf)
get_essential "$result"
fi 
printf "\n"
