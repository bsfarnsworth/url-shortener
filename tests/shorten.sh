#!/bin/bash

#DOMAIN=http://localhost:3000
DOMAIN=https://wee.fly.dev
ROUTE=api/v1

curl -X POST -d "fullUrl=${1}" ${DOMAIN}/${ROUTE}/shorten
