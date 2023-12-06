// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package entity

import (
	"time"

	"github.com/go-sonic/sonic/consts"
)

const TableNameAttachment = "attachment"

// Attachment mapped from table <attachment>
type Attachment struct {
	ID         *int32                `gorm:"column:id;type:integer;primaryKey" json:"id"`
	CreateTime time.Time             `gorm:"column:create_time;type:datetime;not null" json:"create_time"`
	UpdateTime *time.Time            `gorm:"column:update_time;type:datetime" json:"update_time"`
	FileKey    string                `gorm:"column:file_key;type:varchar(2047);not null" json:"file_key"`
	Height     int32                 `gorm:"column:height;type:integer;not null" json:"height"`
	MediaType  string                `gorm:"column:media_type;type:varchar(127);not null" json:"media_type"`
	Name       string                `gorm:"column:name;type:varchar(255);not null" json:"name"`
	Path       string                `gorm:"column:path;type:varchar(1023);not null" json:"path"`
	Size       int64                 `gorm:"column:size;type:bigint;not null" json:"size"`
	Suffix     string                `gorm:"column:suffix;type:varchar(50);not null" json:"suffix"`
	ThumbPath  string                `gorm:"column:thumb_path;type:varchar(1023);not null" json:"thumb_path"`
	Type       consts.AttachmentType `gorm:"column:type;type:bigint;not null" json:"type"`
	Width      int32                 `gorm:"column:width;type:integer;not null" json:"width"`
}

// TableName Attachment's table name
func (*Attachment) TableName() string {
	return TableNameAttachment
}
