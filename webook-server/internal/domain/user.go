package domain

import "time"

// User领域对象，是DDD中的entity
// BO(busuness object)
type User struct {
	Id           int64
	Email        string
	Password     string
	Nickname     string
	Birthday     string
	Introduction string
	Avatar       string
	Ctime        time.Time
}

//type Address struct {
//}
