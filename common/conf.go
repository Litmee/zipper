package common

import (
	"encoding/json"
	"fmt"
	"os"
	"zipper/logger"
)

// framework version information
var zipperVersion = "	Zipper v1.01, go 1.19"

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
		MaxConnect:  30,
	}
	// read json configuration file
	file, err := os.ReadFile("conf/zipper.json")
	if err != nil {
		logger.OutLog(err.Error(), logger.WARN)
		logger.OutLog("The json configuration file was not read, path = ./conf/zipper.json", logger.WARN)
		logger.OutLog("The system will enable the default configuration", logger.WARN)
		return
	}
	// parse json file content
	err = json.Unmarshal(file, GlobalConfig)
	if err != nil {
		logger.OutLog("json configuration file parsing exception", logger.ERROR)
		panic(err)
	}
}
