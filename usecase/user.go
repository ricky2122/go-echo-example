package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/ricky2122/go-echo-example/domain"
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserNotFound      = errors.New("user not found")
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

type GetUserUseCaseInput struct {
	ID int
}

type GetUserUseCaseOutput struct {
	ID       int
	Name     string
	Email    string
	BirthDay string
}

type GetUsersUseCaseOutput struct {
	Users []GetUserUseCaseOutput
}

type IUserRepository interface {
	IsExist(ctx context.Context, name string) (bool, error)
	Create(ctx context.Context, newUser domain.User) (*domain.User, error)
	GetUserByID(ctx context.Context, id domain.UserID) (*domain.User, error)
	GetUsers(ctx context.Context) ([]domain.User, error)
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

func (uc *UserUseCase) GetUser(input GetUserUseCaseInput) (*GetUserUseCaseOutput, error) {
	// get user by UserID
	ctx := context.Background()
	userID := domain.UserID(input.ID)
	user, err := uc.ur.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// check if user exists
	if user == nil {
		return nil, ErrUserNotFound
	}

	output := &GetUserUseCaseOutput{
		ID:       user.GetID().Int(),
		Name:     user.GetName(),
		Email:    user.GetEmail(),
		BirthDay: user.GetBirthDay().String(),
	}
	return output, nil
}

func (uc *UserUseCase) GetUsers() (*GetUsersUseCaseOutput, error) {
	ctx := context.Background()
	users, err := uc.ur.GetUsers(ctx)
	if err != nil {
		return nil, err
	}

	outputUsers := make([]GetUserUseCaseOutput, 0, len(users))
	for _, user := range users {
		outputUser := GetUserUseCaseOutput{
			ID:       user.GetID().Int(),
			Name:     user.GetName(),
			Email:    user.GetEmail(),
			BirthDay: user.GetBirthDay().String(),
		}
		outputUsers = append(outputUsers, outputUser)
	}

	output := &GetUsersUseCaseOutput{Users: outputUsers}

	return output, nil
}
