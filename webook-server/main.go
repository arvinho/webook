package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
	"strings"
	"time"
	"webook/webook-server/internal/repository"
	"webook/webook-server/internal/repository/dao"
	"webook/webook-server/internal/service"
	"webook/webook-server/internal/web"
	"webook/webook-server/internal/web/middleware"
	"webook/webook-server/pkg/ginx/middlewares/ratelimit"
)

func main() {
	//db := initDB()
	//server := initWebServer()
	//
	//u := initUser(db)
	//u.RegisterRouters(server)

	server := gin.Default()
	server.GET("/hello", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "你好，你来了")
	})

	server.Run(":8080")
}

func initWebServer() *gin.Engine {
	server := gin.Default()
	server.Use(func(ctx *gin.Context) {
		println("这是第一个middleware")
	})

	server.Use(func(cxt *gin.Context) {
		println("这是第二个middleware")
	})

	//redis 实现访问限流
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	server.Use(ratelimit.NewBuilder(redisClient, time.Second, 100).Build())

	server.Use(cors.New(cors.Config{
		AllowHeaders:  []string{"Content-Type", "Authorization"},
		ExposeHeaders: []string{"x-jwt-token"},
		//是否允许带cookie等身份验证的类
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				return true
			}
			return strings.Contains(origin, "yourcompany.com")
		},
		MaxAge: 12 * time.Hour,
	}))

	//步骤1
	//1、session基于cookie的实现
	//store := cookie.NewStore([]byte("secret"))
	//2、基于内存的session实现
	store := memstore.NewStore([]byte("jH2xG8kO9bW7mV9nO1zD8aD3bS4cV4hQ"), []byte("zH5eF6lF7rS1fQ3kJ8sL4tK3wM9cO9mQ"))
	//3、基于redis的session实现
	//store, err := redis.NewStore(16, "tcp", "localhost:6379", "", []byte("jH2xG8kO9bW7mV9nO1zD8aD3bS4cV4hQ"), []byte("zH5eF6lF7rS1fQ3kJ8sL4tK3wM9cO9mQ"))
	//if err != nil {
	//	panic(err)
	//}
	server.Use(sessions.Sessions("mysession", store))
	//步骤3
	server.Use(middleware.NewLoginMiddlewareBuilder().
		IgnorePaths("/users/signup").
		IgnorePaths("/users/login").Build())
	return server
}

func initUser(db *gorm.DB) *web.UserHandler {
	ud := dao.NewUserDAO(db)
	repo := repository.NewUserRepository(ud)
	svc := service.NewUserService(repo)
	u := web.NewUserHandler(svc)
	return u
}

func initDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:root@tcp(localhost:13316)/webook"))
	if err != nil {
		//只会在初始化过程中panic
		//panic相当于整个 goroutine 结束
		//一旦初始化过程出错，就不需要启动应用
		panic(err)
	}

	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}
	return db
}
