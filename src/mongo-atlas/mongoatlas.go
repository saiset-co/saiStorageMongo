package main

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"net"
	"crypto/tls"
)

func main() {
	//dialInfo, err := mgo.ParseURL("Kirill:2PF8qDc84ciLaeg4@cluster0-qglca.mongodb.net/test")//sai_user:maxdata123@127.0.0.1/sai")
	//
	//fmt.Println(err)
	//tlsConfig := &tls.Config{}
	//
	//dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
	//	conn, err := tls.Dial("tcp", addr.String(), tlsConfig)
	//	return conn, err
	//}
	//
	//session, err := mgo.DialWithInfo(dialInfo)
	//
	//fmt.Println(err)
	//defer session.Close()

	dialInfo, err := mgo.ParseURL("Kirill:2PF8qDc84ciLaeg4@cluster0-qglca.mongodb.net/test")
	fmt.Println(err)
	dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
		tlsConfig := &tls.Config{}
		tlsConfig.InsecureSkipVerify = true
		conn, err := tls.Dial("tcp", addr.String(), tlsConfig)
		fmt.Println(err)
		return conn, err
	}
	session, err := mgo.DialWithInfo(dialInfo)
	fmt.Println(err)

	collection := session.DB("test").C("sai")

	var document interface{}
	if err := collection.Find(nil).One(&document); err != nil {
		fmt.Println(err)
	}

	fmt.Println(document)
}
