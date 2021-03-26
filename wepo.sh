#!/bin/bash 

URL=$WEPO_URL

if [ -p /dev/stdin ]; then
  body=`cat -`
  # escape LF, tab and double quotes
  body=`echo -n "$body" | \
   perl -pe 's/$/\\\\n/' | perl -pe 's/\R//' | \
   perl -pe 's/\t/\    /g' | \
   perl -pe 's/"/\\\\"/g'`
else
  body="$@"
fi

# echo "$body"
if [ -n "$body" ]; then
  curl -X POST -H 'Content-type: application/json' -d @- $URL << EOS
  {
    "content": "$body"
  }
EOS
else
  echo "empty value"
fi

