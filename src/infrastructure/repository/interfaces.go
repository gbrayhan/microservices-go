package repository

import (
	"gorm.io/gorm"
)

// DatabaseInterface defines the interface for database operations
type DatabaseInterface interface {
	Find(dest interface{}, conds ...interface{}) DatabaseInterface
	Create(value interface{}) DatabaseInterface
	Limit(limit int) DatabaseInterface
	Where(query interface{}, args ...interface{}) DatabaseInterface
	First(dest interface{}, conds ...interface{}) DatabaseInterface
	Model(value interface{}) DatabaseInterface
	Select(query interface{}, args ...interface{}) DatabaseInterface
	Updates(values interface{}) DatabaseInterface
	Delete(value interface{}, conds ...interface{}) DatabaseInterface
	Error() error
	RowsAffected() int64
}

// GormDBAdapter adapts gorm.DB to DatabaseInterface
type GormDBAdapter struct {
	db *gorm.DB
}

// NewGormDBAdapter creates a new GormDBAdapter
func NewGormDBAdapter(db *gorm.DB) DatabaseInterface {
	return &GormDBAdapter{db: db}
}

func (g *GormDBAdapter) Find(dest interface{}, conds ...interface{}) DatabaseInterface {
	return &GormDBAdapter{db: g.db.Find(dest, conds...)}
}

func (g *GormDBAdapter) Create(value interface{}) DatabaseInterface {
	return &GormDBAdapter{db: g.db.Create(value)}
}

func (g *GormDBAdapter) Limit(limit int) DatabaseInterface {
	return &GormDBAdapter{db: g.db.Limit(limit)}
}

func (g *GormDBAdapter) Where(query interface{}, args ...interface{}) DatabaseInterface {
	return &GormDBAdapter{db: g.db.Where(query, args...)}
}

func (g *GormDBAdapter) First(dest interface{}, conds ...interface{}) DatabaseInterface {
	return &GormDBAdapter{db: g.db.First(dest, conds...)}
}

func (g *GormDBAdapter) Model(value interface{}) DatabaseInterface {
	return &GormDBAdapter{db: g.db.Model(value)}
}

func (g *GormDBAdapter) Select(query interface{}, args ...interface{}) DatabaseInterface {
	return &GormDBAdapter{db: g.db.Select(query, args...)}
}

func (g *GormDBAdapter) Updates(values interface{}) DatabaseInterface {
	return &GormDBAdapter{db: g.db.Updates(values)}
}

func (g *GormDBAdapter) Delete(value interface{}, conds ...interface{}) DatabaseInterface {
	return &GormDBAdapter{db: g.db.Delete(value, conds...)}
}

func (g *GormDBAdapter) Error() error {
	return g.db.Error
}

func (g *GormDBAdapter) RowsAffected() int64 {
	return g.db.RowsAffected
}
