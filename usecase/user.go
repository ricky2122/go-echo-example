package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/ricky2122/go-echo-example/domain"
)

var ErrUserAlreadyExists = errors.New("user already exists")

type SignUpUseCaseInput struct {
	Name     string
	Password string
	Email    string
	BirthDay time.Time
}

type SignUpUseCaseOutput struct {
	ID   int
	Name string
}

type IUserRepository interface {
	IsExist(ctx context.Context, name string) (bool, error)
	Create(ctx context.Context, newUser domain.User) (*domain.User, error)
}

type UserUseCase struct {
	ur IUserRepository
}

func NewUserUseCase(ur IUserRepository) *UserUseCase {
	return &UserUseCase{ur: ur}
}

func (uc *UserUseCase) SignUp(input SignUpUseCaseInput) (*SignUpUseCaseOutput, error) {
	user := domain.NewUser(input.Name, input.Password, input.Email, input.BirthDay)

	ctx := context.Background()

	// check if user already exists
	isExist, err := uc.ur.IsExist(ctx, user.GetName())
	if err != nil {
		return nil, err
	}
	if isExist {
		return nil, ErrUserAlreadyExists
	}

	// create user
	createdUser, err := uc.ur.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	// response
	output := &SignUpUseCaseOutput{
		ID:   int(createdUser.GetID()),
		Name: createdUser.GetName(),
	}

	return output, nil
}
