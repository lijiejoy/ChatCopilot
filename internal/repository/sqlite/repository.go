package sqlite

import (
	"context"
	"fmt"

	"github.com/lw396/WeComCopilot/internal/repository"
	"github.com/lw396/WeComCopilot/pkg/sqlcipher"
	"gorm.io/gorm"
)

type SQLiteClient interface {
	// DB
	OpenDB(ctx context.Context, dbName string) (tx *gorm.DB, err error)
	// Message
	FindMessage(ctx context.Context, tx *gorm.DB, dbName string) (sequence *repository.SQLiteSequence, err error)
	BindMessage(ctx context.Context, tx *gorm.DB, dbName, msgName string) (err error)
	UnbindMessage(ctx context.Context, dbName, msgName string) (err error)
	// Group

}

const (
	MessageDB = "Message/msg_%d.db"  // 消息库
	GroupDB   = "Group/group_new.db" // 群联系
)

type SQLite struct {
	key  string
	path string
	db   map[string]*DB
}

type DB struct {
	tx      *gorm.DB
	msgName []string
}

func NewSQLiteClient(key, path string) *SQLite {
	return &SQLite{
		key:  key,
		path: path,
	}
}

func (s *SQLite) OpenDB(ctx context.Context, dbName string) (tx *gorm.DB, err error) {
	if s.db[dbName] != nil {
		return
	}

	dsn := fmt.Sprintf("%s/%s?_pragma_key=x'%s'", s.path, dbName, s.key)
	tx, err = gorm.Open(sqlcipher.Open(dsn), &gorm.Config{})
	if err != nil {
		return
	}

	return
}
