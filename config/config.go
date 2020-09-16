package config

import (
	"flag"
	"os"
	"strconv"
	"sync"

	"github.com/cyrildever/treee/common/logger"
)

//--- TYPES

// Config ...
type Config struct {
	HTTPPort       string
	Host           string
	InitPrime      uint64
	IndexPath      string
	UsePersistence bool
}

var singleton *Config
var once sync.Once
var err error
var isTestEnv bool

//--- METHODS

// IsTestEnvironment ...
func (c *Config) IsTestEnvironment() bool {
	return isTestEnv
}

func (c *Config) populateWithEnv() {
	setString("HTTP_PORT", &c.HTTPPort)
	setString("HOST", &c.Host)
	setUintOrPanic("INIT_PRIME", &c.InitPrime)
	setString("INDEX_PATH", &c.IndexPath)
	setBoolean("PERMANENT_INDEX", &c.UsePersistence)
}

//--- FUNCTIONS

// InitConfig ...
func InitConfig(isTest bool) (*Config, error) {
	isTestEnv = isTest
	return GetConfig()
}

// GetConfig ...
func GetConfig() (*Config, error) {
	once.Do(func() {
		singleton = &Config{}
		httpPort := flag.String("t.port", "7000", "HTTP port number")
		host := flag.String("t.host", "0.0.0.0", "Host address")
		indexPath := flag.String("t.file", "", "File path to an existing index")
		initPrime := flag.String("t.init", "0", "Initial prime number to use for the index")
		usePersistence := flag.Bool("t.persist", true, "Activate persistence")

		flag.Parse()

		singleton.HTTPPort = *httpPort
		singleton.Host = *host
		singleton.IndexPath = *indexPath
		if *initPrime != "0" {
			p, e := strconv.ParseUint(*initPrime, 10, 64)
			if e != nil {
				err = e
				return
			}
			singleton.InitPrime = p
		}
		singleton.UsePersistence = *usePersistence

		singleton.populateWithEnv()
	})
	return singleton, err
}

func setBoolean(envName string, shouldChange *bool) {
	str := os.Getenv(envName)
	if v, err := strconv.ParseBool(str); err == nil {
		*shouldChange = v
	}
}

func setString(envName string, shouldChange *string) {
	str := os.Getenv(envName)
	if str != "" {
		*shouldChange = str
	}
}

func setUintOrPanic(envName string, shouldChange *uint64) {
	str := os.Getenv(envName)
	if str != "" {
		i, err := strconv.ParseUint(str, 10, 64)
		if err == nil {
			(*shouldChange) = i
		} else {
			log := logger.Init("config", "Environment variables")
			log.Error("You passed " + str + " for " + envName + ". A Uint64 is mandatory for this data!")
			os.Exit(1)
		}
	}
}
