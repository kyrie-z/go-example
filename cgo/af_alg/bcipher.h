#ifndef BCIPHER_H
#define BCIPHER_H
#include <stdio.h>

struct crypt_cipher {
	int tfmfd;
	int opfd;
};

void crypt_cipher_init(struct crypt_cipher **ctx,char *key);
void crypt_cipher_encrypt(int opfd,char *in,size_t length,char *out);
void crypt_cipher_decrypt(int opfd,char *in,size_t length,char *out);
void crypt_cipher_destroy(struct crypt_cipher *ctx);
#endif //BCIPHER_H
