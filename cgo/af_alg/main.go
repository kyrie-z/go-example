package main

// #cgo CFLAGS: -I./
// #cgo LDFLAGS: -L${SRCDIR}/ -lbcipher
//
// #include "bcipher.h"
// #include <stdlib.h>
import "C"
import (
	"encoding/hex"
	"fmt"
	"unsafe"
)

func main() {
	key := []byte{0x06, 0xa9, 0x21, 0x40, 0x36, 0xb8, 0xa1, 0x5b, 0x51, 0x2e, 0x03, 0xd5, 0x34, 0x12, 0x00, 0x06}
	var h *C.struct_crypt_cipher
	var ckey = C.CBytes(key) //byte slice(GO) convert to void*(C)
	defer C.free(ckey)
	C.crypt_cipher_init((**C.struct_crypt_cipher)(unsafe.Pointer(&h)), (*C.char)(ckey))
	var plain string = "plain test test test"
	var cplain = C.CString(plain)
	defer C.free(unsafe.Pointer(cplain))

	var outbuf = make([]byte, 16)
	C.crypt_cipher_encrypt(h.opfd, cplain, C.ulong(16), (*C.char)(unsafe.Pointer(&outbuf[0])))
	fmt.Println(hex.EncodeToString(outbuf))

}
