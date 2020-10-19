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

FROM alpine AS tftp 
RUN apk add --no-cache wget
RUN apk add --no-cache syslinux
WORKDIR /files/
RUN wget https://stable.release.flatcar-linux.net/amd64-usr/current/flatcar_production_pxe.vmlinuz
RUN wget https://stable.release.flatcar-linux.net/amd64-usr/current/flatcar_production_pxe_image.cpio.gz

FROM scratch
WORKDIR /files/
COPY --from=tftp /files/flatcar_production_pxe.vmlinuz .
COPY --from=tftp /files/flatcar_production_pxe_image.cpio.gz .
COPY pxe-config.ign .
COPY ssh ssh
WORKDIR /tftp/
COPY --from=tftp /usr/share/syslinux/lpxelinux.0 .
COPY --from=tftp /usr/share/syslinux/ldlinux.c32 .
WORKDIR /tftp/pxelinux.cfg/
COPY default default
WORKDIR /tftp/
COPY --from=builder /main /tftp/main
ENTRYPOINT ["/tftp/main"]
