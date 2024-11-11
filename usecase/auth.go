package usecase

import "errors"

var ErrLoginFailed = errors.New("failed login")

type LoginUseCaseInput struct {
	Name     string
	Password string
}

type AuthUseCase struct{}

func NewAuthUseCase() *AuthUseCase {
	return &AuthUseCase{}
}

func (au *AuthUseCase) Login(input LoginUseCaseInput) error {
	// Todo: Implement
	return nil
}
