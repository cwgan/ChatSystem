package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type Server struct{
	Ip string
	Port int

	// online users table
	OnLineMap map[string]*User
	mapLock sync.RWMutex

	// message broadcast channel
	Message chan string
}

// server interface
func NewServer(ip string, port int) *Server{
	server := &Server{
		Ip: ip,
		Port: port,
		OnLineMap: make(map[string]*User),
		Message: make(chan string),
	}
	return server
}

// broadcast channel message to all the users
func (this *Server) ListenMessager(){
	for {
		// check if there are messages in channel
		msg := <- this.Message

		// send message to all users
		this.mapLock.Lock()
		for _, cli := range this.OnLineMap {
			cli.C <- msg
		}
		this.mapLock.Unlock()
	}
}

func (this *Server) BroadCast(user *User, msg string){
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" +msg
	this.Message <- sendMsg
}

func (this *Server) Handler(conn net.Conn){
	
	
	fmt.Println("tcp-connection established")
	// user online
	user := NewUser(conn,this)

	fmt.Println("user address:%s",user.Addr)

	
	user.Online()

	// checkk if this user online
	isLive := make(chan bool)

	// receive message from client
	go func(){
		buf := make([]byte, 4096)
		for{
			n, err := conn.Read(buf)
			if n == 0{
				user.Offline()
				return
			}

			if err != nil && err != io.EOF{
				fmt.Println("conn read err:", err)
				return
			}

			// extract user message
			msg := string(buf[:n-1])

			user.DoMessage(msg)

			// any message of user indicate this user is active
			isLive <- true
		}
	}()

	// timer
	for{
		select{
		case <-isLive:
			// active, reset timer
			// do nothing, just for activate select sentence
		case <-time.After(time.Second * 120):
			// out of time
			// user suspend
			user.SendMsg("you're suspend\n")

			// destroy
			close(user.C)
			conn.Close()

			// quit this handler
			return
		}
	}

}

func (this *Server) Start(){
	// listen
	listener, err := net.Listen("tcp",fmt.Sprintf("%s:%d",this.Ip,this.Port))
	if err != nil {
		fmt.Println("net.Listen err:" , err)
		return
	}

	// close listen socket
	defer listener.Close()

	// launch jianting go routine
	go this.ListenMessager()

	for{
		// accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener accept failed:", err)
			continue
		}

		// do handler
		go this.Handler(conn)
	}
}

