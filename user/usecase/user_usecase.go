package usecase

import (
	"context"
	"errors"
	"fmt"

	"api-mygram-go/domain"
	"api-mygram-go/helpers"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

type userUseCase struct {
	userRepository domain.UserRepository
}

func NewUserUseCase(userRepository domain.UserRepository) *userUseCase {
	return &userUseCase{userRepository}
}

func (userUseCase *userUseCase) Register(ctx context.Context, user *domain.User) (err error) {
	ID, _ := gonanoid.New(16)
	user.ID = fmt.Sprintf("user-%s", ID)

	if err = userUseCase.userRepository.Create(ctx, user); err != nil {
		return err
	}

	return
}

func (userUseCase *userUseCase) Login(ctx context.Context, user *domain.User) (err error) {
	password := user.Password

	if err = userUseCase.userRepository.GetUserByEmail(ctx, user); err != nil {
		return err
	}
	if isValid := helpers.Compare([]byte(user.Password), []byte(password)); !isValid {
		return errors.New("the credential you entered are wrong")
	}

	return
}

func (userUseCase *userUseCase) GetAllUsers(ctx context.Context, users *[]domain.User) (err error) {
	if err = userUseCase.userRepository.GetAllUsers(ctx, users); err != nil {
		return err
	}

	return nil
}

func (userUseCase *userUseCase) Update(ctx context.Context, user domain.User) (u domain.User, err error) {
	if u, err = userUseCase.userRepository.Update(ctx, user); err != nil {
		return u, err
	}

	return u, nil
}

func (userUseCase *userUseCase) Delete(ctx context.Context, id string) (err error) {
	if err = userUseCase.userRepository.Delete(ctx, id); err != nil {
		return err
	}

	return
}
