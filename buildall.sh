#!/bin/bash
echo "Building Linux x86"
env GOOS=linux GOARCH=386 go build -o bin/ucspm-linux-x86
echo "Building Linux amd64"
env GOOS=linux GOARCH=amd64 go build -o bin/ucspm-linux-amd64
echo "Building Mac amd64"
env GOOS=darwin GOARCH=amd64 go build -o bin/ucspm-mac
echo "Building Solaris amd64"
env GOOS=solaris GOARCH=amd64 go build -o bin/ucspm-solaris
echo "Building Windows x86"
env GOOS=windows GOARCH=386 go build -o bin/ucspm-windows-x86
echo "Building Windows amd64"
env GOOS=windows GOARCH=amd64 go build -o bin/ucspm-windows-amd64
echo "Building complete"


# android	arm
# darwin	386
# darwin	amd64
# darwin	arm
# darwin	arm64
# dragonfly	amd64
# freebsd	386
# freebsd	amd64
# freebsd	arm
# linux	386
# linux	amd64
# linux	arm
# linux	arm64
# linux	ppc64
# linux	ppc64le
# linux	mips
# linux	mipsle
# linux	mips64
# linux	mips64le
# netbsd	386
# netbsd	amd64
# netbsd	arm
# openbsd	386
# openbsd	amd64
# openbsd	arm
# plan9	386
# plan9	amd64
# solaris	amd64
# windows	386
# windows	amd64