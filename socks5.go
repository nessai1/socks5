package main

import (
	"fmt"
	"github.com/txthinking/socks5"
	"log"
	"net"
	"os"
	"strconv"
)

type Config struct {
	IP   string
	Port string

	Username string
	Password string

	TCPTimeout int
	UDPTimeout int
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

	server, err := socks5.NewClassicServer(fmt.Sprintf("%s:%s", config.IP, config.Port), config.IP, config.Username, config.Password, config.TCPTimeout, config.UDPTimeout)
	if err != nil {
		log.Fatalf("Cannot create socks5 proxy: %s", err.Error())
	}

	var authType string
	if config.Username != "" && config.Password != "" {
		authType = "username/password"
	} else {
		authType = "no-auth"
	}

	log.Printf("Socks5 server starts: address=%s; authType=%s", server.Addr, authType)

	middleware := HandlerMiddleware{
		handler: &socks5.DefaultHandle{},
	}
	err = server.ListenAndServe(&middleware)
	if err != nil {
		log.Fatalf("Error while listen proxy: %s", err.Error())
	}
}

func fetchConfig() (Config, error) {
	port, ok := os.LookupEnv("PROXY_PORT")
	if !ok {
		port = "1080"
	}

	ip, ok := os.LookupEnv("PROXY_IP")
	if !ok {
		return Config{}, fmt.Errorf("PROXY_IP env required")
	}

	username, _ := os.LookupEnv("PROXY_USERNAME")
	password, _ := os.LookupEnv("PROXY_PASSWORD")

	var tcp, udp int
	var err error
	tcpTimeout, ok := os.LookupEnv("PROXY_TCP_TIMEOUT")
	if !ok {
		tcp = 30
	} else {
		tcp, err = strconv.Atoi(tcpTimeout)
		if err != nil {
			return Config{}, fmt.Errorf("invalid format of given TCP_TIMEOUT env: %w", err)
		}
	}

	udpTimeout, ok := os.LookupEnv("PROXY_UDP_TIMEOUT")
	if !ok {
		udp = 30
	} else {
		udp, err = strconv.Atoi(udpTimeout)
		if err != nil {
			return Config{}, fmt.Errorf("invalid format of given UDP_TIMEOUT env: %w", err)
		}
	}

	return Config{
		IP:   ip,
		Port: port,

		Username: username,
		Password: password,

		TCPTimeout: tcp,
		UDPTimeout: udp,
	}, nil
}
