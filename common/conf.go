package common

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
)

// ZipperConfig middleware configuration structure
type ZipperConfig struct {
	// server name
	Name string
	// port parameter
	Port uint16
	// a queue message capacity parameter
	QueueSize uint16
	// queue pool sub-capacity parameters
	PoolSize uint8
	// maximum number of bytes allowed in a TCP word packet
	MaxPackSize uint16
	// maximum number of TCP connections allowed for a single service
	MaxConnect uint16
	// logs channel size
	LogSize uint16
}

// GlobalConfig global configuration parameters
var GlobalConfig *ZipperConfig

// configuration initialization function
func init() {
	fmt.Print(Banner)
	fmt.Println(zipperVersion)
	fmt.Println()
	GlobalConfig = &ZipperConfig{
		Name:        "Zipper-Server",
		Port:        8066,
		QueueSize:   200,
		PoolSize:    6,
		MaxPackSize: 1024,
		MaxConnect:  50,
		LogSize:     200,
	}
	// read json configuration file
	file, err := os.ReadFile("conf/zipper.json")
	if err != nil {
		log.Println("WARN " + err.Error())
		log.Println("WARN The json configuration file was not read, path = ./conf/zipper.json")
		log.Println("WARN The system will enable the default configuration")
		goto logs
	}
	// parse json file content
	err = json.Unmarshal(file, GlobalConfig)
	if err != nil {
		log.Println("ERROR json configuration file parsing exception")
		panic(err)
	}
logs:
	log.Println("INFO The port is " + strconv.Itoa(int(GlobalConfig.Port)))
	log.Println("INFO The QueueSize is " + strconv.Itoa(int(GlobalConfig.QueueSize)))
	log.Println("INFO The PoolSize is " + strconv.Itoa(int(GlobalConfig.PoolSize)))
	log.Println("INFO The MaxPackSize is " + strconv.Itoa(int(GlobalConfig.MaxPackSize)))
	log.Println("INFO The MaxConnect is " + strconv.Itoa(int(GlobalConfig.MaxConnect)))
	log.Println("INFO The LogSize is " + strconv.Itoa(int(GlobalConfig.LogSize)))
}
