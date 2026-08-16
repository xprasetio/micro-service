package seeds

import (
	"user-service/internal/core/domain/model"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func SeedRole(db *gorm.DB) {
	roles := []model.Role{
		{Name: "Super Admin"},
		{Name: "Customer"},
	}

	for _, role := range roles {
		if err := db.FirstOrCreate(&role, model.Role{Name: role.Name}).Error; err != nil {
			log.Error().Err(err).Msg("Failed to create role")
		} else {
			log.Info().Msg("Role seeded successfully")
		}
	}
}
