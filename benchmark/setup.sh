#!/bin/sh


install_if_not()
{
        package_name=$1
if [ $(dpkg-query -W -f='${Status}' $package_name 2>/dev/null | grep -c "ok installed") -eq 0 ];
then
        echo "$package_name install"
        apt-get install -y package_name;
else 
        echo "$package_name already installed"
fi
}

install_if_not httpie 
install_if_not apache2-utils 
install_if_not postgresql 

#This changes user ubuntu (if exists) password into 'ubuntu'
#and drops 'users' table (if exists) of 'users' database
if psql postgres -tAc "SELECT 1 FROM pg_roles WHERE rolname='ubuntu'" | grep -q 0 ;then
        sudo -u postgres createuser ubuntu
fi 
if psql -lqt | cut -d \| -f 1 | grep -qw users |grep 0; then 
       sudo -u postgres createdb users
fi 
sudo -u postgres psql -c "alter user ubuntu with encrypted password 'ubuntu'"
PGPASSWORD=ubuntu psql -d users -U ubuntu -c 'drop table users;' -h localhost 
PGPASSWORD=ubuntu psql -d users -U ubuntu -c 'create table users( username text primary key, password text,role text );' -h localhost
echo "Benchmark authentication data store created"
PGPASSWORD=ubuntu psql -d users -U ubuntu -c "select * from users" -h localhost

