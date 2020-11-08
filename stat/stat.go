package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/client"
	"github.com/docker/docker/api/types"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"log"
	"net/http"
	"os/exec"
	"time"
)
	var ctx = context.Background()

func status(w http.ResponseWriter, req *http.Request) {
	v, err := mem.VirtualMemory()

	if err != nil {
		fmt.Println("err:", err)
	}
	times, err1 := cpu.PercentWithContext(ctx, 0, false)
	if err1 != nil {
		log.Fatal(err1)
	}

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())

	if err != nil {
		panic(err)
	}
	info, err2 := cli.Info(ctx)
	if err2 != nil {
		panic(err2)
	}

	infoStat, err3 := host.InfoWithContext(ctx)
	if err3 != nil {
		panic(err3)
	}

	usageStat, err4 := disk.UsageWithContext(ctx,"/var/lib/docker")
	if err4 != nil {
		panic(err4)
	}

	m := JumpStarterStatus{
		fmt.Sprintf("%.2fGB", float64(v.Total)/1000000000),
		fmt.Sprintf("%.2fGB", float64(v.Available)/1000000000),
		//fmt.Sprintf("%.2f%%", float64(float64(v.Available/1000000)/float64(v.Total/1000000))*100),
		fmt.Sprintf("%.1f%%", times[0]),
		info.Containers,
		info.ContainersRunning,
		info.ContainersPaused,
		info.ContainersStopped,
		time.Duration(time.Duration(infoStat.Uptime) * time.Second).String(),
		fmt.Sprintf("%.2fGB", float64(usageStat.Total)/1000000000),
		fmt.Sprintf("%.2fGB", float64(usageStat.Free)/1000000000),
	}
	b, err := json.MarshalIndent(m, "", " ")
	if err != nil {
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func containers(w http.ResponseWriter, req *http.Request) {

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())

	if err != nil {
		panic(err)
	}

	containerList, err5 := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err5 != nil {
		panic(err5)
	}
	b, err := json.MarshalIndent(containerList, "", " ")
	if err != nil {
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func main() {
        exec.Command("/bin/mkdir","/var/lib/docker").CombinedOutput()
        out, err2 := exec.Command("/usr/bin/mount","-t","ext4","/dev/sda1","/var/lib/docker").CombinedOutput()

        if err2 != nil {
                fmt.Printf("Error Message:%s", err2)
        }
       output := string(out[:])
        fmt.Println(output)


	mux := http.NewServeMux()
	mux.HandleFunc("/", status)
	mux.HandleFunc("/containers", containers)
	err := http.ListenAndServe(":4", mux)
	log.Fatal(err)
}
