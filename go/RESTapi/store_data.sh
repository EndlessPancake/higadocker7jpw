#!/bin/sh

## sample
# curl -i -H 'Content-Type: application/json'  -d '{"ServiceCode":"FXY9009999","IPv4":"10.100.0.1"}' http://127.0.0.1:8080/services

URL="http://127.0.0.1:8080/services"
FILE=test.csv
ECHO=

while IFS=\, read var1 var2 unused; do

  # echo $var1 $var2
  $ECHO curl -i -H 'Content-Type: application/json' -d '{"ServiceCode":"'$var1'","IPv4":"'$var2'"}' $URL

done < $FILE
