package repository

import (
	"context"
	"webook/webook-server/internal/domain"
	"webook/webook-server/internal/repository/cache"
	"webook/webook-server/internal/repository/dao"
)

var (
	ErrUserDuplicateEmail = dao.ErrUserDuplicateEmail
	ErrUserNotFound       = dao.ErrUserNotFound
)

type UserRepository struct {
	dao   *dao.UserDAO
	cache *cache.UserCache
}

func NewUserRepository(dao *dao.UserDAO, c *cache.UserCache) *UserRepository {
	return &UserRepository{
		dao:   dao,
		cache: c,
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
	u, err := r.cache.Get(ctx, userId)
	if err != nil {
		return u, nil
	}

	ue, err := r.dao.FindbyId(ctx, userId)
	if err != nil {
		return domain.User{}, err
	}
	u = domain.User{
		Id:           ue.Id,
		Email:        ue.Email,
		Password:     ue.Password,
		Nickname:     ue.Nickname,
		Birthday:     ue.Birthday,
		Introduction: ue.Introduction,
		Avatar:       ue.Avatar,
	}

	go func() {
		err = r.cache.Set(ctx, u)
		if err != nil {
			//可能是redis挂了
			//打日志，做监控
		}
	}()

	return u, err
	// 这里怎么办？ err = io.EOF
	// 要不要去数据库加载？
	// 看起来我不应该加载？
	// 看起来我好像也要加载？

	// 选加载 —— 做好兜底，万一 Redis 真的崩了，你要保护住你的数据库
	// 我数据库限流呀！

	// 选不加载 —— 用户体验差一点

	// 缓存里面有数据
	// 缓存里面没有数据
	// 缓存出错了，你也不知道有没有数据

}
