package repository

import (
	"context"
	"webook/webook-server/internal/domain"
	"webook/webook-server/internal/repository/dao"
)

var (
	ErrUserDuplicateEmail = dao.ErrUserDuplicateEmail
	ErrUserNotFound       = dao.ErrUserNotFound
)

type UserRepository struct {
	dao *dao.UserDAO
}

func NewUserRepository(dao *dao.UserDAO) *UserRepository {
	return &UserRepository{
		dao: dao,
	}
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := r.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return domain.User{
		Id:       u.Id,
		Email:    u.Email,
		Password: u.Password,
	}, nil
}

func (r *UserRepository) Create(ctx context.Context, u domain.User) error {
	return r.dao.Insert(ctx, dao.User{
		Email:    u.Email,
		Password: u.Password,
	})
}

func (r *UserRepository) Update(ctx context.Context, u domain.User) error {
	return r.dao.Update(ctx, dao.User{
		Id:           u.Id,
		Nickname:     u.Nickname,
		Birthday:     u.Birthday,
		Introduction: u.Introduction,
		Avatar:       u.Avatar,
	})
}

func (r *UserRepository) FindById(ctx context.Context, userId int64) (domain.User, error) {
	//先从 cache 里面找
	//再从到里面找
	//找到回写cache
	u, err := r.dao.FindbyId(ctx, userId)
	if err != nil {
		return domain.User{}, err
	}
	return domain.User{
		Id:           u.Id,
		Email:        u.Email,
		Password:     u.Password,
		Nickname:     u.Nickname,
		Birthday:     u.Birthday,
		Introduction: u.Introduction,
		Avatar:       u.Avatar,
	}, nil
}
