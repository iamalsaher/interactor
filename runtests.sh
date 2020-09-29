#!/bin/bash

gcc test.c -o test
go run interactor.go
rm test