

## create root cert:
`openssl req -new -x509 -newkey rsa:2048 -subj "/CN=u_cert root/" -keyout root.key -out root.crt -days 3650 -nodes -sha256`

## verify cert chain
`openssl verify -verbose -CAfile root.crt  usrCA.crt`
