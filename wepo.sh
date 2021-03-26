#!/bin/bash 

if [ -z $WEPO_URL ]; then
  echo 'The environment variable "WEPO_URL" is not set.'
  exit 1
fi
URL=$WEPO_URL

if [ -p /dev/stdin ]; then
  body=`cat -`
  # escape LF, tab and double quotes
  body=`echo -n "$body" | \
    perl -CS -pe 's/$/\\\\n/' | perl -CS -pe 's/\R//' | \
    perl -CS -pe 's/\t/\    /g' | \
    perl -CS -pe 's/"/\\\\"/g'`
else
  body="$@"
fi

while getopts d OPT
do
  case $OPT in
    d) echo "$body"
    ;;
  esac
done

if [ -n "$body" ]; then
  curl -X POST -H 'Content-type: application/json' -d @- $URL << EOS
  {
    "content": "$body"
  }
EOS
  echo
else
  echo "empty value"
fi

