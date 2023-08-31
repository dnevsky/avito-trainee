package migrations

import (
	"avito-trainee/internal/models"
	"fmt"
	"log"

	"gorm.io/gorm"
)

type Migration interface {
	Migrate(*gorm.DB) error
}

func RunMigrations(ormDB *gorm.DB) {
	if err := ormDB.AutoMigrate(
		models.Segment{},
		models.User{},
		models.SegmentHistory{},
		models.UserSegment{},
	); err != nil {
		log.Fatalf(fmt.Sprintf("migraion process was failed: %s", err))
	}

	if err := ormDB.SetupJoinTable(&models.User{}, "Segments", &models.UserSegment{}); err != nil {
		log.Fatalf(fmt.Sprintf("migration process was failed: %s", err))
	}
}
