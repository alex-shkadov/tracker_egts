FROM golang:alpine3.14 as debug

RUN apk update && apk add mc bash gcc musl-dev dpkg
ENV GOPATH=/usr/bin/go
RUN go install github.com/go-delve/delve/cmd/dlv@latest
RUN go mod init example/hello
RUN go install github.com/cosmtrek/air@latest
VOLUME ["/var/www/html"]

EXPOSE 2350

WORKDIR /var/www/html
COPY dlv.sh /
RUN chmod +x /dlv.sh

ENTRYPOINT ["/dlv.sh"]

FROM golang:alpine3.14 as devel
RUN apk update && apk add mc

ENV GOPATH=/usr/bin/go

RUN go install github.com/cosmtrek/air@latest
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest && \
    ln -s /usr/bin/go/bin/migrate /usr/local/bin/migrate

#RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/database/postgres@latest
#RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/source/github@latest
RUN go mod init example/hello

VOLUME ["/var/www/html"]

WORKDIR /var/www/html

COPY air.sh /
RUN chmod +x /air.sh
CMD ["/air.sh"]

FROM golang:alpine3.14 as prod
RUN apk update && apk add mc

ENV GOPATH=/usr/bin/go

RUN go install github.com/cosmtrek/air@latest
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest && \
    ln -s /usr/bin/go/bin/migrate /usr/local/bin/migrate

#RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/database/postgres@latest
#RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/source/github@latest
RUN go mod init example/hello

VOLUME ["/var/www/html"]

WORKDIR /var/www/html

COPY prod.sh /
RUN chmod +x /prod.sh
CMD ["/prod.sh"]

