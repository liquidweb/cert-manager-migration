package main

import (
	"flag"
	"github.com/boltdb/bolt"
	"log"
	"os"
	"path"
)

func main() {
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
	log.Printf("Done")
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
			log.Printf("************** %v **************", bucket)

			b := tx.Bucket([]byte(bucket))

			c := b.Cursor()

			for k, v := c.First(); k != nil; k, v = c.Next() {
				log.Printf("key=%s, value=%s\n", k, v)
			}

			log.Printf("*********************************")
		}
		return nil
	})
}