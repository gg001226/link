package account

import (
	"math/rand"
)

type Account struct {
	id			uint64
	name		string
	password	string
}

type Group struct {
	members	[]uint64
	
}

const (
	SYSTEM_ACCOUNT = iota
	UNREGISTERED_ACCOUNT
)

func (a *Account) GetID() uint64 {return a.id}
func (a *Account) GetName() string {return a.name}
func (a *Account) SetName(name string) {a.name = name}
func (a *Account) SetPassword(psw string) {a.password = psw}

func (a *Account) run() error {
	//Account的事件响应函数
	return nil
}

//创建一个账号对象
func NewAccount(name string, password string) *Account {
	return &Account{id:generateID(), name:name, password:password, }
}

func LoginAccount(id uint64, psw string) *Account {
	return &Account{id:id, password:psw}
}


//给新账号创建一个id
func generateID() uint64 {
	return rand.Uint64()
}