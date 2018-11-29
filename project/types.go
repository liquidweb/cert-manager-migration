package project

import (
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

/*
	PostgreSQL Database Schema
 */
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

/*
	Configuration File Definition
 */
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
		SourceConfigFile string `yaml:"src_config_file"`
		DestConfigFile string `yaml:"dest_config_file"`
	}
}

/*
	Certificate Custom Resource Definition
 */
type Certificate struct {
	metav1.TypeMeta `json:",inline"`
	Metadata        metav1.ObjectMeta `json:"metadata"`
	Spec            CertificateSpec   `json:"spec"`
	Status          CertificateStatus `json:"status,omitempty"`
}

type CertificateStatus struct {
	Provisioned string `json:"provisioned,omitempty"`
	CreatedDate string `json:"created,omitempty"`
	ExpiresDate string `json:"expires,omitempty"`
	ErrorMsg    string `json:"error_msg,omitempty"`
	ErrorDate   string `json:"error_date,omitempty"`
}

type CertificateSpec struct {
	Domain     string   `json:"domain"`
	Provider   string   `json:"provider"`
	Email      string   `json:"email"`
	SecretName string   `json:"secretName"`
	AltNames   []string `json:"altNames"`
}

type CertificateList struct {
	metav1.TypeMeta `json:",inline"`
	Metadata        metav1.ListMeta `json:"metadata"`
	Items           []Certificate   `json:"items"`
}

func ( Certificate) DeepCopyObject() runtime.Object {
	log.Warn("Certificate DeepCopyObject Not Implemented")
	return nil
}

func ( CertificateList) DeepCopyObject() runtime.Object {
	log.Warn("CertificateList DeepCopyObject Not Implemented")
	return nil
}

type KubeClient struct {
	KubeConfig *rest.Config
	Client *kubernetes.Clientset
	CertClient *rest.RESTClient
}

/*
	Utility Functions
 */
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
