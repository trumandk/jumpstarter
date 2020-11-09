package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
)

func getHeader(w http.ResponseWriter) {
	var result JumpStarterStatus
	v := reflect.ValueOf(result)
	typeOfS := v.Type()

	for i := 0; i < v.NumField(); i++ {
		fmt.Fprintf(w, "<th>%s</th>", typeOfS.Field(i).Name)
	}
}

func getStat(w http.ResponseWriter, ip string) {
	res, err := http.Get("http://" + ip + ":4")

	if err != nil {
		return
	}

	data, _ := ioutil.ReadAll(res.Body)

	res.Body.Close()

	var result JumpStarterStatus
	json.Unmarshal([]byte(data), &result)
	v := reflect.ValueOf(result)

	for i := 0; i < v.NumField(); i++ {
		if result.Running > 0 {
			fmt.Fprintf(w, "<td>\n")
		} else {
			fmt.Fprintf(w, "<td style=\"background-color:#cce5ff\">\n")
		}
		fmt.Fprintf(w, "%v</td>", v.Field(i).Interface())
	}
}
func menu(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "<head>")
	fmt.Fprintf(w, "<title>JumpStarter</title>")
	fmt.Fprintf(w, " <link rel=\"icon\" type=\"image/png\" href=\"files/jumpstarter.png\">")
	fmt.Fprintf(w, "<meta http-equiv=\"refresh\" content=\"5\">")
	fmt.Fprintf(w, "</head>")
	fmt.Fprintf(w, "<link rel=\"stylesheet\" href=\"files/bootstrap.css\">")
	fmt.Fprintf(w, "<script src=\"files/bootstrap.js\"></script>")
	fmt.Fprintf(w, "<body>")
	fmt.Fprintf(w, "<nav class=\"navbar navbar-expand-lg navbar-dark bg-dark\">")
	fmt.Fprintf(w, "<a class=\"navbar-brand\" href=\"/\">Servers</a>")
	fmt.Fprintf(w, "<a class=\"navbar-brand\" href=\"/containers\">Containers</a>")
	fmt.Fprintf(w, "<a class=\"navbar-brand\" href=\"/git\">Git</a>")
	fmt.Fprintf(w, "</nav>")

}
func servers(w http.ResponseWriter, req *http.Request) {
	menu(w, req)
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
	fmt.Fprintf(w, "<th scope=col>Ping</th>")
	getHeader(w)
	fmt.Fprintf(w, "</tr>")
	fmt.Fprintf(w, "</thead>")
	fmt.Fprintf(w, "<tbody>")

	for _, f := range nodes {
		if f.Name() != "env" && f.Name() != "all" {
			result, time := pingTest(f.Name())
			online := false
			if result {
				online = statOnline(f.Name())
			}
			if result && online {
				fmt.Fprintf(w, "<tr class=\"table-success\">\n")
			} else if result {
				fmt.Fprintf(w, "<tr class=\"table-warning\">\n")
			} else {
				fmt.Fprintf(w, "<tr class=\"table-danger\">\n")
			}
			fmt.Fprintf(w, "<td>\n")
			fmt.Fprintf(w, "<a href=\"/ssh?ip=%s&command=sudo reboot\" class=\"btn btn-warning\">Reboot</a>\n", f.Name())
			fmt.Fprintf(w, "<a href=\"/ssh?ip=%s&command=sudo poweroff\" class=\"btn btn-danger\">Poweroff</a>\n", f.Name())
			fmt.Fprintf(w, "</td>\n")
			fmt.Fprintf(w, "<td>%s</td>\n", f.Name())
			fmt.Fprintf(w, "<td>%s</td>\n", time)
			if online {
				getStat(w, f.Name())
			}
			fmt.Fprintf(w, "</tr>")

		}
	}
	fmt.Fprintf(w, "</tbody>")
	fmt.Fprintf(w, "</table>")
}
