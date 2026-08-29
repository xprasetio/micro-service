package repository

import (
	"context"
	"errors"
	"log"
	"user-service/internal/core/domain/entity"
	"user-service/internal/core/domain/model"

	"gorm.io/gorm"
)

type RoleRepositoryInterface interface {
	GetAll(ctx context.Context, search string) ([]entity.RoleEntity, error)
	GetByID(ctx context.Context, id int64) (*entity.RoleEntity, error)
	Create(ctx context.Context, req entity.RoleEntity) error
	Update(ctx context.Context, req entity.RoleEntity) error
	Delete(ctx context.Context, id int64) error
}

type roleRepository struct {
	db *gorm.DB
}

// Create implements [RoleRepositoryInterface].
func (r *roleRepository) Create(ctx context.Context, req entity.RoleEntity) error {
	modelRole := model.Role{
		Name: req.Name,
	}
	if err := r.db.WithContext(ctx).Create(&modelRole).Error; err != nil {
		log.Printf("[RoleRepository-1] Create: %v", err)
		return err
	}
	return nil
}

// Delete implements [RoleRepositoryInterface].
func (r *roleRepository) Delete(ctx context.Context, id int64) error {
	modelRole := model.Role{}
	if err := r.db.WithContext(ctx).Where("id = ?", id).Preload("Users").First(&modelRole).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.New("404")
			log.Printf("[RoleRepository-1] Delete: %v", err)
			return err
		}
		log.Printf("[RoleRepository-1] Delete: %v", err)
		return err
	}

	if len(modelRole.Users) > 0 {
		err := errors.New("400")
		log.Printf("[RoleRepository-2] Delete: %v", err)
		return err
	}

	if err := r.db.WithContext(ctx).Delete(&modelRole).Error; err != nil {
		log.Printf("[RoleRepository-2] Delete: %v", err)
		return err
	}
	return nil
}

// GetAll implements [RoleRepositoryInterface].
func (r *roleRepository) GetAll(ctx context.Context, search string) ([]entity.RoleEntity, error) {
	modelRoles := []model.Role{}

	if err := r.db.WithContext(ctx).Where("name LIKE ?", "%"+search+"%").Find(&modelRoles).Error; err != nil {
		log.Printf("[RoleRepository-1] GetAll: %v", err)
		return nil, err
	}
	if len(modelRoles) == 0 {
		err := errors.New("role not found")
		log.Printf("[RoleRepository-2] GetAll: %v", err)
		return nil, err
	}
	entityRole := []entity.RoleEntity{}
	for _, role := range modelRoles {
		entityRole = append(entityRole, entity.RoleEntity{
			ID:   role.ID,
			Name: role.Name,
		})
	}
	return entityRole, nil

}

// GetByID implements [RoleRepositoryInterface].
func (r *roleRepository) GetByID(ctx context.Context, id int64) (*entity.RoleEntity, error) {
	modelRole := model.Role{}
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&modelRole).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.New("role not found")
			log.Printf("[RoleRepository-1] GetByID: %v", err)
			return nil, err
		}
		log.Printf("[RoleRepository-1] GetByID: %v", err)
		return nil, err
	}

	return &entity.RoleEntity{
		ID:   modelRole.ID,
		Name: modelRole.Name,
	}, nil
}

// Update implements [RoleRepositoryInterface].
func (r *roleRepository) Update(ctx context.Context, req entity.RoleEntity) error {
	modelRole := model.Role{}
	if err := r.db.WithContext(ctx).Where("id = ?", req.ID).First(&modelRole).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.New("role not found")
			log.Printf("[RoleRepository-1] Update: %v", err)
			return err
		}
		log.Printf("[RoleRepository-1] Update: %v", err)
		return err
	}

	modelRole.Name = req.Name
	if err := r.db.WithContext(ctx).Save(&modelRole).Error; err != nil {
		log.Printf("[RoleRepository-2] Update: %v", err)
		return err
	}
	return nil
}

func NewRoleRepository(db *gorm.DB) RoleRepositoryInterface {
	return &roleRepository{
		db: db,
	}
}
