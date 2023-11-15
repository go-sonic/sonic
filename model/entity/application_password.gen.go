// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package entity

import (
	"time"
)

const TableNameApplicationPassword = "application_password"

// ApplicationPassword mapped from table <application_password>
type ApplicationPassword struct {
	ID               *int32     `gorm:"column:id;type:INTEGER" json:"id"`
	CreateTime       time.Time  `gorm:"column:create_time;type:datetime;not null" json:"create_time"`
	UpdateTime       *time.Time `gorm:"column:update_time;type:datetime" json:"update_time"`
	Name             string     `gorm:"column:name;type:varchar(32);not null" json:"name"`
	Password         string     `gorm:"column:password;type:varchar(256);not null" json:"password"`
	UserID           int32      `gorm:"column:user_id;type:integer;not null" json:"user_id"`
	LastActivateTime *time.Time `gorm:"column:last_activate_time;type:datetime" json:"last_activate_time"`
	LastActivateIP   string     `gorm:"column:last_activate_ip;type:varchar(128);not null;default:'' not null" json:"last_activate_ip"`
}

// TableName ApplicationPassword's table name
func (*ApplicationPassword) TableName() string {
	return TableNameApplicationPassword
}
