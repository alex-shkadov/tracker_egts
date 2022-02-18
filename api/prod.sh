#!/bin/sh
# FOR DEVELOPMENT
cd /var/www/html
go mod tidy -v
# порт - опциональный параметр, он также может быть определен в окружении
go run index.go #8005