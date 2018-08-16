#curl --header "AuthToken:$(http POST localhost:6000/auth/signin username=XXXtentacion password=makeouthill666)" http://localhost:9000/video/test.jpg
ab -c 2000 -n 40000 -H "AuthToken: $(http POST localhost:6000/auth/signin username=XXXtentacion password=makeouthill666)" http://localhost:9000/video/test.css
