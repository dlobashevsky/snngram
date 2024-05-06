#!/bin/bash


TARGET=http://127.0.0.1:7234

curl -X POST -d '["sdtjaxulbhdbstnkekgmsthveskqjqrqgx"]' $TARGET/split2
echo 'high-entropy requests expect to long time processing by deep splitting'
curl -X POST -d '["sdtjaxulbhdbstnkekgmsthveskqjqrqgx"]' $TARGET/score
