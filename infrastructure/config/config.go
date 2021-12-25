package config

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type App struct {
	Name          string        `yaml:"name"`
	Domain        string        `yaml:"domain"`
	Log           Log           `yaml:"log"`
	DB            DB            `yaml:"db"`
	Sonic         Sonic         `yaml:"sonic"`
	ReviewService ReviewService `yaml:"review"`
	DictService   DictService   `yaml:"dict"`
}

type Log struct {
	Level         int  `yaml:"level"`
	DisableCaller bool `yaml:"disableCaller"`
}

type DB struct {
	Name              string        `yaml:"name"`
	URI               string        `yaml:"uri"`
	ConnectionTimeout time.Duration `yaml:"connectTimeout"`
	UserTable         TableConfig   `yaml:"userTable"`
	ReviewItemTable   TableConfig   `yaml:"reviewItemTable"`
	InspirationTable  TableConfig   `yaml:"inspirationTable"`
	WordTable         TableConfig   `yaml:"wordTable"`
	DictTable         TableConfig   `yaml:"dictTable"`
}

type Sonic struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
	Pwd  string `yaml:"pwd"`
}

type TableConfig struct {
	Name string `yaml:"name"`
}

type ReviewService struct {
	Mail Mail `yaml:"mail"`
}

type DictService struct {
	Langs   []string `yaml:"langs"`
	Schemas []string `yaml:"schemas"`
}

type Mail struct {
	Smtp Smtp `yaml:"smtp"`
}

type Smtp struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Sender   string `yaml:"sender"`
	Email    string `yaml:"email"`
	Password string `yaml:"password"`
}

const (
	EmptyString       = ""
	EnvAppConfigURL   = "APP_CONFIG_URL"
	EnvAppConfigPath  = "APP_CONFIG_PATH"
	DefaultConfigPath = "config.yaml"
)

var appConfigMutex sync.Mutex
var appConfig App

func Get() App {
	return appConfig
}

func init() {
	appConfigMutex.Lock()
	defer appConfigMutex.Unlock()
	// load readConfig file
	configPath := os.Getenv(EnvAppConfigURL)
	if configPath == EmptyString {
		configPath = os.Getenv(EnvAppConfigPath)
	}
	if configPath == EmptyString {
		configPath = DefaultConfigPath
	}

	//parse readConfig file
	var err error
	var fileBytes []byte
	if fileBytes, err = readConfig(configPath); err != nil {
		fmt.Printf("Load readConfig File Error: %v\n", err)
		return
	}

	//unmarshal readConfig file
	if err = unmarshal(fileBytes, &appConfig, false); err != nil {
		fmt.Printf("unmarshal readConfig File Error: %v\n", err)
		return
	}
}

// internal

func readConfig(location string) (bytes []byte, err error) {
	if strings.HasPrefix(location, "http") {
		return remoteConfig(location)
	}
	return localConfig(location)
}

func localConfig(filePath string) (bytes []byte, err error) {
	fileBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return fileBytes, nil
}

func remoteConfig(url string) (bytes []byte, err error) {
	if url == "" {
		return nil, errors.New("can't get readConfig url")
	}
	result, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadAll(result.Body)
	if err != nil {
		return nil, err
	}
	return b, err
}

func unmarshal(in []byte, out interface{}, isStrict bool) (err error) {
	if in == nil {
		err = error(errors.New("can't unmarshal empty byte"))
		return err
	}
	if isStrict {
		err = yaml.UnmarshalStrict(in, out)
		if err != nil {
			return err
		}
	} else {
		err = yaml.Unmarshal(in, out)
		if err != nil {
			return err
		}
	}
	return nil
}
