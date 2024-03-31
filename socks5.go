package main

import (
	"fmt"
	"github.com/txthinking/socks5"
	"log"
	"net"
	"os"
)

type Config struct {
	Address string
	IP      string

	Username string
	Password string
}

type HandlerMiddleware struct {
	handler socks5.Handler
}

func (h *HandlerMiddleware) TCPHandle(server *socks5.Server, conn *net.TCPConn, request *socks5.Request) error {
	dest := net.IP(request.DstAddr)
	log.Printf("Got TCP message: address=%s; destination=%s", conn.RemoteAddr(), dest.String())

	return h.handler.TCPHandle(server, conn, request)
}

func (h *HandlerMiddleware) UDPHandle(server *socks5.Server, addr *net.UDPAddr, datagram *socks5.Datagram) error {
	dest := net.IP(datagram.DstAddr)
	log.Printf("Got UDP message: address=%s; destination=%s", addr.IP.String(), dest.String())

	return h.handler.UDPHandle(server, addr, datagram)
}

func main() {
	config, err := fetchConfig()
	if err != nil {
		log.Fatalf("Cannot get config for start proxy: %s", err.Error())
	}

	server, _ := socks5.NewClassicServer(config.Address, config.IP, config.Username, config.Password, 30, 30)

	middleware := HandlerMiddleware{
		handler: &socks5.DefaultHandle{},
	}
	err = server.ListenAndServe(&middleware)
	if err != nil {
		log.Fatalf("Error while listen proxy: %s", err.Error())
	}
}

func fetchConfig() (Config, error) {
	addr, ok := os.LookupEnv("PROXY_ADDRESS")
	if !ok {
		return Config{}, fmt.Errorf("PROXY_ADDRESS env required")
	}

	ip, ok := os.LookupEnv("PROXY_IP")
	if !ok {
		return Config{}, fmt.Errorf("PROXY_IP env required")
	}

	username, _ := os.LookupEnv("PROXY_USERNAME")
	password, _ := os.LookupEnv("PROXY_PASSWORD")

	return Config{
		Address: addr,
		IP:      ip,

		Username: username,
		Password: password,
	}, nil
}
