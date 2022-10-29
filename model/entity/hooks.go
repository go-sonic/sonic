package entity

import (
	"time"

	"gorm.io/gorm"
)

// -----------------------Attachment----------------

func (m *Attachment) BeforeCreate(tx *gorm.DB) (err error) {
	m.CreateTime = time.Now()
	return nil
}

func (m *Attachment) BeforeUpdate(tx *gorm.DB) (err error) {
	tx.Statement.SetColumn("update_time", time.Now())
	return nil
}

// ---------------------- Category ----------------

func (m *Category) BeforeCreate(tx *gorm.DB) (err error) {
	m.CreateTime = time.Now()
	return nil
}

func (m *Category) BeforeUpdate(tx *gorm.DB) (err error) {
	tx.Statement.SetColumn("update_time", time.Now())
	return nil
}

// ------------------ Comment -----------

func (m *Comment) BeforeCreate(tx *gorm.DB) (err error) {
	m.CreateTime = time.Now()
	return nil
}

func (m *Comment) BeforeUpdate(tx *gorm.DB) (err error) {
	tx.Statement.SetColumn("update_time", time.Now())
	return nil
}

// ---------------------- CommentBlack -------------------

func (m *CommentBlack) BeforeCreate(tx *gorm.DB) (err error) {
	m.CreateTime = time.Now()
	return nil
}

func (m *CommentBlack) BeforeUpdate(tx *gorm.DB) (err error) {
	tx.Statement.SetColumn("update_time", time.Now())
	return nil
}

// ---------------------- Journal -----------

func (m *Journal) BeforeCreate(tx *gorm.DB) (err error) {
	m.CreateTime = time.Now()
	return nil
}

func (m *Journal) BeforeUpdate(tx *gorm.DB) (err error) {
	tx.Statement.SetColumn("update_time", time.Now())
	return nil
}

// --------------------- Link -------------------------

func (m *Link) BeforeCreate(tx *gorm.DB) (err error) {
	m.CreateTime = time.Now()
	return nil
}

func (m *Link) BeforeUpdate(tx *gorm.DB) (err error) {
	tx.Statement.SetColumn("update_time", time.Now())
	return nil
}

// ------------------- Log ---------------------------

func (m *Log) BeforeCreate(tx *gorm.DB) (err error) {
	m.CreateTime = time.Now()
	return nil
}

func (m *Log) BeforeUpdate(tx *gorm.DB) (err error) {
	tx.Statement.SetColumn("update_time", time.Now())
	return nil
}

// -------------------- Menu -----------------

func (m *Menu) BeforeCreate(tx *gorm.DB) (err error) {
	m.CreateTime = time.Now()
	return nil
}

func (m *Menu) BeforeUpdate(tx *gorm.DB) (err error) {
	tx.Statement.SetColumn("update_time", time.Now())
	return nil
}

// --------------------------- Option -----------------------

func (m *Option) BeforeCreate(tx *gorm.DB) (err error) {
	m.CreateTime = time.Now()
	return nil
}

func (m *Option) BeforeUpdate(tx *gorm.DB) (err error) {
	tx.Statement.SetColumn("update_time", time.Now())
	return nil
}

// -------------------------- Photo ---------------------

func (m *Photo) BeforeCreate(tx *gorm.DB) (err error) {
	m.CreateTime = time.Now()
	return nil
}

func (m *Photo) BeforeUpdate(tx *gorm.DB) (err error) {
	tx.Statement.SetColumn("update_time", time.Now())
	return nil
}

// ----------------------- Post -------------------------

func (m *Post) BeforeCreate(tx *gorm.DB) (err error) {
	if m.CreateTime == (time.Time{}) {
		m.CreateTime = time.Now()
	}
	m.CreateTime = time.Now()
	return nil
}

func (m *Post) BeforeUpdate(tx *gorm.DB) (err error) {
	tx.Statement.SetColumn("update_time", time.Now())
	return nil
}

// ------------------------- PostCategory ----------------

func (m *PostCategory) BeforeCreate(tx *gorm.DB) (err error) {
	m.CreateTime = time.Now()
	return nil
}

func (m *PostCategory) BeforeUpdate(tx *gorm.DB) (err error) {
	tx.Statement.SetColumn("update_time", time.Now())
	return nil
}

// -------------------- PostTag ----------------------

func (m *PostTag) BeforeCreate(tx *gorm.DB) (err error) {
	m.CreateTime = time.Now()
	return nil
}

func (m *PostTag) BeforeUpdate(tx *gorm.DB) (err error) {
	tx.Statement.SetColumn("update_time", time.Now())
	return nil
}

// ------------------------- Tag ------------------------

func (m *Tag) BeforeCreate(tx *gorm.DB) (err error) {
	m.CreateTime = time.Now()
	return nil
}

func (m *Tag) BeforeUpdate(tx *gorm.DB) (err error) {
	tx.Statement.SetColumn("update_time", time.Now())
	return nil
}

// ------------------------- ThemeSetting --------------------------

func (m *ThemeSetting) BeforeCreate(tx *gorm.DB) (err error) {
	m.CreateTime = time.Now()
	return nil
}

func (m *ThemeSetting) BeforeUpdate(tx *gorm.DB) (err error) {
	tx.Statement.SetColumn("update_time", time.Now())
	return nil
}

// ----------------------- User ---------------------

func (m *User) BeforeCreate(tx *gorm.DB) (err error) {
	m.CreateTime = time.Now()
	return nil
}

func (m *User) BeforeUpdate(tx *gorm.DB) (err error) {
	tx.Statement.SetColumn("update_time", time.Now())
	return nil
}

// ----------------------- Meta ---------------------

func (m *Meta) BeforeCreate(tx *gorm.DB) (err error) {
	m.CreateTime = time.Now()
	return nil
}

func (m *Meta) BeforeUpdate(tx *gorm.DB) (err error) {
	tx.Statement.SetColumn("update_time", time.Now())
	return nil
}
