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
	BoltDataDir string `yaml:"bolt_data_dir"`
	BoltDataFile string `yaml:"bolt_data_file"`
	PostgresHost string `yaml:"pg_host"`
	PostgresPort int `yaml:"pg_port"`
	PostgresDatabaseName string `yaml:"pg_db_name"`
	PostgresUser string `yaml:"pg_user"`
	PostgresPassword string `yaml:"pg_password"`
	PostgresSslMode string `yaml:"pg_ssl_mode"`
	LogFile string `yaml:"log_file"`
	KubeConfigFile string `yaml:"kube_config_file"`
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
