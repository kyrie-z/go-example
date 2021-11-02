#include <stdio.h>
#include <unistd.h>
#include <sys/socket.h>
#include <linux/if_alg.h>
#include <linux/socket.h>
#include <string.h>
#include <stdlib.h>
#include "bcipher.h"

#ifndef SOL_ALG
#define SOL_ALG 279
#endif


void crypt_cipher(int opfd,char *in,size_t length,char *out,int alg_option){
  char cbuf[CMSG_SPACE(4) + CMSG_SPACE(20)] = {0};
  struct msghdr msg = {};
  struct cmsghdr *cmsg;

  struct af_alg_iv *iv;
  struct iovec iov;

  msg.msg_control = cbuf;
  msg.msg_controllen = sizeof(cbuf);

  cmsg = CMSG_FIRSTHDR(&msg);
  cmsg->cmsg_level = SOL_ALG;
  cmsg->cmsg_type = ALG_SET_OP;
  cmsg->cmsg_len = CMSG_LEN(4);
  *(__u32 *)CMSG_DATA(cmsg) = alg_option;

  cmsg = CMSG_NXTHDR(&msg, cmsg);
  cmsg->cmsg_level = SOL_ALG;
  cmsg->cmsg_type = ALG_SET_IV;
  cmsg->cmsg_len = CMSG_LEN(20);
  iv = (void *)CMSG_DATA(cmsg);
  iv->ivlen = 16;
  memcpy(iv->iv, "\x3d\xaf\xba\x42\x9d\x9e\xb4\x30"
           "\xb4\x22\xda\x80\x2c\x9f\xac\x41", 16);

  iov.iov_base = in;
  iov.iov_len = length;

  msg.msg_iov = &iov;
  msg.msg_iovlen = 1;

  sendmsg(opfd, &msg, 0);
  read(opfd, out, 16);

}

void crypt_cipher_encrypt(int opfd,char *in,size_t length,char *out)
{
	return crypt_cipher(opfd, in, length, out,ALG_OP_ENCRYPT);
}

void crypt_cipher_decrypt(int opfd,char *in,size_t length,char *out)
{
	return crypt_cipher(opfd, in, length, out,ALG_OP_DECRYPT);
}


void crypt_cipher_init(struct crypt_cipher **ctx,char *key)
{
  *ctx = malloc(sizeof(**ctx));

  struct sockaddr_alg sa = {
    .salg_family = AF_ALG,
    .salg_type = "skcipher",
    .salg_name = "cbc(aes)"
  };

  (*ctx)->tfmfd = socket(AF_ALG, SOCK_SEQPACKET, 0);
  bind((*ctx)->tfmfd, (struct sockaddr *)&sa, sizeof(sa));
  setsockopt((*ctx)->tfmfd, SOL_ALG, ALG_SET_KEY,key, 16);
  (*ctx)->opfd = accept((*ctx)->tfmfd, NULL, 0);
  
}

void crypt_cipher_destroy(struct crypt_cipher *ctx)
{
	if (ctx->tfmfd >= 0)
		close(ctx->tfmfd);
	if (ctx->opfd >= 0)
		close(ctx->opfd);
	memset(ctx, 0, sizeof(*ctx));
	free(ctx);
}

// int main(void)
// {

//   char buf[16];
//   int i;
//   char key[] ={0x06,0xa9,0x21,0x40,0x36,0xb8,0xa1,0x5b,0x51,0x2e,0x03,0xd5,0x34,0x12,0x00,0x06};

//   struct crypt_cipher *h;
//   crypt_cipher_init(&h,key);

//   // 明文
//   char *plain="plain test";
//   printf("data: %s\n",plain);

//   // 加密
//   printf("ENCRYPT: ");
//   crypt_cipher_encrypt(h->opfd,plain,16,buf);
//   for (i = 0; i < 16; i++) {
//     printf("%02x", (unsigned char)buf[i]);
//   }
//   printf("\n");

//   // 解密
//   printf("DECRYPT: ");
//   char plain_out[16];
//   crypt_cipher_decrypt(h->opfd,buf,16,plain_out);
//   printf("%s\n", plain_out);

//   crypt_cipher_destroy(h);
//   return 0;
// }
