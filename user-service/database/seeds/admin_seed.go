package seeds

import (
	"user-service/internal/core/domain/model"
	"user-service/utils/conv"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func SeedAdmin(db *gorm.DB) {
	bytes, err := conv.HashPassword("admin123")
	if err != nil {
		log.Error().Err(err).Msg("Failed to hash admin password")
		return
	}

	// Create the admin role
	existingRole := model.Role{
		Name: "Super Admin",
	}
	if err := db.FirstOrCreate(&existingRole, model.Role{Name: "Super Admin"}).Error; err != nil {
		log.Error().Err(err).Msg("Failed to create admin role")
		return
	}

	// Create the admin role
	admin := model.User{
		Name:       "Super Admin",
		Email:      "admin@admin.com",
		Password:   bytes,
		IsVerified: true,
		Roles:      []model.Role{existingRole},
	}
	if err := db.FirstOrCreate(&admin, model.User{Email: "admin@admin.com"}).Error; err != nil {
		log.Error().Err(err).Msg("Failed to create admin role")
		return
	}

	log.Info().Msg("Admin role seeded successfully")
}
