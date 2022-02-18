#!/bin/sh
cd /var/www/html
go build -o app
/usr/bin/go/bin/dlv debug --headless --log -l 0.0.0.0:2350 --api-version=2