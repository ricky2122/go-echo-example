package usecase

import (
	"errors"
	"time"

	"github.com/ricky2122/go-echo-example/domain"
)

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
	IsExist(name string) (bool, error)
	Create(newUser domain.User) (*domain.User, error)
}

type UserUseCase struct {
	ur IUserRepository
}

func NewUserUseCase(ur IUserRepository) *UserUseCase {
	return &UserUseCase{ur: ur}
}

func (uc *UserUseCase) SignUp(input SignUpUseCaseInput) (*SignUpUseCaseOutput, error) {
	user := domain.NewUser(input.Name, input.Password, input.Email, input.BirthDay)

	// check if user already exists
	isExist, err := uc.ur.IsExist(user.GetName())
	if err != nil {
		return nil, err
	}
	if isExist {
		return nil, errors.New("user already exists")
	}

	// create user
	createdUser, err := uc.ur.Create(user)
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
