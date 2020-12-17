package utils

import "github.com/pkg/errors"

var (
	UserNotExist   = errors.New("User doesn't exist")
	UserNotUpdated = errors.New("User cannot be updated")
	ForumExist     = errors.New("Forum is already exist")
	ForumNotExist  = errors.New("Forum doesn't exist")
	ParentNotExist = errors.New("Parent doesn't exist")
	ThreadNotExist = errors.New("Thread doesn't exist")
	PostNotExist   = errors.New("Post doesn't exist")
)
