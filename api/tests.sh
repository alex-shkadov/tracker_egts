#!/bin/sh
# FOR DEVELOPMENT
cd /var/www/html
#/usr/bin/go/bin/air index.go 8005
go mod tidy -v
go run tracker/tests/tests.go