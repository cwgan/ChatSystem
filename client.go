package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	ServerIp string
	ServerPort int
	Name string
	conn net.Conn
	flag int // current mode of client
}

func NewClient(serverIp string, serverPort int) *Client {
	// create a client object 
	client := &Client{
		ServerIp: serverIp,
		ServerPort: serverPort,	
		flag: 999,
	}
	// link to server
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d",serverIp,serverPort))
	if err != nil {
		fmt.Println("net.Dial error:",err)
		return nil
	}
	client.conn = conn
	return client
}

var serverIp string
var serverPort int
// ./client -ip 127.0.0.1 -port 9999
func init(){
	flag.StringVar(&serverIp,"ip" , "127.0.0.1", "set server ip address (default:127.0.0.1)")
	flag.IntVar(&serverPort,"port",9999, "set server port (default:8888)")
}

func main(){
	// command line parse
	flag.Parse()

	client := NewClient(serverIp,serverPort)
	if client == nil {
		fmt.Println(">>>>>link to server failed")
	}

	fmt.Println(">>>>>link to server success")

	// single goroutine to deal with server response
	go client.DealResponse()

	// launch a client 
	client.Run()
}

func (client *Client) menu() bool {
	var flag int
	fmt.Println("1.public chat")
	fmt.Println("2.private chat")
	fmt.Println("3.rename")
	fmt.Println("0.quit")

	fmt.Scanln(&flag)
	if flag >= 0 && flag <= 3{
		client.flag = flag
		return true
	}else{
		fmt.Println("input invalid")
		return false
	}
}

func (client *Client) Run() {
	for client.flag != 0 {
		for client.menu() != true {
		}
		switch client.flag {
		case 1:
			//public
			fmt.Println("public chat mode ")
			client.PublicChat()
			break
		case 2:
			//private
			fmt.Println("private chat mode")
			client.PrivateChat()
			break
		case 3:
			//rename
			fmt.Println("rename")
			client.UpdateName()
			break
		}
	}
}

func (client *Client) UpdateName() bool {
	fmt.Println(">>>>> please enter new name:")
	fmt.Scanln(&client.Name)

	sendMsg := "rename|" + client.Name + "\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn writer err:" , err)
		return false
	}
	return true
}

func (client *Client) DealResponse(){
	// once conn has data, data will copy to stdout, forever listening
	io.Copy(os.Stdout, client.conn)
}

func (client *Client) PublicChat(){
	// 
	var chatMsg string
	fmt.Println("begin your public talk,type exit to quit")
	chatMsg = ReadLine()
	for chatMsg != "exit" {
		// send it to server
		if len(chatMsg) != 0 {
			sendMsg := chatMsg + "\n"
			_, err := client.conn.Write([]byte(sendMsg))
			if err != nil{
				fmt.Println("conn write err:", err)
				break
			}
		}
		chatMsg = ""
		chatMsg = ReadLine()
	}
}

func (client *Client) SelectUsers() {
	sendMsg := "who\n"
	_,err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn write err:", err)
	}
}

func ReadLine () string{
	reader := bufio.NewReader(os.Stdin)
	str,_,err := reader.ReadLine()
	if err != nil{	
		fmt.Println("reader.readline error:",err)
	}
	chatMsg := string(str)
	return chatMsg
}


func (client *Client) PrivateChat(){
	var remoteName string
	var chatMsg string

	client.SelectUsers()
	fmt.Println(">>>>>please enter user name,exit for quit.")
	fmt.Scan(&remoteName)
	fmt.Println(">>>>>begin your private talk,exit for quit.")
	for remoteName != "exit" {
		chatMsg = ReadLine()
		for chatMsg != "exit" {
			if len(remoteName) != 0 {
				sendMsg := "to|" + remoteName + "|" + chatMsg + "\n\n"
				_, err := client.conn.Write([]byte(sendMsg))
				if err != nil {
					fmt.Println("conn write err:",err)
					break
				}
			}
			chatMsg = ""
			chatMsg = ReadLine()
		}	
	}
	client.SelectUsers()
	fmt.Println(">>>>>please enter user name:")
	fmt.Scan(&remoteName)
}


