FROM golang:alpine AS builder

RUN apk update
RUN apk add --no-cache git
WORKDIR /app/

RUN go get github.com/pin/tftp
#RUN go get github.com/miekg/dns
#RUN go get github.com/sparrc/go-ping
COPY main.go main.go
#RUN go get -d -v
RUN CGO_ENABLED=0 go build -o /main

FROM alpine
RUN apk add --no-cache wget
WORKDIR /files/
RUN wget https://stable.release.flatcar-linux.net/amd64-usr/current/flatcar_production_pxe.vmlinuz
RUN wget https://stable.release.flatcar-linux.net/amd64-usr/current/flatcar_production_pxe_image.cpio.gz
WORKDIR /tftp/
RUN apk add --no-cache syslinux
RUN apk add --no-cache mlocate
RUN updatedb
RUN locate ldlinux.c32
RUN cp /usr/share/syslinux/lpxelinux.0 .
RUN cp /usr/share/syslinux/ldlinux.c32 .
COPY --from=builder /main /tftp/main
RUN mkdir pxelinux.cfg/
COPY default pxelinux.cfg/default
COPY pxe-config.ign /files/pxe-config.ign
ENTRYPOINT ["/tftp/main"]
