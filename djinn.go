package main

import (
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func read(conn net.Conn) {
	buf := make([]byte, 340)
	conn.Read([]byte(buf))
	for i := 0; i <= 1010; i++ {
		conn.Read([]byte(buf))
		r, _ := regexp.Compile("[0-9], '.', [0-9]")
		find := r.FindString(string(buf))
		find = strings.ReplaceAll(find, "'", "")
		find = strings.ReplaceAll(find, ",", "")
		test := strings.Fields(find)
		fmt.Print(string(buf))
		buf = make([]byte, 400)
		var a [3]string
		for i, v := range test {
			a[i] = v
		}
		firt, _ := strconv.Atoi(a[0])
		second, _ := strconv.Atoi(a[2])
		operator := a[1]
		var data string
		switch operator {
		case "+":
			data = strconv.Itoa(firt + second)
		case "-":
			data = strconv.Itoa(firt - second)
		case "/":
			data = strconv.Itoa(firt / second)
		case "*":
			data = strconv.Itoa(firt * second)
		}
		conn.Write([]byte(data + "\n"))
		time.Sleep(1 * time.Second)
	}

}

func main() {
	fmt.Println("Connecting to 10.0.1.8")
	conn, err := net.Dial("tcp", "10.0.1.8:1337")
	if err != nil {
		fmt.Print("Cannot Connect")
	}
	read(conn)

}
