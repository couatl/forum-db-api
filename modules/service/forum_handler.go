package service

import (
	"github.com/couatl/forum-db-api/restapi/operations"
	"github.com/go-openapi/runtime/middleware"
)

// ForumHandler ...
type ForumHandler interface {
	Clear(params operations.ClearParams) middleware.Responder
	Status(params operations.StatusParams) middleware.Responder

	ForumCreate(params operations.ForumCreateParams) middleware.Responder
	ForumGetOne(params operations.ForumGetOneParams) middleware.Responder
	ForumGetThreads(params operations.ForumGetThreadsParams) middleware.Responder
	ForumGetUsers(params operations.ForumGetUsersParams) middleware.Responder

	PostGetOne(params operations.PostGetOneParams) middleware.Responder
	PostUpdate(params operations.PostUpdateParams) middleware.Responder
	PostsCreate(params operations.PostsCreateParams) middleware.Responder

	ThreadCreate(params operations.ThreadCreateParams) middleware.Responder
	ThreadGetOne(params operations.ThreadGetOneParams) middleware.Responder
	ThreadGetPosts(params operations.ThreadGetPostsParams) middleware.Responder
	ThreadUpdate(params operations.ThreadUpdateParams) middleware.Responder
	ThreadVote(params operations.ThreadVoteParams) middleware.Responder

	UserCreate(params operations.UserCreateParams) middleware.Responder
	UserGetOne(params operations.UserGetOneParams) middleware.Responder
	UserUpdate(params operations.UserUpdateParams) middleware.Responder
}
