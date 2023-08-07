#!/bin/bash

ROUTE=api/v1

curl -X POST -d "fullUrl=${1}" ${DOMAIN}/${ROUTE}/shorten
