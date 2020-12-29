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

	ar := repository.App{
		User:   newUser(dbConn),
		Forum:  newForum(dbConn),
		Thread: newThread(dbConn),
		Post:   newPost(dbConn),
	}

	if err = ar.User.Prepare(); err != nil {
		return nil, err
	}

	if err = ar.Forum.Prepare(); err != nil {
		return nil, err
	}

	if err = ar.Thread.Prepare(); err != nil {
		return nil, err
	}

	if err = ar.Post.Prepare(); err != nil {
		return nil, err
	}

	return &ar, nil
}
