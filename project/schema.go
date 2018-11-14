package project

import (
	"github.com/jinzhu/gorm"
	"gopkg.in/gormigrate.v1"
)

func schemaMigration(db gorm.DB) *gormigrate.Gormigrate {
	m := gormigrate.New(&db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID: "initial",
			Migrate: func(tx *gorm.DB) error {
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

				return tx.AutoMigrate(&CertDetail{}, &DomainAltname{}, &UserInfo{}).Error
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.DropTable("cert_details", "domain_altnames", "user_infos").Error
			},
		},
	})

	return m
}

func migrateTables(db gorm.DB) error {
	m := schemaMigration(db)

	err := m.Migrate()
	if err != nil {
		return err
	}

	return nil
}

func rollbackTables(db gorm.DB) error {
	m := schemaMigration(db)

	err := m.RollbackTo("")
	if err != nil {
		return err
	}

	return nil
}