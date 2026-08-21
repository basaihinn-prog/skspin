#!/bin/bash

protoc -I=. --go_out=.  message.proto