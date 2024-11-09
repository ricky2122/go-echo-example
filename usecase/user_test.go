package usecase_test

import (
	"errors"
	"testing"
	"time"

	"github.com/ricky2122/go-echo-example/domain"
	"github.com/ricky2122/go-echo-example/usecase"
	"github.com/stretchr/testify/assert"
)

type TestStubUserRepository struct {
	userStore []domain.User
}

func (s *TestStubUserRepository) IsExist(name string) (bool, error) {
	for _, user := range s.userStore {
		if name == user.GetName() {
			return true, nil
		}
	}
	return false, nil
}

func (s *TestStubUserRepository) Create(newUser domain.User) (*domain.User, error) {
	newUser.SetID(len(s.userStore) + 1)
	s.userStore = append(s.userStore, newUser)
	return &newUser, nil
}

func TestSignUpUseCase(t *testing.T) {
	t.Run("Success SignUp", func(t *testing.T) {
		uuc := usecase.NewUserUseCase(&TestStubUserRepository{})

		cases := []struct {
			name  string
			input usecase.SignUpUseCaseInput
			want  *usecase.SignUpUseCaseOutput
		}{
			{
				name: "first user",
				input: usecase.SignUpUseCaseInput{
					Name:     "test01",
					Password: "test01",
					Email:    "test01@test.com",
					BirthDay: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				},
				want: &usecase.SignUpUseCaseOutput{
					ID:   1,
					Name: "test01",
				},
			},
			{
				name: "second user",
				input: usecase.SignUpUseCaseInput{
					Name:     "test02",
					Password: "test02",
					Email:    "test02@test.com",
					BirthDay: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
				},
				want: &usecase.SignUpUseCaseOutput{
					ID:   2,
					Name: "test02",
				},
			},
		}

		for _, tt := range cases {
			t.Run(tt.name, func(t *testing.T) {
				got, err := uuc.SignUp(tt.input)
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			})
		}
	})

	t.Run("User already exists", func(t *testing.T) {
		uuc := usecase.NewUserUseCase(&TestStubUserRepository{})

		input := usecase.SignUpUseCaseInput{
			Name:     "test01",
			Password: "test01",
			Email:    "test01@test.com",
			BirthDay: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
		}

		_, _ = uuc.SignUp(input)
		_, err := uuc.SignUp(input)
		wantErr := errors.New("user already exists")
		assert.Equal(t, wantErr, err)
	})
}
