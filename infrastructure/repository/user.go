package repository

import (
	"context"
	"time"

	"github.com/ricky2122/go-echo-example/domain"
	"github.com/uptrace/bun"
)

type UserModel struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	ID       int       `bun:"id,pk,autoincrement"`
	Name     string    `bun:"name,notnull,unique"`
	Password string    `bun:"password,notnull"`
	Email    string    `bun:"email,notnull,unique"`
	BirthDay time.Time `bun:"birth_day,notnull"`
}

type UserRepository struct {
	db *bun.DB
}

func NewUserRepository(db *bun.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (ur *UserRepository) IsExist(ctx context.Context, name string) (bool, error) {
	exists, err := ur.db.NewSelect().
		Model((*UserModel)(nil)).
		Where("name = ?", name).
		Exists(ctx)
	if err != nil {
		return true, err
	}

	if exists {
		return true, nil
	}

	return false, nil
}

func (ur *UserRepository) Create(ctx context.Context, newUser domain.User) (*domain.User, error) {
	newUserModel := convertToUserModel(newUser)
	_, err := ur.db.NewInsert().
		Model(&newUserModel).
		Returning("id").
		Exec(ctx)
	if err != nil {
		return nil, err
	}
	createdUser := convertToUser(newUserModel)

	return &createdUser, nil
}

func (ur *UserRepository) GetUserByID(ctx context.Context, userID domain.UserID) (*domain.User, error) {
	var userModel UserModel
	if err := ur.db.NewSelect().
		Model(&userModel).
		Where("id = ?", userID).
		Scan(ctx); err != nil {
		return nil, err
	}

	// check if user exists
	if userModel.ID == 0 {
		return nil, nil
	}

	user := convertToUser(userModel)

	return &user, nil
}

func convertToUserModel(user domain.User) UserModel {
	return UserModel{
		ID:       user.GetID().Int(),
		Name:     user.GetName(),
		Password: user.GetPassword(),
		Email:    user.GetEmail(),
		BirthDay: user.GetBirthDay().Time(),
	}
}

func convertToUser(userModel UserModel) domain.User {
	user := domain.NewUser(
		userModel.Name,
		userModel.Password,
		userModel.Email,
		userModel.BirthDay,
	)
	user.SetID(userModel.ID)

	return user
}