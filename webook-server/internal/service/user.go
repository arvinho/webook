package service

import (
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"webook/webook-server/internal/domain"
	"webook/webook-server/internal/repository"
)

var ErrUserDuplicateEmail = repository.ErrUserDuplicateEmail
var ErrInvalidUserOrPassword = errors.New("账号、邮箱或密码不对")

type UserService struct {
	rep *repository.UserRepository
}

func NewUserService(rep *repository.UserRepository) *UserService {
	return &UserService{
		rep: rep,
	}
}

func (svc *UserService) Login(ctx context.Context, email string, password string) (domain.User, error) {
	//先找用户
	u, err := svc.rep.FindByEmail(ctx, email)
	if err == repository.ErrUserNotFound {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	if err != nil {
		return domain.User{}, err
	}

	//比较密码
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		//DEBUG
		return domain.User{}, ErrInvalidUserOrPassword
	}
	return u, nil
}

func (svc *UserService) SignUp(ctx context.Context, u domain.User) error {
	//密码加密
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	//存储
	return svc.rep.Create(ctx, u)
}

func (svc *UserService) Edit(ctx context.Context, u domain.User) error {
	_, err := svc.rep.FindById(ctx, u.Id)
	if err != nil {
		return err
	}
	return svc.rep.Update(ctx, u)
}

func (svc *UserService) Profile(ctx context.Context, userId int64) (domain.User, error) {
	u, err := svc.rep.FindById(ctx, userId)
	if err != nil {
		return domain.User{}, err
	}
	return u, nil
}
