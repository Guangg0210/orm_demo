package testdata

import (
	"database/sql"
	"sync"
)

type User struct {
	Name     string
	Age      *int
	NickName *sql.NullString
	Picture  []byte
}

type UserDetail struct {
	Address string
}

type SyncMap struct {
	sync.Mutex
	Map map[string]*User
}
