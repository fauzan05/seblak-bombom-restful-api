package repository

import (
	"seblak-bombom-restful-api/internal/entity"

	"gorm.io/gorm"
)

type Repository[T any] struct {
	DB *gorm.DB
}

func (r *Repository[T]) Create(db *gorm.DB, entity *T) error {
	return db.Create(&entity).Error
}

func (r *Repository[T]) Update(db *gorm.DB, entity *T) error {
	return db.Save(entity).Error
}

func (r *Repository[T]) FindTokenByUserId(db *gorm.DB, token *T, userId int) error {
	return db.Where("user_id = ?", userId).First(&token).Error
}

func (r *Repository[T]) FindUserByToken(db *gorm.DB, user *T, token_code string) error {
	token := new(entity.Token)
	// temukan data user_id
	tokenWithUser := db.Where("token = ?", token_code).Joins("User").Find(&token).Error
	if tokenWithUser != nil {
		return tokenWithUser //return errornya
	}
	return db.Where("id = ?", token.UserId).Preload("Token").Preload("Addresses").Find(user).Error
}

func (r *Repository[T]) Delete(db *gorm.DB, entity *T) error {
	return db.Delete(entity).Error
}

func (r *Repository[T]) FindById(db *gorm.DB, entity *T) error {
	return db.First(&entity).Error
}

func(c *Repository[T]) DeleteToken(db *gorm.DB, entity *T, token string) *gorm.DB {
	result := db.Where("token = ?", token).Delete(&entity)
	return result
}

func (r *Repository[T]) FindByEmail(db *gorm.DB, entity *T, email string) error {
	return db.Where("email = ?", email).First(&entity).Error
}

func (r *Repository[T]) CheckEmailIsExists(db *gorm.DB, currentEmail string,requestEmail string) (int64, error) {
	var total int64
	err :=  db.Model(&entity.User{}).Where("email = ? AND email != ?", requestEmail, currentEmail).Count(&total).Error
	return total, err
}

func (r *Repository[T]) FindUserById(db *gorm.DB, entity *T, userId uint64) error {
	return db.Where("id = ?", userId).Preload("Token").Preload("Addresses").Find(&entity).Error
}

func (r *Repository[T]) FindUserByIdWithAddress(db *gorm.DB, entity *T, userId uint64) error {
	return db.Where("id = ?", userId).Preload("Addresses").Find(&entity).Error
}

func (r *Repository[T]) UserCountByEmail(db *gorm.DB, entity *T, email string) (int64, error) {
	var total int64
	err := db.Model(new(T)).Where("email = ?", email).Count(&total).Error
	return total, err
}

func (r *Repository[T]) DeleteAllAddressByUserId(db *gorm.DB, entity *T, userId uint64) *gorm.DB {
	result := db.Where("user_id = ?", userId).Delete(&entity)
	return result
}

func (r *Repository[T]) FindAll(db *gorm.DB, entities *[]T) error {
	return db.Find(&entities).Error
}
