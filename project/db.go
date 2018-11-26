package project

import (
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"path"
)

func openBoltDb(conf Conf) *bolt.DB {
	dbPath := path.Join(conf.Bolt.DataDir, conf.Bolt.DataFile)

	boltDb, err := bolt.Open(dbPath, 0600, nil)
	if err != nil {
		log.Fatalf("Error while creating bolt database file at %v: %v", dbPath, err)
	}

	return boltDb
}

func openPostgresDb(conf Conf) *gorm.DB {
	args := fmt.Sprintf("host=%v port=%d user=%v dbname=%v password=%v sslmode=%v", conf.Psql.Host, conf.Psql.Port, conf.Psql.User, conf.Psql.DatabaseName, conf.Psql.Password, conf.Psql.SslMode)
	postgresDb, err := gorm.Open("postgres", args)
	if err != nil {
		log.Fatalf("Error while connecting to postgres database: %v", err)
	}
	return postgresDb
}

func createTables(db gorm.DB) {
	PrintLogMsg("Creating Tables: CERT_DETAILS, DOMAIN_ALTNAMES, USER_INFOS")
	log.Info()

	err := migrateTables(db)
	if err != nil {
		log.Fatalf("Error while creating database tables: %v", err)
	}
}

func dropTables(db gorm.DB) {
	PrintLogMsg("Dropping Tables: CERT_DETAILS, DOMAIN_ALTNAMES, USER_INFOS")
	log.Info()

	err := rollbackTables(db)
	if err != nil {
		log.Fatalf("Error while dropping database tables: %v", err)
	}
}

func doMigration(postgresDb gorm.DB, boltDb bolt.DB) {
	PrintLogMsg("Migrating Tables")
	log.Info()

	err := boltDb.View(func(tx *bolt.Tx) error {
		PrintLogMsg("Migrating cert-details into CERT_DETAILS\n")
		b := tx.Bucket([]byte("cert-details"))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			log.Infof("Migrating Domain: %s", k)
			domain := string(k[:])
			value := string(v[:])
			postgresDb.Create(&CertDetail{Domain: domain, Value: value})
		}

		log.Info()
		PrintLogMsg("Migrating domain-altnames into DOMAIN_ALTNAMES\n")
		b = tx.Bucket([]byte("domain-altnames"))
		c = b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			log.Infof("Migrating Domain: %s", k)
			domain := string(k[:])
			value := string(v[:])
			postgresDb.Create(&DomainAltname{Domain: domain, Value: value})
		}

		log.Info()
		PrintLogMsg("Migrating user-info into USER_INFOS\n")
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
	if err != nil {
		log.Fatalf("Error viewing bolt data: %v", err)
	}
}

func getBucketNames(db bolt.DB) []string {
	var buckets []string
	err := db.View(func(tx *bolt.Tx) error {
		err := tx.ForEach(func(name []byte, b *bolt.Bucket) error {
			buckets = append(buckets, string(name))
			return nil
		})
		if err != nil {
			log.Fatalf("Error grabbing bolt data: %v", err)
		}
		return nil
	})
	if err != nil {
		log.Fatalf("Error viewing bolt data: %v", err)
	}
	return buckets
}

func printBoltKeyValuePairs(db bolt.DB, bucketNames []string) {
	err := db.View(func(tx *bolt.Tx) error {
		for _, bucket := range bucketNames {
			PrintLogMsg(fmt.Sprintf("%v", bucket))

			b := tx.Bucket([]byte(bucket))

			c := b.Cursor()

			for k, v := c.First(); k != nil; k, v = c.Next() {
				log.Infof("key=%s, value=%s\n", k, v)
			}
		}
		return nil
	})
	if err != nil {
		log.Fatalf("Error viewing bolt data: %v", err)
	}
}
