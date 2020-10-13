package main

import (
	"fmt"
	"os"
	"time"
	"io"
	"log"
	"net/http"
	"github.com/pin/tftp"
)

func readHandler(filename string, r io.ReaderFrom) error {
    file, err := os.Open(filename)
    fmt.Printf("open: %s\n", filename)
      raddr := r.(tftp.OutgoingTransfer).RemoteAddr()
        laddr := r.(tftp.RequestPacketInfo).LocalIP()
        log.Println("RRQ from", raddr.String(), "To ",laddr.String())
        log.Println("")

    if err != nil {
        fmt.Fprintf(os.Stderr, "%v\n", err)
        return err
    }
    // Optional tsize support.
    // Set transfer size before calling ReadFrom.
    if t, ok := r.(tftp.OutgoingTransfer); ok {
        if fi, err := file.Stat(); err == nil {
            t.SetSize(fi.Size())
        }
    }
    //buf := bytes.NewBufferString("R29waGVycyBydWxlIQ==")
    n, err := r.ReadFrom(file)
    if err != nil {
        fmt.Fprintf(os.Stderr, "%v\n", err)
        return err
    }
    fmt.Printf("%d bytes sent\n", n)
    return nil
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("/files/")))
     go func() {
                http.ListenAndServe(":80", nil)
        }()

    // use nil in place of handler to disable read or write operations
    s := tftp.NewServer(readHandler, nil)
    s.SetTimeout(5 * time.Second) // optional
    err := s.ListenAndServe(":69") // blocks until s.Shutdown() is called
    if err != nil {
        fmt.Fprintf(os.Stdout, "server: %v\n", err)
        os.Exit(1)
    }
}
