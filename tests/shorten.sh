#!/bin/bash

DOMAIN=http://localhost:3000
ROUTE=api/v1

curl -X POST -d "fullUrl=${1}" ${DOMAIN}/${ROUTE}/shorten
