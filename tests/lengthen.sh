#!/bin/bash

DOMAIN=http://localhost:3000
ROUTE=api/v1

curl ${DOMAIN}/${ROUTE}/lengthen/${1}
