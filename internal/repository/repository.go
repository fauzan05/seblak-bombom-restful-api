package repository

import "gorm.io/gorm"

type Repository[T any] struct {
	DB *gorm.DB
}

func (r *Repository[T]) Create(db *gorm.DB, entity *T) error {
	return db.Create(entity).Error
}

func (r *Repository[T]) Update(db *gorm.DB, entity *T) error {
	return db.Save(entity).Error
}

func (r *Repository[T]) FindTokenByUserId(db *gorm.DB, token *T, userId int) error {
	return db.Where("user_id = ?", userId).First(token).Error
}

func (r *Repository[T]) FindUserByToken(db *gorm.DB, token *T, token_code string) error {
	return db.Where("token = ?", token_code).Find(token).Error
}

func (r *Repository[T]) Delete(db *gorm.DB, entity *T) error {
	return db.Delete(entity).Error
}

func (r *Repository[T]) FindByEmail(db *gorm.DB, entity *T, email string) error {
	return db.Where("email = ?", email).First(entity).Error
}

func (r *Repository[T]) FindUserById(db *gorm.DB, entity *T, userId uint64) error {
	return db.Where("id = ?", userId).First(entity).Error
}

func (r *Repository[T]) FindUserByIdWithAddress(db *gorm.DB, entity *T, userId uint64) error {
	return db.Where("id = ?", userId).Preload("Addresses").Find(entity).Error
}

func (r *Repository[T]) UserCountByEmail(db *gorm.DB, entity *T, email string) (int64, error) {
	var total int64
	err := db.Model(new(T)).Where("email = ?", email).Count(&total).Error
	return total, err
}