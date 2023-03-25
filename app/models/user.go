package models

import "strconv"

type User struct {
	ID       int64  `json:"id" gorm:"primaryKey"`
	Name     string `json:"name" gorm:"not null;comment:用戶名稱"`
	Mobile   string `json:"mobile" gorm:"not null;index;comment:手機號碼"`
	Password string `json:"-" gorm:"not null;default:'';comment:密碼"`
	Timestamps
	SoftDeletes
}

func (user User) GetId() int64 {
	return user.ID
}

func (user User) GetUid() string {
	return strconv.Itoa(int(user.ID))
}

func (user User) GetName() string {
	return user.Name
}
