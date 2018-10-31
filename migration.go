package main

import (
	"flag"
	"github.com/boltdb/bolt"
	log "github.com/sirupsen/logrus"
	"os"
	"path"
)

func main() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
		TimestampFormat : "2006-01-02 15:04:05",
	})

	var workingDir, err = os.Getwd()
	if err != nil {
		log.Fatal("Cannot get current working directory.")
	}

	var dataDir string
	flag.StringVar(&dataDir, "data-dir", workingDir, "Data directory path")
	flag.Parse()

	dbPath := path.Join(dataDir, "data.db")

	db, err := bolt.Open(dbPath, 0600, nil)
	if err != nil {
		log.Fatalf("Error while creating bolt database file at %v: %v", dbPath, err)
	}
	defer db.Close()

	printKeyValuePairs(*db, getBucketNames(*db))
	log.Info("****************************************")
	log.Info("*************** Done *******************")
	log.Info("****************************************")
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

func printKeyValuePairs(db bolt.DB, bucketNames []string) {
	db.View(func(tx *bolt.Tx) error {
		for _, bucket := range bucketNames {
			log.Info("*******************************************")
			log.Infof("************** %v **************", bucket)
			log.Info("*******************************************")

			b := tx.Bucket([]byte(bucket))

			c := b.Cursor()

			for k, v := c.First(); k != nil; k, v = c.Next() {
				log.Infof("key=%s, value=%s\n", k, v)
			}
		}
		return nil
	})
}