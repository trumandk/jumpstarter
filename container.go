package main

import (
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	//"time"
)

func getContainer(w http.ResponseWriter, ip string) {
	res, err := http.Get("http://" + ip + ":4/containers")

	if err != nil {
		return
	}

	data, _ := ioutil.ReadAll(res.Body)

	res.Body.Close()

	var containers []types.Container
	json.Unmarshal([]byte(data), &containers)
	for _, container := range containers {
		if strings.Contains(container.Status, "minutes") || strings.Contains(container.Status, "hour") {
			fmt.Fprintf(w, "<tr class=\"table-success\">\n")
		} else {
			fmt.Fprintf(w, "<tr class=\"table-danger\">\n")
		}
		fmt.Fprintf(w, "<td><a href=\"/ssh?ip=%s&command=sudo docker restart %s\" class=\"btn btn-warning\">Restart</a>", ip, container.ID)
		fmt.Fprintf(w, " <a href=\"/ssh?ip=%s&command=sudo docker pull %s\" class=\"btn btn-warning\">RePull</a></td>", ip, container.Image)
		fmt.Fprintf(w, "<td>%v</td>", ip)
		fmt.Fprintf(w, "<td>")
		fmt.Fprintf(w, "<a href=\"/sshout?ip=%s&command=sudo docker logs --tail 500 %s\" target=\"_blank\">", ip, container.ID)
		for _, name := range container.Names {
			fmt.Fprintf(w, "%v ", name)
		}
		fmt.Fprintf(w, "</a>")
		fmt.Fprintf(w, "</td>")
		fmt.Fprintf(w, "<td>")
		for _, port := range container.Ports {
			if port.PublicPort > 0 {
				fmt.Fprintf(w, "<a href=http://%s:%v target=\"_blank\">%v</a>:%v:%v ", ip, port.PublicPort, port.PublicPort, port.PrivatePort, port.Type)
			}
		}
		fmt.Fprintf(w, "</td>")
		fmt.Fprintf(w, "<td>")
		if len(container.HostConfig.NetworkMode) > 0 {
			fmt.Fprintf(w, "%v", container.HostConfig.NetworkMode)
		}

		fmt.Fprintf(w, "</td>")
		fmt.Fprintf(w, "<td>%v</td>", container.Image)
		fmt.Fprintf(w, "<td>%v</td>", container.State)
		fmt.Fprintf(w, "<td>%v</td>", container.Status)
		//fmt.Fprintf(w, "<td>%s</td>", time.Unix(container.Created, 0))

		fmt.Fprintf(w, "</tr>")
	}
}

func containers(w http.ResponseWriter, req *http.Request) {
	menu(w, req)
	/*
		fmt.Fprintf(w, "<head>")
		fmt.Fprintf(w, "<title>JumpStarter</title>")
		fmt.Fprintf(w, " <link rel=\"icon\" type=\"image/png\" href=\"files/jumpstarter.png\">")
		fmt.Fprintf(w, "</head>")
		fmt.Fprintf(w, "<link rel=\"stylesheet\" href=\"files/css/bootstrap.css\">")
		fmt.Fprintf(w, "<body>")
	*/
	fmt.Fprintf(w, "<center>")
	nodes, err := ioutil.ReadDir("/git/docker/")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "<table class=\"table table-hover\">")
	fmt.Fprintf(w, "<thead>")
	fmt.Fprintf(w, "<tr>")
	fmt.Fprintf(w, "<th scope=col></th>")
	fmt.Fprintf(w, "<th scope=col>IP</th>")
	fmt.Fprintf(w, "<th scope=col>Name</th>")
	fmt.Fprintf(w, "<th scope=col>Port</th>")
	fmt.Fprintf(w, "<th scope=col>Network-mode</th>")
	fmt.Fprintf(w, "<th scope=col>Image</th>")
	fmt.Fprintf(w, "<th scope=col>State</th>")
	fmt.Fprintf(w, "<th scope=col>Status</th>")
	//fmt.Fprintf(w, "<th scope=col>Created</th>")
	//getHeader(w)
	fmt.Fprintf(w, "</tr>")
	fmt.Fprintf(w, "</thead>")
	fmt.Fprintf(w, "<tbody>")

	for _, f := range nodes {
		if f.Name() != "env" && f.Name() != "all" {
			result, _ := pingTest(f.Name())
			online := false
			if result {
				online = statOnline(f.Name())
			}
			if result && online {
				//		fmt.Fprintf(w, "<tr class=\"table-success\">\n")
			} else if result {
				fmt.Fprintf(w, "<tr class=\"table-warning\">\n")
				fmt.Fprintf(w, "<td scope=col>%v</td>", f.Name())
				fmt.Fprintf(w, "</tr>")
			} else {
				fmt.Fprintf(w, "<tr class=\"table-danger\">\n")
				fmt.Fprintf(w, "<td scope=col>%v</td>", f.Name())
				fmt.Fprintf(w, "</tr>")
			}
			/*
				fmt.Fprintf(w, "<td>\n")
				fmt.Fprintf(w, "<a href=\"/ssh?ip=%s&command=sudo reboot\" class=\"btn btn-warning\">Reboot</a>\n", f.Name())
				fmt.Fprintf(w, "<a href=\"/ssh?ip=%s&command=sudo poweroff\" class=\"btn btn-danger\">Poweroff</a>\n", f.Name())
				fmt.Fprintf(w, "</td>\n")
				fmt.Fprintf(w, "<td>%s</td>\n", f.Name())
				fmt.Fprintf(w, "<td>%s</td>\n", time)
				if online {
					getContainer(w, f.Name())
				}
				fmt.Fprintf(w, "</tr>")
			*/
			if online {
				getContainer(w, f.Name())
			}
		}
	}
	fmt.Fprintf(w, "</tbody>")
	fmt.Fprintf(w, "</table>")
}
