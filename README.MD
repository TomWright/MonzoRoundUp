# Monzo Round Up

Rounds up purchases to the nearest whole number and deposits that remainder to a RoundUp pot.

# Instructions

## Build the binary

```
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
go get github.com/tomwright/monzoroundup
cd $GOPATH/src/github.com/monzoroundup
dep ensure
cd cmd/roundup
go build -i -o $GOPATH/bin/roundup
```

## Start the server

```
roundup aaa bbb mydomain.com:8080
```

Where:
- `aaa` is your oauth client id
- `bbb` is your oauth client secret
- `mydomain.com:8080` is your bind address