#!/usr/bin/env bash

ab -p test.json -T application/json -c 10 -n 10 http://127.0.0.1:8818/