FROM golang:1.10.2

RUN mkdir -p /go/bin

RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

RUN mkdir -p /go/src/github.com/tomwright/monzoroundup

WORKDIR /go/src/github.com/tomwright/monzoroundup
COPY . .

RUN dep ensure -v

WORKDIR /go/src/github.com/tomwright/monzoroundup/cmd/roundup
RUN go build -i -o /go/bin/roundup

WORKDIR /go/src/github.com/tomwright/monzoroundup/cmd/goose
RUN go build -i -o /go/bin/goose

COPY ./docker/app/migrate_up.sh /usr/local/bin/migrate_up.sh
COPY ./docker/app/migrate_down.sh /usr/local/bin/migrate_down.sh
RUN chmod +x /usr/local/bin/migrate_up.sh
RUN chmod +x /usr/local/bin/migrate_down.sh

EXPOSE 8000

WORKDIR /go

CMD ["/go/bin/roundup"]
