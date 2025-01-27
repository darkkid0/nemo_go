package db

import (
	"fmt"
	"gorm.io/gorm"
	"time"
)

type KeyWord struct {
	Id             int       `gorm:"primaryKey"`
	OrgId          int       `gorm:"column:org_id"`
	KeyWord        string    `gorm:"column:key_word"`
	SearchTime     string    `gorm:"column:search_time"`
	ExcludeWords   string    `gorm:"column:exclude_words"`
	CheckMod       string    `gorm:"column:check_mod"`
	IsDelete       bool      `gorm:"column:is_delete"`
	Count          int       `gorm:"column:count"`
	WorkspaceId    int       `gorm:"column:workspace_id"`
	CreateDatetime time.Time `gorm:"column:create_datetime"`
	UpdateDatetime time.Time `gorm:"column:update_datetime"`
}

func (KeyWord) TableName() string {
	return "key_word"
}

// Add 插入一条新的记录，返回主键ID及成功标志
func (t *KeyWord) Add() (success bool) {
	t.CreateDatetime = time.Now()
	t.UpdateDatetime = time.Now()
	t.IsDelete = false
	db := GetDB()
	defer CloseDB(db)
	if result := db.Create(t); result.RowsAffected > 0 {
		return true
	} else {
		return false
	}
}

// Get 根据ID查询记录
func (t *KeyWord) Get() (success bool) {
	db := GetDB()
	defer CloseDB(db)

	if result := db.First(t, t.Id); result.RowsAffected > 0 {
		return true
	} else {
		return false
	}
}

//// GetByOrgId OrgId（不是数据库ID）精确查询一条记录
//func (t *KeyWord) GetByOrgId() (success bool) {
//	db := GetDB()
//	defer CloseDB(db)
//	if result := db.Where("org_id", t.OrgId).First(t); result.RowsAffected > 0 {
//		return true
//	} else {
//		return false
//	}
//}

// Update 更新指定ID的一条记录，列名和内容位于map中
func (t *KeyWord) Update(updateMap map[string]interface{}) (success bool) {
	updateMap["update_datetime"] = time.Now()

	db := GetDB()
	defer CloseDB(db)
	if result := db.Model(t).Updates(updateMap); result.RowsAffected > 0 {
		return true
	} else {
		return false
	}
}

// Delete 删除指定主键ID的一条记录
func (t *KeyWord) Delete() (success bool) {
	db := GetDB()
	defer CloseDB(db)
	if result := db.Delete(t, t.Id); result.RowsAffected > 0 {
		return true
	} else {
		return false
	}
}

// Count 统计指定查询条件的记录数量
func (t *KeyWord) GetCount(searchMap map[string]interface{}) (count int) {
	db := t.makeWhere(searchMap).Model(t)
	defer CloseDB(db)
	var result int64
	db.Count(&result)
	return int(result)
}

// makeWhere 根据查询条件的不同的字段，组合生成count和search的查询条件
func (t *KeyWord) makeWhere(searchMap map[string]interface{}) *gorm.DB {
	db := GetDB()
	for column, value := range searchMap {
		switch column {
		case "key_word":
			db = db.Where("key_word like ?", fmt.Sprintf("%%%s%%", value))
		case "search_time":
			db = db.Where("search_time like ?", fmt.Sprintf("%%%s%%", value))
		case "exclude_words":
			db = db.Where("exclude_words like ?", fmt.Sprintf("%%%s%%", value))
		case "org_id":
			db = db.Where("org_id", value)
		case "check_mod":
			db = db.Where("check_mod like ?", fmt.Sprintf("%%%s%%", value))
		case "date_delta":
			daysToHour := 24 * value.(int)
			dayDelta, err := time.ParseDuration(fmt.Sprintf("-%dh", daysToHour))
			if err == nil {
				db = db.Where("update_datetime between ? and ?", time.Now().Add(dayDelta), time.Now())
			}
		case "workspace_id":
			db = db.Where("workspace_id", value)
		default:

			db = db.Where(column, value)
		}
	}
	return db
}

// Gets 根据指定的条件，查询满足要求的记录
func (t *KeyWord) Gets(searchMap map[string]interface{}, page, rowsPerPage int) (results []KeyWord, count int) {
	orderBy := "update_datetime desc"
	searchMap["is_delete"] = 0
	db := t.makeWhere(searchMap).Model(t)
	defer CloseDB(db)
	//统计满足条件的总记录数
	var total int64
	db.Count(&total)
	//获取分页查询结果
	if rowsPerPage > 0 && page > 0 {
		db = db.Offset((page - 1) * rowsPerPage).Limit(rowsPerPage)
	}
	db.Order(orderBy).Find(&results)

	return results, int(total)
}

//// SaveOrUpdate 保存、更新一条记录
//func (t *KeyWord) SaveOrUpdate() (success bool) {
//	oldRecord := &KeyWord{Id: t.Id}
//	if t.Id >0 {
//		updateMap := map[string]interface{}{}
//		if t.KeyWord != "" {
//			updateMap["key_word"] = t.KeyWord
//		}
//		if t.State != "" {
//			updateMap["state"] = t.State
//		}
//		if t.Result != "" {
//			updateMap["result"] = t.Result
//		}
//		if t.ReceivedTime != nil {
//			updateMap["received"] = t.ReceivedTime
//		}
//		if t.RetriedTime != nil {
//			updateMap["retried"] = t.RetriedTime
//		}
//		if t.RevokedTime != nil {
//			updateMap["revoked"] = t.RevokedTime
//		}
//		if t.StartedTime != nil {
//			updateMap["started"] = t.StartedTime
//		}
//		if t.SucceededTime != nil {
//			updateMap["succeeded"] = t.SucceededTime
//		}
//		if t.FailedTime != nil {
//			updateMap["failed"] = t.FailedTime
//		}
//		if t.ProgressMessage != "" {
//			updateMap["progress_message"] = t.ProgressMessage
//		}
//		t.Id = oldRecord.Id
//		return t.Update(updateMap)
//	} else {
//		return t.Add()
//	}
//}
