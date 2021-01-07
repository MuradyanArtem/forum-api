package persistence

import (
	"forum-api/internal/domain/repository"
	"time"

	"github.com/jackc/pgx"
)

type DBConfig struct {
	Host                 string
	Port                 uint16
	User                 string
	Database             string
	Password             string
	PreferSimpleProtocol bool
	MaxConnections       int
	AcquireTimeout       time.Duration
}

func New(conf *DBConfig) (*repository.App, error) {
	dbConn, err := pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig: pgx.ConnConfig{
			Host:                 conf.Host,
			Port:                 conf.Port,
			User:                 conf.User,
			Database:             conf.Database,
			Password:             conf.Password,
			PreferSimpleProtocol: conf.PreferSimpleProtocol,
		},
		MaxConnections: conf.MaxConnections,
		AcquireTimeout: conf.AcquireTimeout,
	})
	if err != nil {
		return nil, err
	}

	return &repository.App{
		User:   newUser(dbConn),
		Forum:  newForum(dbConn),
		Thread: newThread(dbConn),
		Post:   newPost(dbConn),
	}, nil
}
