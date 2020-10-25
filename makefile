all:
#	gofmt -w main.go
	./ct -pretty  -in-file pxe-config -out-file pxe-config.ign
	docker-compose stop
	docker-compose build
	docker-compose up
