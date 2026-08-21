#!/bin/bash
go tool pprof -http :9999  http://127.0.0.1:21004/debug/pprof/profile?seconds=30
