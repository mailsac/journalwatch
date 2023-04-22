.PHONY: all clean build-arm64 build-x86

build: build-arm64 build-x86

build-arm64:
	env CGO_ENABLED=1 GOOS=linux GOARCH=arm64 CC=aarch64-linux-gnu-gcc go build -o out/journalwatch-arm64 ./cmd

build-x86:
	env CGO_ENABLED=1 GOOS=linux GOARCH=amd64 CC=x86_64-linux-gnu-gcc go build -o out/journalwatch-x86 ./cmd

clean:
	rm -f out/*

prepare-macos:
	brew install FiloSottile/musl-cross/musl-cross
	ln -s /usr/local/bin/x86_64-linux-musl-gcc /usr/local/bin/x86_64-linux-gnu-gcc
	ln -s /usr/local/bin/aarch64-linux-musl-gcc /usr/local/bin/aarch64-linux-gnu-gcc

prepare-linux:
	sudo apt-get install gcc-aarch64-linux-gnu
	sudo apt-get install gcc-x86-64-linux-gnu
