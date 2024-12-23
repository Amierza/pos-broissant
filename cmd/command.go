package cmd

import (
	"log"
	"os"

	"github.com/Amierza/pos-broissant/migrations"
	"gorm.io/gorm"
)

func Command(db *gorm.DB) {
	migrate := false
	seed := false
	rollback := false

	for _, arg := range os.Args[1:] {
		if arg == "--migrate" {
			migrate = true
		}
		if arg == "--seed" {
			seed = true
		}
		if arg == "--rollback" {
			rollback = true
		}
	}

	if migrate {
		if err := migrations.Migrate(db); err != nil {
			log.Fatal("error migration: %v", err)
		}
		log.Println("migration completed successfully!")
	}

	if seed {
		if err := migrations.Seeder(db); err != nil {
			log.Fatal("error migration seed: %v", err)
		}
		log.Println("seeder completed successfully!")
	}

	if rollback {
		if err := migrations.Rollback(db); err != nil {
			log.Fatalf("error rolling back tables: %v", err)
		}
		log.Println("rollback completed successfully!")
	}
}
