package main

import (
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/liquidweb/cert-manager-migration/project"
	log "github.com/sirupsen/logrus"
	"os"
)

func main() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
		TimestampFormat : "2006-01-02 15:04:05",
	})

	var conf project.Conf
	conf.GetConf()

	if conf.LogFile != "" {
		f, err := os.OpenFile(conf.LogFile, os.O_WRONLY|os.O_CREATE, 0755)
		if err != nil {
			log.Fatalf("Error while reading in file: %v", err)
		}
		log.SetOutput(f)
	}

	project.PrintLogMsg("Starting Cert Manager Migration")
	log.Info()

	if len(os.Args) <= 1 {
		project.ArgumentError()
	}

	switch os.Args[1] {
	case "print-bolt-data":
		project.PrintBoltData(conf)
	case "create-tables":
		project.CreateTables(conf)
	case "drop-tables":
		project.DropTables(conf)
	case "migrate":
		project.Migrate(conf)
	case "kube-test":
		project.KubeTest(conf)
	default:
		project.ArgumentError()
	}

	log.Info()
	project.PrintLogMsg("Done")
}

