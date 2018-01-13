package service

import (
	"fmt"
	"log"

	"github.com/couatl/forum-db-api/models"
	"github.com/couatl/forum-db-api/restapi/operations"
	"github.com/go-openapi/runtime/middleware"
	_ "github.com/lib/pq"
)

const (
	ERR_NOT_FOUND      = "Can't find!"
	ERR_ALREADY_EXISTS = "Already exists!"
)

type ID struct {
	ID int64 `db:"id"`
}

type ForumPgSQL struct {
	ForumGeneric
}

func NewForumPgSQL(dataSourceName string) ForumHandler {
	return ForumPgSQL{ForumGeneric: NewForumGeneric("postgres", dataSourceName)}
}

//Clear ...
func (dbManager ForumPgSQL) Clear(params operations.ClearParams) middleware.Responder {
	tx := dbManager.db.MustBegin()
	defer tx.Rollback()

	_, err := tx.Exec(`TRUNCATE TABLE forums, threads, users CASCADE`)
	check(err)
	check(tx.Commit())

	return operations.NewClearOK()
}

//ForumCreate ...
func (dbManager ForumPgSQL) ForumCreate(params operations.ForumCreateParams) middleware.Responder {
	tx := dbManager.db.MustBegin()

	id := ID{}
	forum := models.Forum{}
	err := tx.Get(&id, `SELECT id FROM users WHERE nickname = $1`, params.Forum.User)

	if err != nil {
		log.Println(err)
		tx.Rollback()
		return operations.NewForumCreateNotFound().WithPayload(&models.Error{Message: ERR_NOT_FOUND})
	}

	err = tx.Get(&forum, `SELECT posts, slug, threads, title, users.nickname as user
		FROM forums, users WHERE lower(slug) = lower($1) AND users.id = forums.user_id`, params.Forum.Slug)

	if err == nil {
		log.Println(err)
		tx.Rollback()
		return operations.NewForumCreateConflict().WithPayload(&forum)
	}

	tx.Get(&forum, `INSERT INTO forums (slug, user_id, title)
		VALUES ($1, $2, $3) RETURNING slug, title, posts, threads`, params.Forum.Slug, id.ID, params.Forum.Title)
	forum.User = params.Forum.User

	tx.Commit()
	return operations.NewForumCreateCreated().WithPayload(&forum)
}

//ForumGetOne ...
func (dbManager ForumPgSQL) ForumGetOne(params operations.ForumGetOneParams) middleware.Responder {
	tx := dbManager.db.MustBegin()

	forum := models.Forum{}
	err := tx.Get(&forum, `SELECT posts, slug, threads, title, users.nickname as user
		FROM forums, users
		WHERE forums.user_id = users.id AND lower(slug) = lower($1)`, params.Slug)

	if err != nil {
		log.Println(err)
		tx.Rollback()
		return operations.NewForumGetOneNotFound().WithPayload(&models.Error{Message: ERR_NOT_FOUND})
	}

	tx.Commit()
	return operations.NewForumGetOneOK().WithPayload(&forum)
}

//ForumGetThreads TODO
func (dbManager ForumPgSQL) ForumGetThreads(params operations.ForumGetThreadsParams) middleware.Responder {
	tx := dbManager.db.MustBegin()

	id := ID{}
	threads := models.Threads{}
	err := tx.Get(&id, `SELECT id FROM forums WHERE slug = $1`, params.Slug)

	if err != nil {
		log.Println(err)
		tx.Rollback()
		return operations.NewForumGetThreadsNotFound().WithPayload(&models.Error{Message: ERR_NOT_FOUND})
	}

	query := `SELECT threads.id, threads.created, threads.message, threads.slug, threads.votes,
	 threads.title, users.nickname as author, forums.title as forum
	 FROM threads JOIN forums ON (forums.id = threads.forum_id AND forums.id = $1)
	 JOIN users ON (users.id = threads.author_id AND threads.forum_id = $1)
	 `

	desc := params.Desc != nil && *params.Desc
	if params.Since != nil {
		query += ` WHERE threads.created `
		if desc {
			query += `<= ` + params.Since.String()
		} else {
			query += `>= ` + params.Since.String()
		}
	}
	query += ` AND forums.id = $1`
	query += ` ORDER BY threads.created`
	if desc {
		query += ` DESC`
	}
	if params.Limit != nil {
		query += fmt.Sprintf("LIMIT $%d", params.Limit)
	}
	tx.Select(&threads, query, id.ID)

	return operations.NewForumGetOneOK()
}

//ForumGetUsers TODO
func (dbManager ForumPgSQL) ForumGetUsers(params operations.ForumGetUsersParams) middleware.Responder {
	tx := dbManager.db.MustBegin()
	defer tx.Rollback()

	return middleware.NotImplemented("operation .ForumGetUsers has not yet been implemented")
}

func (dbManager ForumPgSQL) PostGetOne(params operations.PostGetOneParams) middleware.Responder {
	return middleware.NotImplemented("operation .PostGetOneParams has not yet been implemented")
}

func (dbManager ForumPgSQL) PostUpdate(params operations.PostUpdateParams) middleware.Responder {
	return middleware.NotImplemented("operation has not yet been implemented")
}

func (dbManager ForumPgSQL) PostsCreate(params operations.PostsCreateParams) middleware.Responder {
	return middleware.NotImplemented("operation .PostsCreateParams has not yet been implemented")
}

func (dbManager ForumPgSQL) Status(params operations.StatusParams) middleware.Responder {
	return middleware.NotImplemented("operation .StatusParams has not yet been implemented")
}

func (dbManager ForumPgSQL) ThreadCreate(params operations.ThreadCreateParams) middleware.Responder {
	return middleware.NotImplemented("operation .ThreadCreateParams has not yet been implemented")
}

func (dbManager ForumPgSQL) ThreadGetOne(params operations.ThreadGetOneParams) middleware.Responder {
	return middleware.NotImplemented("operation .ThreadGetOneParams has not yet been implemented")
}

func (dbManager ForumPgSQL) ThreadGetPosts(params operations.ThreadGetPostsParams) middleware.Responder {
	return middleware.NotImplemented("operation .ThreadUpdateParams has not yet been implemented")
}

func (dbManager ForumPgSQL) ThreadUpdate(params operations.ThreadUpdateParams) middleware.Responder {
	return middleware.NotImplemented("operation .ThreadVoteParams has not yet been implemented")
}

func (dbManager ForumPgSQL) ThreadVote(params operations.ThreadVoteParams) middleware.Responder {
	return middleware.NotImplemented("operation .ForumGetThreadsParams has not yet been implemented")
}

//UserCreate ...
func (dbManager ForumPgSQL) UserCreate(params operations.UserCreateParams) middleware.Responder {
	tx := dbManager.db.MustBegin()

	user := models.User{}
	users := models.Users{}

	tx.Select(&users, "SELECT nickname, fullname, about, email FROM users WHERE users.nickname = $1 OR users.email = $2", params.Nickname, params.Profile.Email)

	log.Println(len(users))
	if len(users) != 0 {
		return operations.NewUserCreateConflict().WithPayload(users)
	}

	tx.Get(&user, "INSERT INTO users (nickname, fullname, about, email) VALUES ($1, $2, $3, $4) RETURNING nickname, fullname, about, email",
		params.Nickname, params.Profile.Fullname, params.Profile.About, params.Profile.Email)

	tx.Commit()
	return operations.NewUserCreateCreated().WithPayload(&user)
}

//UserGetOne ...
func (dbManager ForumPgSQL) UserGetOne(params operations.UserGetOneParams) middleware.Responder {
	tx := dbManager.db.MustBegin()
	defer tx.Rollback()

	users := []models.User{}
	check(tx.Select(&users, "SELECT nickname, fullname, about, email FROM users WHERE users.nickname = $1", params.Nickname))
	check(tx.Commit())

	if len(users) == 0 {
		return operations.NewUserGetOneNotFound().WithPayload(&models.Error{Message: ERR_NOT_FOUND})
	}

	return operations.NewUserGetOneOK().WithPayload(&users[0])
}

//UserUpdate TODO: Change to full query without COALESCE
func (dbManager ForumPgSQL) UserUpdate(params operations.UserUpdateParams) middleware.Responder {
	tx := dbManager.db.MustBegin()

	user := models.User{}
	err := tx.Get(&user, "SELECT nickname FROM users WHERE lower(users.nickname) = lower($1)", params.Nickname)
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return operations.NewUserUpdateNotFound().WithPayload(&models.Error{Message: ERR_NOT_FOUND})
	}

	err = tx.Get(&user, "UPDATE users SET fullname = COALESCE($1, fullname), email = COALESCE($2, email), about = COALESCE($3, about) WHERE nickname = $4 RETURNING about, email, fullname, nickname",
		params.Profile.Fullname, params.Profile.Email, params.Profile.About, params.Nickname)

	if err != nil {
		log.Println(err)
		tx.Rollback()
		return operations.NewUserUpdateConflict().WithPayload(&models.Error{Message: ERR_ALREADY_EXISTS})
	}

	tx.Commit()
	return operations.NewUserUpdateOK().WithPayload(&user)
}
