package seeds

import (
	"gonews/internal/core/domain/model"
	"gonews/lib/conv"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func SeedRoles(db *gorm.DB) {
	bytes, err := conv.HashPassword("admin123")
	if err != nil {
		log.Fatal().Err(err).Msg("Error when creating hash password")
	}

	admin := model.User{
		Name: "Admin",
		Email: "admin@mail.com",
		Password: string(bytes),
	}

	if err := db.FirstOrCreate(&admin, model.User{Email: "admin@mail.com"}).Error; err != nil {
		log.Fatal().Err(err).Msg("Error seeding")
	} else {
		log.Info().Msg("Seeded successfuly")
	}
}