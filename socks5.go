package main

import (
	"fmt"
	"github.com/txthinking/socks5"
	"log"
	"os"
)

type Config struct {
	Address string
	IP      string

	Username string
	Password string
}

func main() {
	config, err := fetchConfig()
	if err != nil {
		log.Fatalf("Cannot get config for start proxy: %s", err.Error())
	}

	server, _ := socks5.NewClassicServer(config.Address, config.IP, config.Username, config.Password, 30, 30)

	err = server.ListenAndServe(nil)
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
