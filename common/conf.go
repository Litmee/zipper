package common

import (
	"encoding/json"
	"fmt"
	"os"
	"zipper/logger"
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
		logger.OutLog("The port is 8066", logger.INFO)
		logger.OutLog("The QueueSize is 200", logger.INFO)
		logger.OutLog("The PoolSize is 6", logger.INFO)
		logger.OutLog("The MaxPackSize is 1024", logger.INFO)
		logger.OutLog("The MaxConnect is 30", logger.INFO)
		return
	}
	// parse json file content
	err = json.Unmarshal(file, GlobalConfig)
	if err != nil {
		logger.OutLog("json configuration file parsing exception", logger.ERROR)
		panic(err)
	}
}
