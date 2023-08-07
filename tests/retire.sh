#!/bin/bash

ROUTE=api/v1

curl ${DOMAIN}/${ROUTE}/retire/${1}
