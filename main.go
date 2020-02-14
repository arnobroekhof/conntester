package main

import (
	"context"
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

const (
	emojiCheckMark = "\xE2\x9C\x85"
	emojiCrossMark = "\xE2\x9D\x8C"
	checkRemoteIPUrl = "http://checkip.amazonaws.com/"
)

type Connections struct {
	Services []Service `yaml:"services"`
}

type Service struct {
	Name     string `yaml:"name"`
	Port     int    `yaml:"port"`
	Host     string `yaml:"host"`
	Protocol string `yaml:"protocol"`
}

var (
	configFile string
	withHostInfo bool
)


func init() {
	flag.StringVar(&configFile, "config", "config.yaml", "the config file to use")
	flag.BoolVar(&withHostInfo, "print-host-info", true, "print network info from the host")
	flag.Parse()
}

func main() {

	config, err := ioutil.ReadFile(configFile)
	if err != nil {
		fmt.Printf("error reading file: %s\n", err)
	}

	connections := Connections{}
	err = yaml.Unmarshal(config, &connections)
	if err != nil {
		fmt.Printf("error while parsing config file: %s\n", err)
	}

	if withHostInfo {
		printHostNetworkInfo()
	}

	waitGroup := sync.WaitGroup{}
	waitGroup.Add(1)
	go func() {
		for _, srv := range connections.Services {
			waitGroup.Add(1)
			printOkOrError(testService(srv))
			waitGroup.Done()

		}
		waitGroup.Done()
	}()
	waitGroup.Wait()
}

func testService(srv Service) (bool, string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	switch strings.ToLower(srv.Protocol) {
	case "tcp":
		ok, msg := testTCP(srv, ctx)
		if ok {
			return true, fmt.Sprintf("able to connect to %s = %s:%s:%d", srv.Name, srv.Protocol, srv.Host, srv.Port)
		} else {
			return false, fmt.Sprintf("unable to connect to %v, error: %s", srv, msg)
		}
	default:
		return false, fmt.Sprintf("protocol %s not supported", srv.Protocol)
	}
}

func testTCP(srv Service, ctx context.Context) (bool, string) {

	_, err := net.ResolveTCPAddr(strings.ToLower(srv.Protocol), fmt.Sprintf("%s:%d", srv.Host, srv.Port))
	if err != nil {
		return false, err.Error()
		ctx.Err()
	}

	conn, err := net.DialTimeout(strings.ToLower(srv.Protocol), fmt.Sprintf("%s:%d", srv.Host, srv.Port), time.Second*5)
	if err != nil {
		ctx.Err()
		return false, err.Error()
	}
	defer conn.Close()
	ctx.Done()
	return true, ""
}

func printOkOrError(ok bool, msg string) {
	switch ok {
	case true:
		fmt.Printf("%s  --> %s\n", []byte(emojiCheckMark), msg)
	default:
		fmt.Printf("%s  --> %s\n", []byte(emojiCrossMark), msg)
	}
}

func printHostNetworkInfo() {
	fmt.Printf("Local IP Address: %s\n", getLocalIP())
	fmt.Printf("Remote IP Address: %s\n", getRemoteIP())
}

func getRemoteIP() string {
	res, err := http.Get(checkRemoteIPUrl)
	if err != nil {
		return ""
	}
	body, err := ioutil.ReadAll(res.Body)
	return string(body)
}

func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}
