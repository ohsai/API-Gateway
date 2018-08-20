#!/bin/sh

sudo apt-get install httpie
sudo apt-get install apache2-utils
sudo apt-get install postgresql
netstat -ntlp | grep postgresql
sudo -u postgres createuser ubuntu
sudo -u postgres createdb users
sudo -u postgres psql -c "alter user ubuntu with encrypted password 'ubuntu'"
psql -d users -U ubuntu -c 'create table users( username text \
        primary key, password text,role text );' -h localhost
password : ubuntu 
# this could be changed
make signup
psql -d users -U ubuntu -c "select * from users" -h localhost
password : ubuntu 
# this could be changed
