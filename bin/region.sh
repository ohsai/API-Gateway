curl --header "AuthToken:$(http POST localhost:9000/auth/signin username=XXXtentacion password=makeouthill666)" --header "Accept-Language: en_US" http://localhost:9000/video/test.css  -v
ab -c 500 -n 30000 -H "AuthToken: $(http POST localhost:9000/auth/signin username=XXXtentacion password=makeouthill666)" -H "Accept-Language: en_US" http://localhost:9000/video/index.html 
