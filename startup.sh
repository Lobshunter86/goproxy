#!/usr/bin/bash

./bin/server --addr 0.0.0.0:8443 \
--cacert ./bin/client.crt \
--cert ./bin/server.crt \
--key ./bin/server.key
