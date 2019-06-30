package main

import (
	"context"
	"io/ioutil"

	"github.com/grandcat/zeroconf"
	"gopkg.in/yaml.v2"

	"flag"
	"fmt"
	"log"
	"os"
)

var (
	configFile string
)

func init() {
	flag.StringVar(&configFile, "config", "", "Yaml formatted configuration file")
	flag.Parse()

	if len(configFile) == 0 {
		fmt.Println("No configuration file provided. please set the -config flag.")
		os.Exit(2)
	}
}

type Config struct {
	Metadata                                             []string
	Port                                                 int
	NetworkInterface, Name, Type, Protocol, Domain, Mode string
}

func serviceMode(config Config) {
	log.Println("Running in service mode.")
	service, err := zeroconf.Register(
		config.Name,                     //  instance name
		config.Type+"."+config.Protocol, //  type and protocl
		config.Domain,                   // service domain
		config.Port,                     // service port
		config.Metadata,                 // service metadata
		nil,                             // register on all network interfaces
	)

	if err != nil {
		log.Fatal(err)
	}

	defer service.Shutdown()

	// Sleep forever
	select {}
}

func discoveryMode(config Config) {
	log.Println("Running in discovery mode")
	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		log.Fatal(err)
	}

	// Channel to receive discoveryed service entries
	entries := make(chan *zeroconf.ServiceEntry)

	go func(results <-chan *zeroconf.ServiceEntry) {
		for entry := range results {
			log.Println("Found service:", entry.ServiceInstanceName(), entry.Text, "at", entry.AddrIPv4)
		}

	}(entries)

	ctx := context.Background()

	err = resolver.Browse(ctx, config.Type+"."+config.Protocol, config.Domain, entries)
	if err != nil {
		log.Fatalln("Failed to browse:", err.Error())
	}

	<-ctx.Done()
}

func main() {
	sourceBytes, err := ioutil.ReadFile(configFile)
	if err != nil {
		panic(err)
	}

	var config Config

	if err := yaml.Unmarshal(sourceBytes, &config); err != nil {
		panic(err)
	}

	if config.Mode == "service" {
		serviceMode(config)
	} else if config.Mode == "discovery" {
		discoveryMode(config)
	} else {
		log.Println("The mode provided doesn't exist. Please choose either \"discover\" or \"service.\"")
		os.Exit(1)
	}
}
