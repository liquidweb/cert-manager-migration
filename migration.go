package main

import (
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path"
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

type conf struct {
	BoltDataDir string `yaml:"bolt_data_dir"`
	BoltDataFile string `yaml:"bolt_data_file"`
	PostgresHost string `yaml:"pg_host"`
	PostgresPort int `yaml:"pg_port"`
	PostgresDatabaseName string `yaml:"pg_db_name"`
	PostgresUser string `yaml:"pg_user"`
	PostgresPassword string `yaml:"pg_password"`
	PostgresSslMode string `yaml:"pg_ssl_mode"`
	LogFile string `yaml:"log_file"`
}

func main() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
		TimestampFormat : "2006-01-02 15:04:05",
	})

	var conf conf
	conf.getConf()

	if conf.LogFile != "" {
		f, err := os.OpenFile(conf.LogFile, os.O_WRONLY|os.O_CREATE, 0755)
		if err != nil {
			log.Fatalf("Error while reading in file: %v", err)
		}
		log.SetOutput(f)
	}

	printLogMsg("Starting Cert Manager Migration")
	log.Info()

	if len(os.Args) <= 1 {
		argumentError()
	}

	switch os.Args[1] {
	case "print-bolt-data":
		printLogMsg("Print Bolt Data")
		log.Info()
		var boltDb = openBoltDb(conf)
		printBoltKeyValuePairs(*boltDb, getBucketNames(*boltDb))
		defer boltDb.Close()
	case "create-tables":
		printLogMsg("Create Tables")
		log.Info()
		var postgresDb = openPostgresDb(conf)
		createTables(*postgresDb)
		defer postgresDb.Close()
	case "drop-tables":
		printLogMsg("Drop Tables")
		log.Info()
		var postgresDb = openPostgresDb(conf)
		dropTables(*postgresDb)
		defer postgresDb.Close()
	case "migrate":
		printLogMsg("BoltDB to PostgreSQL Migration")
		log.Info()
		var postgresDb = openPostgresDb(conf)
		var boltDb = openBoltDb(conf)
		createTables(*postgresDb)
		doMigration(*postgresDb, *boltDb)
		defer postgresDb.Close()
		defer boltDb.Close()
	default:
		argumentError()
	}

	log.Info()
	printLogMsg("Done")
}

func openBoltDb(conf conf) *bolt.DB {
	dbPath := path.Join(conf.BoltDataDir, conf.BoltDataFile)

	boltDb, err := bolt.Open(dbPath, 0600, nil)
	if err != nil {
		log.Fatalf("Error while creating bolt database file at %v: %v", dbPath, err)
	}

	return boltDb
}

func openPostgresDb(conf conf) *gorm.DB {
	args := fmt.Sprintf("host=%v port=%d user=%v dbname=%v password=%v sslmode=%v", conf.PostgresHost, conf.PostgresPort, conf.PostgresUser, conf.PostgresDatabaseName, conf.PostgresPassword, conf.PostgresSslMode)
	postgresDb, err := gorm.Open("postgres", args)
	if err != nil {
		log.Fatalf("Error while connecting to postgres database: %v", err)
	}
	return postgresDb
}

func createTables(db gorm.DB) {
	printLogMsg("Creating Tables: CERT_DETAILS, DOMAIN_ALTNAMES, USER_INFOS")
	log.Info()

	db.AutoMigrate(&CertDetail{})
	db.AutoMigrate(&DomainAltname{})
	db.AutoMigrate(&UserInfo{})
}

func dropTables(db gorm.DB) {
	printLogMsg("Dropping Tables: CERT_DETAILS, DOMAIN_ALTNAMES, USER_INFOS")
	log.Info()

	db.DropTable(&CertDetail{})
	db.DropTable(&DomainAltname{})
	db.DropTable(&UserInfo{})
}

func doMigration(postgresDb gorm.DB, boltDb bolt.DB) {
	printLogMsg("Migrating Tables")
	log.Info()

	boltDb.View(func(tx *bolt.Tx) error {
		printLogMsg("Migrating cert-details into CERT_DETAILS\n")
		b := tx.Bucket([]byte("cert-details"))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			log.Infof("Migrating Domain: %s", k)
			domain := string(k[:])
			value := string(v[:])
			postgresDb.Create(&CertDetail{Domain: domain, Value: value})
		}

		log.Info()
		printLogMsg("Migrating domain-altnames into DOMAIN_ALTNAMES\n")
		b = tx.Bucket([]byte("domain-altnames"))
		c = b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			log.Infof("Migrating Domain: %s", k)
			domain := string(k[:])
			value := string(v[:])
			postgresDb.Create(&DomainAltname{Domain: domain, Value: value})
		}

		log.Info()
		printLogMsg("Migrating user-info into USER_INFOS\n")
		b = tx.Bucket([]byte("user-info"))
		c = b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			log.Infof("Migrating Email: %s", k)
			email := string(k[:])
			value := string(v[:])
			postgresDb.Create(&UserInfo{Email: email, Value: value})
		}
		return nil
	})
}

func getBucketNames(db bolt.DB) []string {
	var buckets []string
	db.View(func(tx *bolt.Tx) error {
		tx.ForEach(func(name []byte, b *bolt.Bucket) error {
			buckets = append(buckets, string(name))
			return nil
		})
		return nil
	})

	return buckets
}

func printBoltKeyValuePairs(db bolt.DB, bucketNames []string) {
	db.View(func(tx *bolt.Tx) error {
		for _, bucket := range bucketNames {
			printLogMsg(fmt.Sprintf("%v", bucket))

			b := tx.Bucket([]byte(bucket))

			c := b.Cursor()

			for k, v := c.First(); k != nil; k, v = c.Next() {
				log.Infof("key=%s, value=%s\n", k, v)
			}
		}
		return nil
	})
}

func (c *conf) getConf() *conf {
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

func printLogMsg(message string) {
	log.Info("****************************************")
	log.Info(message)
	log.Info("****************************************")
}

func argumentError() {
	log.Fatal("No parameter specified: Use: print-bolt-data create-tables drop-tables migrate")
}