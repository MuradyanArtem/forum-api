package infrastructure

import "github.com/pkg/errors"

var (
	UserExist      = errors.New("User is already exist")
	UserNotExist   = errors.New("User doesn't exist")
	UserNotUpdated = errors.New("Can't update user")
	ForumExist     = errors.New("Forum is already exist")
	ForumNotExist  = errors.New("Forum doesn't exist")
	PostNotExist   = errors.New("Post doesn't exist")
	ParentNotExist = errors.New("Parent doesn't exist")
	ThreadNotExist = errors.New("Thread doesn't exist")
	ThreadExist    = errors.New("Thread is already exist")
)
