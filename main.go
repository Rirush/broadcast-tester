package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

const Usage = `Usage:
  broadcast-tester mserver ADDRESS PORT MESSAGE - start multicast server at ADDRESS:PORT and respond to all packets with MESSAGE
  broadcast-tester userver ADDRESS PORT MESSAGE - start UDP server at ADDRESS:PORT and respond to all packets with MESSAGE
  broadcast-tester client ADDRESS PORT MESSAGE - send packet to UDP server at ADDRESS:PORT with MESSAGE`

func main() {
	if len(os.Args) != 5 {
		fmt.Println(Usage)
		return
	}
	switch os.Args[1] {
	case "mserver":
		mserver()
	case "userver":
		userver()
	case "client":
		client()
	default:
		fmt.Println(Usage)
		return
	}
}

func mserver() {
	address, port, message := os.Args[2], os.Args[3], os.Args[4]
	addr, err := net.ResolveUDPAddr("udp", address + ":" + port)
	if err != nil {
		log.Fatal("Cannot resolve address:", err)
	}
	listener, err := net.ListenMulticastUDP("udp", nil, addr)
	if err != nil {
		log.Fatal("Cannot start multicast server:", err)
	}
	log.Println("Listening for multicast packets at", address + ":" + port)
	for {
		b := make([]byte, 1024)
		n, src, err := listener.ReadFromUDP(b)
		if err != nil {
			log.Fatal("Cannot read from UDP:", err)
		}
		log.Println("Received message from", src, "with content", string(b[:n]))
		_, err = listener.WriteToUDP([]byte(message), src)
		if err != nil {
			log.Fatal("Cannot respond to multicast:", err)
		}
	}
}

func client() {
	address, port, message := os.Args[2], os.Args[3], os.Args[4]
	addr, err := net.ResolveUDPAddr("udp", address + ":" + port)
	if err != nil {
		log.Fatal("Cannot resolve address:", err)
	}
	conn, err := net.ListenUDP("udp", nil)
	if err != nil {
		log.Fatal("Cannot dial specified address:", err)
	}
	_, err = conn.WriteTo([]byte(message), addr)
	if err != nil {
		log.Fatal("Cannot send multicast packet:", err)
	}
	b := make([]byte, 1024)
	n, src, err := conn.ReadFrom(b)
	if err != nil {
		log.Fatal("Cannot read from UDP:", err)
	}
	log.Println("Received a response from", src, "with content", string(b[:n]))
}

func userver() {
	address, port, message := os.Args[2], os.Args[3], os.Args[4]
	addr, err := net.ResolveUDPAddr("udp", address + ":" + port)
	if err != nil {
		log.Fatal("Cannot resolve address:", err)
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatal("Cannot dial specified address:", err)
	}
	log.Println("Listening for UDP packets at", address + ":" + port)
	for {
		b := make([]byte, 1024)
		n, src, err := conn.ReadFrom(b)
		if err != nil {
			log.Fatal("Cannot read from UDP:", err)
		}
		log.Println("Received a message from", src, "with content", string(b[:n]))
		_, err = conn.WriteTo([]byte(message), src)
		if err != nil {
			log.Fatal("Cannot respond to packet:", err)
		}
	}
}