package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"math/big"
	"net"
	"os"
	"time"
)

// save key
func SaveKey(key []byte, filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	block, _ := pem.Decode(key)
	err = pem.Encode(f, block)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	/* 读取根证书和私钥*/
	caFile, err := ioutil.ReadFile("./root.crt")
	if err != nil {
		return
	}
	caBlock, _ := pem.Decode(caFile)
	rootCert, err := x509.ParseCertificate(caBlock.Bytes)
	if err != nil {
		fmt.Println("parse cert fail!!")
		return
	}
	// fmt.Println(rootCert.Subject)
	// fmt.Println(rootCert.Issuer)
	keyFile, err := ioutil.ReadFile("./root.key")
	if err != nil {
		fmt.Println("read root.key fail!!")
		return
	}
	keyBlock, rest := pem.Decode(keyFile)
	if rest == nil {
		fmt.Println("decode prikey fail!!")
		return
	}
	// "----BEGIN PRIVATE KEY----"   usr PKCS8 prase
	rootkey, err := x509.ParsePKCS8PrivateKey(keyBlock.Bytes)
	if err != nil {
		fmt.Println("parse prikey fail!!")
		return
	}

	/* 证书公共信息模板 */
	rd, _ := rand.Int(rand.Reader, big.NewInt(100))
	cer := &x509.Certificate{
		SerialNumber: big.NewInt(rd.Int64()), //证书序列号
		Subject: pkix.Name{
			Country:            []string{"CN"},
			Organization:       []string{"u_cert"},
			OrganizationalUnit: []string{"u_cert.unit"},
			Province:           []string{"HuBei"},
			CommonName:         "u_cert CommonNmae",
			Locality:           []string{"WuHan"},
		},
		NotBefore:             time.Now(),                                                                 //证书有效期开始时间
		NotAfter:              time.Now().AddDate(1, 0, 0),                                                //证书有效期结束时间
		BasicConstraintsValid: true,                                                                       //基本的有效性约束
		IsCA:                  false,                                                                      //是否是根证书
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth}, //证书用途(客户端认证，数据加密)
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageDataEncipherment,
		EmailAddresses:        []string{"u_cert@cert.com"},
		IPAddresses:           []net.IP{net.ParseIP("127.0.0.1")},
	}
	usrKeypair, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return
	}
	//create cert
	ca, err := x509.CreateCertificate(rand.Reader, cer, rootCert, &usrKeypair.PublicKey, rootkey.(*rsa.PrivateKey))
	if err != nil {
		return
	}
	caPem := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: ca,
	}
	usrCa := pem.EncodeToMemory(caPem)

	keyDer := x509.MarshalPKCS1PrivateKey(usrKeypair)
	keyPem := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: keyDer,
	}
	usrPrikey := pem.EncodeToMemory(keyPem)
	// save
	SaveKey(usrPrikey, "usrPrikey.key")
	SaveKey(usrCa, "usrCA.crt")

}
