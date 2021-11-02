#include "../bcipher.h"


int main(void)
{

  char buf[16];
  int i;
  char key[] ={0x06,0xa9,0x21,0x40,0x36,0xb8,0xa1,0x5b,0x51,0x2e,0x03,0xd5,0x34,0x12,0x00,0x06};

  struct crypt_cipher *h;
  crypt_cipher_init(&h,key);
  printf("opfd=%d\n",h->opfd);
  printf("tfmfd=%d\n",h->tfmfd);

  // 明文
  char *plain="plain test test test";
  printf("data: %s\n",plain);

  // 加密
  printf("ENCRYPT: ");
  crypt_cipher_encrypt(h->opfd,plain,16,buf);
  for (i = 0; i < 16; i++) {
    printf("%02x", (unsigned char)buf[i]);
  }
  printf("\n");

  // 解密
  printf("DECRYPT: ");
  char plain_out[16];
  crypt_cipher_decrypt(h->opfd,buf,16,plain_out);
  printf("%s\n", plain_out);

  crypt_cipher_destroy(h);
  return 0;
}
