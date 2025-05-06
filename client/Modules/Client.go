package Modules

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func RunEcho(hostname string){
	n := bufio.NewReader(os.Stdin)
	fmt.Print("Enter Username > ")
	name, _ := n.ReadString('\n')
	serverdial, err := net.ResolveUDPAddr("udp4", hostname)
	if err != nil {
		fmt.Println("Client: (ERROR) Could not resolve hostname and addr -> ", err)
	} else {
		fmt.Println("Client: (OK) Hostname resolved")
	}
	clientdial, err := net.DialUDP("udp4", nil, serverdial)
	if err != nil {
		fmt.Println("Client: (ERROR) Unable to dial the hostname or server provided -> ", err)
	} else{
		fmt.Println("Client: (OK) Dial func Successful")
	}
	fmt.Printf("Client: (Info) Server is located at -> %s\n", clientdial.RemoteAddr().String())

	for {
		// data write loop
		t := bufio.NewReader(os.Stdin)
		fmt.Print("Send message > ")
		text, _ := t.ReadString('\n')
		data := []byte(name + ": " + text + "\n")
		_, err = clientdial.Write(data)
		if strings.TrimSpace(string(data)) == "ClientStop"{
			fmt.Println("Client: (OK) Shutting down client...")
			return
		}
		if err != nil {
			fmt.Println("Client: (ERROR) Got error when Writing buffer -> ", err)
			return
		} 
		
	}
}

