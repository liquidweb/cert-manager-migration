package project

import (
	log "github.com/sirupsen/logrus"
)

func PrintBoltData(conf Conf) {
	PrintLogMsg("Print Bolt Data")
	log.Info()
	var boltDb = openBoltDb(conf)

	printBoltKeyValuePairs(*boltDb, getBucketNames(*boltDb))

	defer boltDb.Close()
}

func CreateTables(conf Conf) {
	PrintLogMsg("Create Tables")
	log.Info()
	var postgresDb = openPostgresDb(conf)

	createTables(*postgresDb)

	defer postgresDb.Close()
}

func DropTables(conf Conf) {
	PrintLogMsg("Drop Tables")
	log.Info()
	var postgresDb = openPostgresDb(conf)

	dropTables(*postgresDb)

	defer postgresDb.Close()
}

func Migrate(conf Conf) {
	PrintLogMsg("BoltDB to PostgreSQL Migration")
	log.Info()
	var postgresDb = openPostgresDb(conf)
	var boltDb = openBoltDb(conf)

	createTables(*postgresDb)
	doMigration(*postgresDb, *boltDb)

	defer postgresDb.Close()
	defer boltDb.Close()
}

func KubeTest(conf Conf) {
	PrintLogMsg("Kubernetes Test")
	log.Info()

	kubeTest(conf)
}