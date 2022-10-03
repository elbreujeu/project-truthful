openssl genrsa -out server/cert/id_rsa 4096
openssl rsa -in server/cert/id_rsa -pubout -out server/cert/id_rsa.pub