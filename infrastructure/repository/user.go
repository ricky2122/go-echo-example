package repository

import "github.com/ricky2122/go-echo-example/domain"

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (ur *UserRepository) IsExist(name string) (bool, error) {
	// Todo: Implement
	return false, nil
}

func (ur *UserRepository) Create(newUser domain.User) (*domain.User, error) {
	// Todo: Implement
	return nil, nil
}
