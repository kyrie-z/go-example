
all: test cgo


test: bcipher.c
	gcc -shared -fPIC -o libbcipher.so bcipher.c

cgo: main.go
	go build .


clean:
	rm -rf cgo  libbcipher.so
