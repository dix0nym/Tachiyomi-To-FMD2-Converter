#!/usr/bin/bash

protoc --experimental_allow_proto3_optional --go_opt=paths=source_relative \
  -I ./protos --go_out ./protos protos/tachiyomi.proto
