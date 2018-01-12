package service

import (
	"github.com/couatl/forum-db-api/restapi/operations"
	"github.com/go-openapi/runtime/middleware"
	_ "github.com/lib/pq"
)

type ForumPgSQL struct {
	ForumGeneric
}

func NewForumPgSQL(dataSourceName string) ForumHandler {
	return ForumPgSQL{ForumGeneric: NewForumGeneric("postgres", dataSourceName)}
}

func (self ForumPgSQL) Clear(params operations.ClearParams) middleware.Responder {
	return middleware.NotImplemented("operation .Clear has not yet been implemented")
}

func (self ForumPgSQL) ForumCreate(params operations.ForumCreateParams) middleware.Responder {
	return middleware.NotImplemented("operation .ForumCreate has not yet been implemented")
}

func (self ForumPgSQL) ForumGetOne(params operations.ForumGetOneParams) middleware.Responder {
	return middleware.NotImplemented("operation .ForumGetOne has not yet been implemented")
}

func (self ForumPgSQL) ForumGetThreads(params operations.ForumGetThreadsParams) middleware.Responder {
	return middleware.NotImplemented("operation .ForumGetThreads has not yet been implemented")
}

func (self ForumPgSQL) ForumGetUsers(params operations.ForumGetUsersParams) middleware.Responder {
	return middleware.NotImplemented("operation .ForumGetUsers has not yet been implemented")
}

func (self ForumPgSQL) PostGetOne(params operations.PostGetOneParams) middleware.Responder {
	return middleware.NotImplemented("operation .PostGetOneParams has not yet been implemented")
}

func (self ForumPgSQL) PostUpdate(params operations.PostUpdateParams) middleware.Responder {
	return middleware.NotImplemented("operation has not yet been implemented")
}

func (self ForumPgSQL) PostsCreate(params operations.PostsCreateParams) middleware.Responder {
	return middleware.NotImplemented("operation .PostsCreateParams has not yet been implemented")
}

func (self ForumPgSQL) Status(params operations.StatusParams) middleware.Responder {
	return middleware.NotImplemented("operation .StatusParams has not yet been implemented")
}

func (self ForumPgSQL) ThreadCreate(params operations.ThreadCreateParams) middleware.Responder {
	return middleware.NotImplemented("operation .ThreadCreateParams has not yet been implemented")
}

func (self ForumPgSQL) ThreadGetOne(params operations.ThreadGetOneParams) middleware.Responder {
	return middleware.NotImplemented("operation .ThreadGetOneParams has not yet been implemented")
}

func (self ForumPgSQL) ThreadGetPosts(params operations.ThreadGetPostsParams) middleware.Responder {
	return middleware.NotImplemented("operation .ThreadUpdateParams has not yet been implemented")
}

func (self ForumPgSQL) ThreadUpdate(params operations.ThreadUpdateParams) middleware.Responder {
	return middleware.NotImplemented("operation .ThreadVoteParams has not yet been implemented")
}

func (self ForumPgSQL) ThreadVote(params operations.ThreadVoteParams) middleware.Responder {
	return middleware.NotImplemented("operation .ForumGetThreadsParams has not yet been implemented")
}

func (self ForumPgSQL) UserCreate(params operations.UserCreateParams) middleware.Responder {
	return middleware.NotImplemented("operation .UserCreateParams has not yet been implemented")
}

func (self ForumPgSQL) UserGetOne(params operations.UserGetOneParams) middleware.Responder {
	return middleware.NotImplemented("operation .UserGetOneParams has not yet been implemented")
}

func (self ForumPgSQL) UserUpdate(params operations.UserUpdateParams) middleware.Responder {
	return middleware.NotImplemented("operation .PostGetOneParams has not yet been implemented")
}
