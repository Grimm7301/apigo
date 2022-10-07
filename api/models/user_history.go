package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type UserHistory struct {
	ID            uint64    `gorm:"primary_key;auto_increment" json:"id"`
	Method        string    `gorm:"size:7;not null;unique" json:"method"`
	Query         string    `gorm:"size:255;not null;" json:"query"`
	QueryStatus   string    `gorm:"size:255;not null;" json:"querystatus"`
	Author        User      `json:"author"`
	Browser       string    `gorm:"size:255;" json:"browser"`
	ExecutionTime time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"exectiontime"`
}

func (u *UserHistory) prepare() {
	u.ID = 0
	u.Method = html.EscapeString(strings.TrimSpace(u.Method))
	u.Query = html.EscapeString(strings.TrimSpace(u.Query))
	u.QueryStatus = html.EscapeString(strings.TrimSpace(u.QueryStatus))
	u.Author = User{}
	u.Browser = html.EscapeString(strings.TrimSpace(u.Browser))
	u.ExecutionTime = time.Now()
}

func (u *UserHistory) Validate() error {

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

func (u *UserHistory) CreateUserHistory(db *gorm.DB) (*UserHistory, error) {
	var err error
	err = db.Debug().Model(&UserHistory{}).Create(&u).Error
	if err != nil {
		return &UserHistory{}, err
	}
	return u, nil
}

func (u *UserHistory) FindAllUserHistory(db *gorm.DB) (*[]UserHistory, error) {
	var err error
	histories := []UserHistory{}
	err = db.Debug().Model(&UserHistory{}).Limit(100).Find(&histories).Error
	if err != nil {
		return &[]UserHistory{}, err
	}
	return &histories, nil
}

func (u *UserHistory) FindUserHistoryById(db *gorm.DB, uhid uint64) (*UserHistory, error) {
	var err error
	err = db.Debug().Model(&UserHistory{}).Where("id = ?", uhid).Take(&u).Error
	if err != nil {
		return &UserHistory{}, err
	}
	return u, nil
}

func (u *UserHistory) FindUserHistoryByUserId(db *gorm.DB, uhid uint64) (*UserHistory, error) {
	return u, nil
}

func (u *UserHistory) UpdateUserHistory(db *gorm.DB) (*UserHistory, error) {

	var err error

	err = db.Debug().Model(&UserHistory{}).Where("id = ?", u.ID).Updates(UserHistory{Method: u.Method, Query: u.Query, QueryStatus: u.QueryStatus}).Error
	if err != nil {
		return &UserHistory{}, err
	}
	return u, nil
}

func (u *UserHistory) DeleteUserHistory(db *gorm.DB, uhid uint64) (int64, error) {

	db = db.Debug().Model(&UserHistory{}).Where("id = ?", uhid).Take(&UserHistory{}).Delete(&UserHistory{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("User History not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
