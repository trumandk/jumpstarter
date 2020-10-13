all:
	gofmt -w main.go
	./ct <pxe-config > pxe-config.ign
	docker-compose stop
	docker-compose build
	docker-compose up
test:
	dig @10.0.1.4 -p 53 test.service
	dig @10.0.1.4 -p 53 test2.service
	dig @10.0.1.4 -p 53 test3.service
	dig @10.0.1.4 -p 53 test10.service
	dig @10.0.1.4 -p 53 skod.com
	dig @10.0.1.4 -p 53 nissefar.com

ca:
	openssl genrsa -out conf/rootCA.key 4096
	openssl req -x509 -new -nodes -key conf/rootCA.key -sha256 -days 20240 -out conf/rootCA.crt

ssl:
	openssl genrsa -out conf/test.service.key 2048
	openssl req -new -sha256 -key conf/test.service.key -subj "/C=US/ST=CA/O=MyOrg, Inc./CN=test.service" -out conf/test.service.csr
	openssl req -in conf/test.service.csr -noout -text
	openssl x509 -req -in conf/test.service.csr -CA conf/rootCA.crt -CAkey conf/rootCA.key -CAcreateserial -out conf/test.service.crt -days 20240 -sha256
	openssl x509 -in conf/test.service.crt -text -noout
info:
	openssl genrsa -out conf/info.service.key 2048
	openssl req -new -sha256 -key conf/info.service.key -subj "/C=US/ST=CA/O=MyOrg, Inc./CN=info.service" -out conf/info.service.csr
	openssl req -in conf/info.service.csr -noout -text
	openssl x509 -req -in conf/info.service.csr -CA conf/rootCA.crt -CAkey conf/rootCA.key -CAcreateserial -out conf/info.service.crt -days 20240 -sha256
	openssl x509 -in conf/info.service.crt -text -noout

keystore:
	openssl pkcs12 -export -in conf/info.service.crt -inkey conf/info.service.key -out conf/info.service.p12
	keytool -importkeystore -srckeystore conf/info.service.p12 -srcstoretype PKCS12 -destkeystore info.service.jks  -deststoretype JKS

