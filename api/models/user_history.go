package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type UserHistory struct {
	ID          uint64    `gorm:"primary_key;auto_increment" json:"id"`
	Method      string    `gorm:"size:7;not null;unique" json:"title"`
	Query       string    `gorm:"size:255;not null;" json:"content"`
	QueryStatus string    `gorm:"size:255;not null;" json:"content"`
	Author      User      `json:"author"`
	Browser     string    `gorm:"size:255;" json:"content"`
	Time        time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
}

func (u *UserHistory) prepare() {
	u.ID = 0
	u.Method = html.EscapeString(strings.TrimSpace(u.Method))
	u.Query = html.EscapeString(strings.TrimSpace(u.Query))
	u.QueryStatus = html.EscapeString(strings.TrimSpace(u.QueryStatus))
	u.Author = User{}
	u.Browser = html.EscapeString(strings.TrimSpace(u.Browser))
	u.Time = time.Now()
}

func (u *UserHistory) validate() error {

	if u.Method == "" {
		return errors.New("Required Method")
	}
	if u.Query == "" {
		return errors.New("Required Query")
	}
	if u.QueryStatus == "" {
		return errors.New("Required QueryStatus")
	}
	return nil
}

func (u *UserHistory) createUserHistory(db *gorm.DB) (*UserHistory, error) {
	var err error
	err = db.Debug().Model(&UserHistory{}).Create(&u).Error
	if err != nil {
		return &UserHistory{}, err
	}
	return u, nil
}

func (u *UserHistory) findAllUserHistory(db *gorm.DB) (*[]UserHistory, error) {
	var err error
	histories := []UserHistory{}
	err = db.Debug().Model(&UserHistory{}).Limit(100).Find(&histories).Error
	if err != nil {
		return &[]UserHistory{}, err
	}
	return &histories, nil
}

func (u *UserHistory) findUserHistoryById(db *gorm.DB, uhid uint64) (*UserHistory, error) {
	var err error
	err = db.Debug().Model(&UserHistory{}).Where("id = ?", uhid).Take(&u).Error
	if err != nil {
		return &UserHistory{}, err
	}
	return u, nil
}

func (u *UserHistory) findUserHistoryByUserId(db *gorm.DB, uhid uint64) (*UserHistory, error) {

}

func (u *UserHistory) updateUserHistory(db *gorm.DB) (*UserHistory, error) {

	var err error

	err = db.Debug().Model(&UserHistory{}).Where("id = ?", p.ID).Updates(UserHistory{Title: p.Title, Content: p.Content, UpdatedAt: time.Now()}).Error
	if err != nil {
		return &UserHistory{}, err
	}
	return p, nil
}

func (u *UserHistory) deleteUserHistory(db *gorm.DB, uhid uint64) (int64, error) {

	db = db.Debug().Model(&UserHistory{}).Where("id = ?", uhid).Take(&UserHistory{}).Delete(&UserHistory{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("User History not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
