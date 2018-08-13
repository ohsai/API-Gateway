curl --header "AuthToken:$(http POST localhost:6000/auth/signin username=XXXtentacion password=makeouthill666)" --header "Accept-Language: en_US" http://localhost:6000/video/test.css  
ab -c 500 -n 30000 -H "AuthToken: $(http POST localhost:6000/auth/signin username=XXXtentacion password=makeouthill666)" -H "Accept-Language: en_US" http://localhost:6000/video/index.html 
