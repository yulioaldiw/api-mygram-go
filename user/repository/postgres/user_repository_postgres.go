package repository

import (
	"api-mygram-go/domain"
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{db}
}

func (userRepository *userRepository) Create(ctx context.Context, user *domain.User) (err error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err = userRepository.db.WithContext(ctx).Create(&user).Error; err != nil {
		return err
	}

	return
}

func (userRepository *userRepository) GetUserByEmail(ctx context.Context, user *domain.User) (err error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err = userRepository.db.WithContext(ctx).Where("email = ?", user.Email).Take(&user).Error; err != nil {
		return errors.New("the email you entered are not registered")
	}

	return
}

func (userRepository *userRepository) GetUserByID(ctx context.Context, user *domain.User) (err error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err = userRepository.db.WithContext(ctx).Where("id = ?", user.ID).Take(&user).Error; err != nil {
		return errors.New("the id you entered are not registered")
	}

	return
}

func (userRepository *userRepository) GetAllUsers(ctx context.Context, users *[]domain.User) (err error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)

	defer cancel()

	if err = userRepository.db.WithContext(ctx).
		Select("id", "username", "email", "age", "profile_image_url").Find(&users).Error; err != nil {
		return err
	}

	return
}

func (userRepository *userRepository) Update(ctx context.Context, user domain.User, userID string) (u domain.User, err error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)

	defer cancel()

	u = domain.User{}

	if err = userRepository.db.WithContext(ctx).First(&u, "id = ?", userID).Error; err != nil {
		return u, err
	}

	if err = userRepository.db.WithContext(ctx).Model(&u).Updates(user).Error; err != nil {
		return u, err
	}

	return u, nil
}

func (userRepository *userRepository) Delete(ctx context.Context, id string) (err error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)

	defer cancel()

	if err = userRepository.db.WithContext(ctx).First(&domain.User{}, &id).Error; err != nil {
		return err
	}

	if err = userRepository.db.WithContext(ctx).Where("user_id = ?", id).Delete(&domain.SocialMedia{}).Error; err != nil {
		return err
	}

	if err = userRepository.db.WithContext(ctx).Delete(&domain.User{}, &id).Error; err != nil {
		return err
	}

	return
}
