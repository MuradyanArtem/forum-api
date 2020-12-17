package persistence

import (
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

type AppRepository struct {
	user   *UserDB
	forum  *ForumDB
	thread *ThreadDB
	post   *PostDB
}

func New(conf *DBConfig) (*AppRepository, error) {
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

	var ar AppRepository
	ar.user = NewUser(dbConn)
	ar.forum = NewForum(dbConn)
	ar.thread = NewThread(dbConn)
	ar.post = NewPost(dbConn)

	err = ar.user.Prepare()
	if err != nil {
		return nil, err
	}
	err = ar.user.Prepare()
	if err != nil {
		return nil, err
	}
	err = ar.user.Prepare()
	if err != nil {
		return nil, err
	}
	err = ar.user.Prepare()
	if err != nil {
		return nil, err
	}

	return &ar, nil
}
