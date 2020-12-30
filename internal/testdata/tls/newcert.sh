#!/bin/sh

openssl req -x509 -nodes -new -sha256 -days 1024 -newkey rsa:2048 -keyout rootCA.key -out rootCA.pem -subj "/C=US/CN=DDL-TEST"
openssl x509 -outform pem -in rootCA.pem -out rootCA.crt
openssl req -new -nodes -newkey rsa:2048 -keyout localhost.key -out localhost.csr -subj "/C=US/ST=PA/L=Philadelphia/O=DDL-TEST/CN=localhost.local"
openssl x509 -req -sha256 -days 1024 -in localhost.csr -CA rootCA.pem -CAkey rootCA.key -CAcreateserial -extfile domains.ext -out localhost.crt
