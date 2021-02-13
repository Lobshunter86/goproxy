#!/bin/bash

set -e

# call this script with an email address (valid or not).
# like:
# ./makecert.sh demo@random.com
echo "make server cert"
openssl req -new -nodes -x509 -out server.cert -keyout server.key -days 3650 -subj "/C=DE/ST=NRW/L=Earth/O=Random Company/OU=IT/CN=www.random.com/emailAddress=$1"
echo "make client cert"
openssl req -new -nodes -x509 -out client.cert -keyout client.key -days 3650 -subj "/C=DE/ST=NRW/L=Earth/O=Random Company/OU=IT/CN=www.random.com/emailAddress=$1"