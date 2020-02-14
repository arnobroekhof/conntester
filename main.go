package main

import (
	"context"
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net"
	"strings"
	"sync"
	"time"
)

const (
	emojiCheckMark = "\xE2\x9C\x85"
	emojiCrossMark = "\xE2\x9D\x8C"
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

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "config.yaml", "the config file to use")
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
