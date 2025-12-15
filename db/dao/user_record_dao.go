package dao

import (
	"fmt"
	"gorm.io/gorm"
	"time"
	"wxcloudrun-golang/db"
	"wxcloudrun-golang/db/model"
)

const userRecordTable = "user_record"

func ListUserRecords(uid string, pageNum, pageSize int) ([]model.UserRecord, error) {
	cli := db.Get()
	var records []model.UserRecord
	q := cli.Table(userRecordTable)
	if uid != "" {
		q = q.Where("uid = ?", uid)
	}
	offset := 0
	if pageNum > 1 && pageSize > 0 {
		offset = (pageNum - 1) * pageSize
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	if err := q.Order("id DESC").Limit(pageSize).Offset(offset).Find(&records).Error; err != nil {
		return nil, err
	}
	return records, nil
}

func CountUserRecords(uid string) (int64, error) {
	cli := db.Get()
	var count int64
	q := cli.Table(userRecordTable)
	if uid != "" {
		q = q.Where("uid = ?", uid)
	}
	if err := q.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func CreateUserRecordIfNotExists(uid, name, concat string) error {
	cli := db.Get()
	var rec model.UserRecord
	err := cli.Table(userRecordTable).Where("uid = ?", uid).Take(&rec).Error
	if err == nil {
		return fmt.Errorf("uid already exists")
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	hasConcat := cli.Migrator().HasColumn(&model.UserRecord{}, "concat")
	rec = model.UserRecord{UID: uid, Name: name}
	if hasConcat {
		rec.Concat = concat
		return cli.Table(userRecordTable).Create(&rec).Error
	}
	return cli.Table(userRecordTable).Omit("Concat").Create(&rec).Error
}

func UpdateUserRecordByUID(uid string, name *string, concat *string) error {
	cli := db.Get()
	var rec model.UserRecord
	if err := cli.Table(userRecordTable).Where("uid = ?", uid).Take(&rec).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("uid not found")
		}
		return err
	}
	updates := map[string]interface{}{}
	if name != nil {
		updates["name"] = *name
	}
	if concat != nil {
		if cli.Migrator().HasColumn(&model.UserRecord{}, "concat") {
			updates["concat"] = *concat
		}
	}
	if len(updates) == 0 {
		return fmt.Errorf("no fields to update")
	}
	return cli.Table(userRecordTable).Where("uid = ?", uid).Updates(updates).Error
}

func CheckinUserRecord(uid string) error {
	cli := db.Get()
	var rec model.UserRecord
	if err := cli.Table(userRecordTable).Where("uid = ?", uid).Take(&rec).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("uid not found")
		}
		return err
	}
	return cli.Table(userRecordTable).Where("uid = ?", uid).Update("updated_at", time.Now()).Error
}
