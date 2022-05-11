# Makefile for yangdump

all:
	go build

build_linux:
	GOOS=linux GOARCH=386 go build

clean:
	rm yangdump
