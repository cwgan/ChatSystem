package main

import (
	"net"
	"strings"
)

type User struct{
	Name string
	Addr string
	C chan string
	conn net.Conn
	server *Server
}

// user api
func NewUser(conn net.Conn, server *Server) *User{
	userAddr := conn.RemoteAddr().String()
	user := &User{
		Name: userAddr,
		Addr: userAddr,
		C: make(chan string),
		conn: conn,
		server: server,
	}
	go user.ListenMessage()
	return user
}

//
func (this *User) ListenMessage(){
	for{
		msg := <-this.C
		this.conn.Write([]byte(msg + "\n"))
	}
}

func (this *User) Online(){
	
	// put user in onlinemap
	this.server.mapLock.Lock()
	this.server.OnLineMap[this.Name] = this
	this.server.mapLock.Unlock()

	// broadcast that this user is online
	this.server.BroadCast(this,"online already")
}

func (this *User) Offline(){
	// delete user from onlineusermap
	this.server.mapLock.Lock()
	delete(this.server.OnLineMap, this.Name)
	this.server.mapLock.Unlock()

	// broadcast that this user is offline
	this.server.BroadCast(this,"offline")
}

func (this *User) SendMsg(msg string){
	this.conn.Write([]byte(msg))
}

func (this *User) DoMessage(msg string){
	if msg == "who" {
		// search who is online
		this.server.mapLock.Lock()
		for _, user := range this.server.OnLineMap {
			onlineMsg := "[" + user.Addr + "] " + user.Name + "\n"
			this.SendMsg(onlineMsg)
		}
		this.server.mapLock.Unlock()
	} else if len(msg) > 7 && msg[:7] == "rename|" {
		newName := strings.Split(msg,"|") [1]

		// check if name exist
		_,ok := this.server.OnLineMap[newName]
		if ok {
			this.SendMsg("this name is already used\n")
		}else {
			this.server.mapLock.Lock()
			delete(this.server.OnLineMap, this.Name)
			this.server.OnLineMap[newName] = this
			this.server.mapLock.Unlock()

			this.Name = newName
			this.SendMsg("rename success, your new name: " + newName + "\n")
		}
	} else if len(msg) >4 && msg[:3] == "to|" {
		// get name
		remoteName := strings.Split(msg,"|") [1]
		if remoteName == "" {
			this.SendMsg("format error \n")
			return
		}
		// get user object
		remoteUser, ok := this.server.OnLineMap[remoteName]
		if !ok {
			this.SendMsg("this user is not exist\n")
			return
		}
		// get message , send 
		content := strings.Split(msg,"|")[2]
		if content == ""{
			this.SendMsg("no content, please try again\n")
			return
		}
		remoteUser.SendMsg(this.Name + " (private) :" + content)
	} else {
		this.server.BroadCast(this, msg)
	}
	
}