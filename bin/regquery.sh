curl --header "AuthToken:$(http POST localhost:6000/auth/signin username=XXXtentacion password=makeouthill666)" http://localhost:6000/video/test.css 
ab -c 2000 -n 30000 -H "AuthToken: $(http POST localhost:6000/auth/signin username=XXXtentacion password=makeouthill666)" http://localhost:6000/video/test.css 
