package models

import (
	"errors"
	"gorm.io/gorm"
	"strings"
	"time"
	"tuxiaocao/pkg/logger"
	"tuxiaocao/pkg/platform/database"
)

type Curd[T any] struct {
	localDB *gorm.DB
}

type BaseDbTime struct {
	CreatedAt time.Time  `gorm:"column:created_at;autoCreateTime;" json:"created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at;autoUpdateTime;" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:delete_at;default:null" json:"deleted_at"`
}

func (c *Curd[T]) Create(model *T) error {
	if c.localDB == nil {
		c.localDB = database.DB
	}
	return database.DB.Create(&model).Error
}

func (c *Curd[T]) Delete(model *T) error {
	if c.localDB == nil {
		c.localDB = database.DB
	}
	db := c.localDB.Delete(&model)
	if db.RowsAffected == 0 {
		return errors.New("no affected rows")
	}
	return db.Error
}

func (c *Curd[T]) Updates(model *T) error {
	if c.localDB == nil {
		var t T
		c.localDB = database.DB.Model(&t)
	}
	db := c.localDB.Updates(&model)
	return db.Error
}

func (c *Curd[T]) Table(name string, args ...interface{}) *Curd[T] {
	c.localDB = database.DB.Table(name, args...)
	return c
}

func (c *Curd[T]) Where(query interface{}, args ...interface{}) *Curd[T] {
	if c.localDB == nil {
		var t T
		c.localDB = database.DB.Model(&t)
	}
	c.localDB = c.localDB.Where(query, args...)
	return c
}

func (c *Curd[T]) Count() (count int64) {
	if c.localDB == nil {
		var t T
		c.localDB = database.DB.Model(&t)
	}
	c.localDB = c.localDB.Count(&count)
	return
}

func (c *Curd[T]) OR(query interface{}, args ...interface{}) *Curd[T] {
	if c.localDB == nil {
		var t T
		c.localDB = database.DB.Model(&t)
	}
	c.localDB = c.localDB.Or(query, args...)
	return c
}

func (c *Curd[T]) Preload(query string, args ...interface{}) *Curd[T] {
	if c.localDB == nil {
		var t T
		c.localDB = database.DB.Model(t)
	}
	c.localDB = c.localDB.Preload(query, args...)
	return c
}

func (c *Curd[T]) Take() (t T, err error) {
	if c.localDB == nil {
		c.localDB = database.DB.Model(t)
	}
	if err = c.localDB.Take(&t).Error; err != nil {
		return t, err
	}
	return t, nil
}

func (c *Curd[T]) Select(query interface{}, args ...interface{}) *Curd[T] {
	if c.localDB == nil {
		c.localDB = database.DB
	}
	c.localDB = c.localDB.Select(query, args...)
	return c
}

func (c *Curd[T]) Omit(columns ...string) *Curd[T] {
	if c.localDB == nil {
		c.localDB = database.DB
	}
	c.localDB = c.localDB.Omit(columns...)
	return c
}

func (c *Curd[T]) Scan(dest interface{}) error {
	if c.localDB == nil {
		c.localDB = database.DB
	}
	if err := c.localDB.Scan(dest).Error; err != nil {
		return err
	}
	return nil
}

// List 查询
func (c *Curd[T]) List(op *QueryOption) (data []T, total int64, err error) {
	if c.localDB == nil {
		var t T
		c.localDB = database.DB.Model(&t)
	}
	if op != nil {
		// 大于0表示需要分页
		if op.limit > 0 {
			c.localDB = c.localDB.Limit(op.limit)
		}
		if op.offset > 0 && op.limit > 0 {
			c.localDB = c.localDB.Offset((op.offset - 1) * op.limit)
		}
		// 排序
		if op.order != "" {
			c.localDB = c.localDB.Order(op.order)
		}
	}
	if err = c.localDB.Find(&data).Limit(-1).Offset(-1).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	return data, total, err
}

type QueryOption struct {
	offset int    // 页码
	limit  int    // 每页数量
	order  string // 排序字段
}

// NewOP 查询条件
func NewOP() *QueryOption {
	return &QueryOption{}
}

func (op *QueryOption) SetOffset(offset int) *QueryOption {
	op.offset = offset
	return op
}

func (op *QueryOption) SetLimit(limit int) *QueryOption {
	op.limit = limit
	return op
}

// SetOrder eg: "id desc"
func (op *QueryOption) SetOrder(order string) *QueryOption {
	op.order = order
	return op
}

type Query struct {
	Name     string       `json:"name"`
	PageNo   int          `json:"page_no"`
	PageSize int          `json:"page_size"`
	Sorting  []*SortParam `json:"sorting"`
}

type SortParam struct {
	SortBy     string `json:"sort_by"`
	Descending bool   `json:"descending"`
}

func (c *Curd[T]) SetSortParams(sortParams []*SortParam) (order string) {
	var t T
	if c.localDB == nil {
		c.localDB = database.DB.Model(&t)
	}
	if sortParams == nil || len(sortParams) == 0 {
		sortParams = append(sortParams, &SortParam{SortBy: "id", Descending: true})
	}
	for _, sortParam := range sortParams {
		if sortParam == nil {
			logger.Log.Errorf("sort params is nil: %v", sortParam)
			return ""
		}
		sortBy := sortParam.SortBy
		if strings.HasSuffix(sortBy, "_info") {
			sortBy = strings.TrimSuffix(sortBy, "_info")
		}
		// 判断字段是否存在
		if !c.localDB.Migrator().HasColumn(&t, sortBy) {
			continue
		}
		descending := sortParam.Descending
		direction := "asc"
		if descending {
			direction = "desc"
		}
		if order == "" {
			order = sortBy + " " + direction
		} else {
			order = order + ", " + sortBy + " " + direction
		}
	}
	return order
}
