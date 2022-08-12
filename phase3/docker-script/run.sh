#!/bin/sh

echo "Sed nginx conf"
cp -v /.docker-script/etc/nginx/nginx.conf /etc/nginx/
mkdir -p /etc/nginx/conf.d/
cp -v /.docker-script/etc/nginx/conf.d/default.conf /etc/nginx/conf.d/

echo "Starting api server"
pwd
cd /app
/wait && ./summercamp-filestore -e production &

nginx -g 'daemon off;'

