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
	AuthorID      uint32    `gorm:"not null" json:"author_id"`
	Browser       string    `gorm:"size:255;" json:"browser"`
	ExecutionTime time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"exectiontime"`
}

func (u *UserHistory) Prepare() {
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
	if u.AuthorID < 1 {
		return errors.New("Required Author")
	}
	return nil
}

func (u *UserHistory) CreateUserHistory(db *gorm.DB) (*UserHistory, error) {
	var err error
	err = db.Debug().Model(&UserHistory{}).Create(&u).Error
	if err != nil {
		return &UserHistory{}, err
	}
	if u.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", u.AuthorID).Take(&u.Author).Error
		if err != nil {
			return &UserHistory{}, err
		}
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
	if len(histories) > 0 {
		for i, _ := range histories {
			err := db.Debug().Model(&User{}).Where("id = ?", histories[i].AuthorID).Take(&histories[i].Author).Error
			if err != nil {
				return &[]UserHistory{}, err
			}
		}
	}
	return &histories, nil
}

func (u *UserHistory) FindUserHistoryById(db *gorm.DB, uhid uint64) (*UserHistory, error) {
	var err error
	err = db.Debug().Model(&UserHistory{}).Where("id = ?", uhid).Take(&u).Error
	if err != nil {
		return &UserHistory{}, err
	}
	if u.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", u.AuthorID).Take(&u.Author).Error
		if err != nil {
			return &UserHistory{}, err
		}
	}
	return u, nil
}
