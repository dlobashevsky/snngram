#!/bin/bash

TARGET=http://127.0.0.1:7234

curl -X POST -d '["dmitrilobashevsky","asdfqwerty","mothermary"]' $TARGET/split2
curl -X POST -d '["dmitrilobashevsky","asdfqwerty","mothermary"]' $TARGET/score
