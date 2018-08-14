curl --header "AuthToken:$(http POST localhost:6000/auth/signin username=XXXtentacion password=makeouthill666)" http://localhost:9000/video/test.css 
ab -c 200 -n 30000 -H "AuthToken: $(http POST localhost:6000/auth/signin username=XXXtentacion password=makeouthill666)" http://localhost:9000/video/test.css 
