#!/bin/bash

CURRENT=$(git tag -l | tail -n 1)
MAIN=$(cat main.go | grep 'const version = ' | grep -o '[0-9]\.[0-9]\.[0-9]')

