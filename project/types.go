package project

import (
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type CertDetail struct {
	gorm.Model
	Domain string `gorm:"unique"`
	Value string
}

type DomainAltname struct {
	gorm.Model
	Domain string `gorm:"unique"`
	Value string
}

type UserInfo struct {
	gorm.Model
	Email string `gorm:"unique"`
	Value string
}

type Conf struct {
	LogFile string `yaml:"log_file"`

	Bolt struct {
		DataDir string `yaml:"data_dir"`
		DataFile string `yaml:"data_file"`
	}

	Psql struct {
		Host string `yaml:"host"`
		Port int `yaml:"port"`
		DatabaseName string `yaml:"db_name"`
		User string `yaml:"user"`
		Password string `yaml:"password"`
		SslMode string `yaml:"ssl_mode"`
	}

	Kube struct {
		ConfigFile string `yaml:"config_file"`
	}
}

func (c *Conf) GetConf() *Conf {
	yamlFile, err := ioutil.ReadFile("conf.yaml")
	if err != nil {
		log.Fatalf("yamlFile.Get err   #%v ", err)
	}

	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return c
}

func PrintLogMsg(message string) {
	log.Info("****************************************")
	log.Info(message)
	log.Info("****************************************")
}

func ArgumentError() {
	log.Fatal("No parameter specified: Use: print-bolt-data create-tables drop-tables migrate kube-migrate")
}
