#!/bin/sh
# FOR DEVELOPMENT
cd /var/www/html
go mod tidy -v
/usr/bin/go/bin/air index.go -s true
# порт - опциональный параметр, он также может быть определен в окружении
#go run index.go -s true