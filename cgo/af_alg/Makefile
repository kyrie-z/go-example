export LD_LIBRARY_PATH=$(shell pwd)

all: test cgo


test: bcipher.c
	gcc -shared -fPIC -o libbcipher.so bcipher.c

cgo: main.go
	go build .
	ldd af_alg


clean:
	rm -rf af_alg  libbcipher.so
