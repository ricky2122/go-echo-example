package usecase

type SignUpUseCaseInput struct {
	Name     string
	Password string
	Email    string
	BirthDay string
}

type SignUpUseCaseOutput struct {
	ID   int
	Name string
}

type UserUseCase struct{}

func NewUserUseCase() *UserUseCase {
	return &UserUseCase{}
}

func (uc *UserUseCase) SignUp(input SignUpUseCaseInput) (*SignUpUseCaseOutput, error) {
	// Todo: Implement
	return nil, nil
}
