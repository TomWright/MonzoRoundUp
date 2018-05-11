#!/bin/bash

DSN="$MYSQL_USER:$MYSQL_PASS@tcp($MYSQL_HOST:$MYSQL_PORT)/mysql?parseTime=true"
echo $DSN

goose mysql $DSN up