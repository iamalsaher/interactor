#!/bin/bash

gcc test.c -o binary
go run interactor.go
rm binary