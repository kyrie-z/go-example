package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"u_rsa/u_rsa"
)

func main() {
	/* generate rsa key and save */
	prikey, pubkey := u_rsa.GenRsaKey(2048)
	err := u_rsa.SaveKey(prikey, "prikey.pem")
	if err != nil {
		fmt.Println(err)
	}
	err = u_rsa.SaveKey(pubkey, "pubkey.pem")
	if err != nil {
		fmt.Println(err)
	}

	/* read key.pem file and sign */
	f, err := ioutil.ReadFile("./prikey.pem")
	if err != nil {
		fmt.Println("read file fail!!")
		return
	}
	block, _ := pem.Decode(f)
	pkey, _ := x509.ParsePKCS1PrivateKey(block.Bytes)
	signData := u_rsa.RsaSignWithSha256([]byte("this is a plain txt."), pkey)
	fmt.Println("签名数据：", signData)

	/* read pubkey and signdata verify */
	f, err = ioutil.ReadFile("./pubkey.pem")
	if err != nil {
		fmt.Println("read file fail!!")
		return
	}
	block, _ = pem.Decode(f)
	fkey, err := x509.ParsePKIXPublicKey(block.Bytes)
	ret := u_rsa.RsaVerifyWithSha256([]byte("this is a plain txt."), signData, fkey.(*rsa.PublicKey))
	if ret {
		fmt.Println("verify success!!")
	}

	/* pubkey encrypt and prikey decrypt */
	ciphertxt := u_rsa.RsaEncrypt([]byte("this is a plain txt!!"), fkey.(*rsa.PublicKey))
	fmt.Println("加密文：", ciphertxt)
	txt := u_rsa.RsaDecrypt(ciphertxt, pkey)
	fmt.Println("解密文：", string(txt))
}
