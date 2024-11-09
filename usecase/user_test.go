package usecase_test

import (
	"context"
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

func (s *TestStubUserRepository) IsExist(_ context.Context, name string) (bool, error) {
	for _, user := range s.userStore {
		if name == user.GetName() {
			return true, nil
		}
	}
	return false, nil
}

func (s *TestStubUserRepository) Create(_ context.Context, newUser domain.User) (*domain.User, error) {
	newUser.SetID(len(s.userStore) + 1)
	s.userStore = append(s.userStore, newUser)
	return &newUser, nil
}

func (s *TestStubUserRepository) GetUserByID(_ context.Context, userID domain.UserID) (*domain.User, error) {
	for _, user := range s.userStore {
		if userID == user.GetID() {
			return &user, nil
		}
	}
	return nil, nil
}

func (s *TestStubUserRepository) GetUsers(_ context.Context) ([]domain.User, error) {
	return s.userStore, nil
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

func TestGetUserUseCase(t *testing.T) {
	t.Run("Success GetUser", func(t *testing.T) {
		user01 := domain.NewUser(
			"test01",
			"test01",
			"test01@test.com",
			time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
		)
		user01.SetID(1)

		user02 := domain.NewUser(
			"test02",
			"test02",
			"test02@test.com",
			time.Date(2002, 1, 1, 0, 0, 0, 0, time.UTC),
		)
		user02.SetID(2)

		users := []domain.User{user01, user02}
		uuc := usecase.NewUserUseCase(&TestStubUserRepository{userStore: users})

		cases := []struct {
			name  string
			input usecase.GetUserUseCaseInput
			want  *usecase.GetUserUseCaseOutput
		}{
			{
				name:  "first user",
				input: usecase.GetUserUseCaseInput{ID: 1},
				want: &usecase.GetUserUseCaseOutput{
					ID:       1,
					Name:     "test01",
					Email:    "test01@test.com",
					BirthDay: "2001-01-01",
				},
			},
			{
				name:  "second user",
				input: usecase.GetUserUseCaseInput{ID: 2},
				want: &usecase.GetUserUseCaseOutput{
					ID:       2,
					Name:     "test02",
					Email:    "test02@test.com",
					BirthDay: "2002-01-01",
				},
			},
		}

		for _, tt := range cases {
			t.Run(tt.name, func(t *testing.T) {
				got, err := uuc.GetUser(tt.input)
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			})
		}
	})

	t.Run("user not found", func(t *testing.T) {
		uuc := usecase.NewUserUseCase(&TestStubUserRepository{})
		input := usecase.GetUserUseCaseInput{ID: 1}

		_, err := uuc.GetUser(input)
		assert.Equal(t, usecase.ErrUserNotFound, err)
	})
}

func TestGetUsersUseCase(t *testing.T) {
	t.Run("Success GetUsers", func(t *testing.T) {
		cases := []struct {
			name string
			want *usecase.GetUsersUseCaseOutput
		}{
			{
				name: "empty",
				want: nil,
			},
			{
				name: "two users",
				want: &usecase.GetUsersUseCaseOutput{
					Users: []usecase.GetUserUseCaseOutput{
						{
							ID:       1,
							Name:     "test01",
							Email:    "test01@test.com",
							BirthDay: "2001-01-01",
						},
						{
							ID:       2,
							Name:     "test02",
							Email:    "test02@test.com",
							BirthDay: "2002-01-01",
						},
					},
				},
			},
		}

		for _, tt := range cases {
			var store []domain.User
			switch tt.name {
			case "empty":
				store = nil
			case "two users":
				user01 := domain.NewUser(
					"test01",
					"test01",
					"test01@test.com",
					time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
				)
				user01.SetID(1)

				user02 := domain.NewUser(
					"test02",
					"test02",
					"test02@test.com",
					time.Date(2002, 1, 1, 0, 0, 0, 0, time.UTC),
				)
				user02.SetID(2)

				store = []domain.User{user01, user02}
			}
			uuc := usecase.NewUserUseCase(&TestStubUserRepository{userStore: store})
			got, err := uuc.GetUsers()
			if assert.NoError(t, err) {
				assert.Equal(t, tt.want, got)
			}
		}
	})
}
