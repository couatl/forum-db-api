package service

import (
	"log"
	"strconv"

	"github.com/couatl/forum-db-api/models"
	"github.com/couatl/forum-db-api/restapi/operations"
	"github.com/go-openapi/runtime/middleware"

	_ "github.com/lib/pq"
)

const (
	ERR_NOT_FOUND      = "Can't find!"
	ERR_ALREADY_EXISTS = "Already exists!"
	ERR                = "An error occured!"
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

//Clear ... OK
func (dbManager ForumPgSQL) Clear(params operations.ClearParams) middleware.Responder {
	tx := dbManager.db.MustBegin()
	defer tx.Rollback()

	_, err := tx.Exec(`TRUNCATE TABLE forums, posts, threads, users, votes CASCADE`)
	check(err)
	check(tx.Commit())

	return operations.NewClearOK()
}

//ForumCreate ... OK OK
func (dbManager ForumPgSQL) ForumCreate(params operations.ForumCreateParams) middleware.Responder {
	tx := dbManager.db.MustBegin()

	user := models.User{}
	forum := models.Forum{}
	err := tx.Get(&user, `SELECT nickname FROM users WHERE lower(nickname) = lower($1)`, params.Forum.User)

	if err != nil {
		log.Println(err)
		tx.Rollback()
		return operations.NewForumCreateNotFound().WithPayload(&models.Error{Message: ERR_NOT_FOUND})
	}

	errAlreadyExists := tx.Get(&forum, `SELECT posts, slug, threads, title, author as user
		FROM forums WHERE lower(slug) = lower($1)`, params.Forum.Slug)

	if errAlreadyExists == nil {
		log.Println(errAlreadyExists)
		tx.Rollback()
		return operations.NewForumCreateConflict().WithPayload(&forum)
	}

	tx.Get(&forum, `INSERT INTO forums (slug, author, title)
		VALUES ($1, $2, $3) RETURNING slug, title, posts, threads, author as user`,
		params.Forum.Slug, user.Nickname, params.Forum.Title)

	tx.Commit()
	return operations.NewForumCreateCreated().WithPayload(&forum)
}

//ForumGetOne ... OK OK
func (dbManager ForumPgSQL) ForumGetOne(params operations.ForumGetOneParams) middleware.Responder {
	tx := dbManager.db.MustBegin()

	forum := models.Forum{}
	err := tx.Get(&forum, `SELECT slug, title, author as user, threads, posts
		FROM forums
		WHERE lower(slug) = lower($1)`, params.Slug)

	if err != nil {
		log.Println(err)
		tx.Rollback()
		return operations.NewForumGetOneNotFound().WithPayload(&models.Error{Message: ERR_NOT_FOUND})
	}

	log.Println(forum)
	tx.Commit()
	return operations.NewForumGetOneOK().WithPayload(&forum)
}

//ForumGetThreads ... OK OK
func (dbManager ForumPgSQL) ForumGetThreads(params operations.ForumGetThreadsParams) middleware.Responder {
	tx := dbManager.db.MustBegin()

	forum := models.Forum{}
	threads := models.Threads{}

	err := tx.Get(&forum, `SELECT slug FROM forums WHERE lower(slug) = lower($1)`, params.Slug)
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return operations.NewForumGetThreadsNotFound().WithPayload(&models.Error{Message: ERR_NOT_FOUND})
	}

	query := `SELECT * FROM threads WHERE threads.forum = $1 `

	desc := params.Desc != nil && *params.Desc
	if params.Since != nil {
		query += ` AND threads.created `
		if desc {
			query += ` <= $2`
		} else {
			query += ` >= $2`
		}
	}
	query += ` ORDER BY threads.created`
	if desc {
		query += ` DESC`
	}
	if params.Limit != nil {
		query += ` LIMIT ` + strconv.FormatInt(int64(*params.Limit), 10)
	}

	if params.Since != nil {
		tx.Select(&threads, query, forum.Slug, *params.Since)
	} else {
		tx.Select(&threads, query, forum.Slug)
	}

	tx.Commit()
	return operations.NewForumGetThreadsOK().WithPayload(threads)
}

//ForumGetUsers ...
func (dbManager ForumPgSQL) ForumGetUsers(params operations.ForumGetUsersParams) middleware.Responder {
	tx := dbManager.db.MustBegin()

	forum := models.Forum{}
	users := models.Users{}

	err := tx.Get(&forum, `SELECT slug FROM forums WHERE lower(slug) = lower($1)`, params.Slug)
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return operations.NewForumGetUsersNotFound().WithPayload(&models.Error{Message: ERR_NOT_FOUND})
	}

	query := `SELECT DISTINCT users.about, users.email, users.fullname, users.nickname
	 FROM users LEFT JOIN posts ON (users.nickname = posts.author AND posts.forum = $1)
	 LEFT JOIN threads ON (users.nickname = threads.author AND threads.forum = $1)
	 WHERE (posts.forum = $1 OR threads.forum = $1) `

	desc := params.Desc != nil && *params.Desc
	if params.Since != nil {
		query += ` AND lower(users.nickname) `

		if desc {
			query += `< lower($2)`
		} else {
			query += `> lower($2)`
		}
	}
	query += ` ORDER BY lower(users.nickname)`
	if desc {
		query += ` DESC`
	}
	if params.Limit != nil {
		query += ` LIMIT ` + strconv.FormatInt(int64(*params.Limit), 10)
	}

	if params.Since != nil {
		tx.Select(&users, query, forum.Slug, *params.Since)
	} else {
		tx.Select(&users, query, forum.Slug)
	}

	tx.Commit()
	return operations.NewForumGetUsersOK().WithPayload(users)
}

// PostGetOne TODO Затестить
func (dbManager ForumPgSQL) PostGetOne(params operations.PostGetOneParams) middleware.Responder {
	tx := dbManager.db.MustBegin()

	post := models.Post{}
	postFull := models.PostFull{}

	err := tx.Get(&post, `SELECT * FROM posts WHERE id = $1`, params.ID)

	if err != nil {
		log.Println(err)
		tx.Rollback()
		return operations.NewPostGetOneNotFound().WithPayload(&models.Error{Message: ERR_NOT_FOUND})
	}

	postFull.Post = &post

	for _, item := range params.Related {
		if item == "user" {
			user := models.User{}
			errUnexpected := tx.Get(&user, `SELECT about, email, fullname, nickname
				FROM users WHERE lower(nickname) = lower($1)`, post.Author)

			if errUnexpected != nil {
				log.Println(errUnexpected)
			}
			postFull.Author = &user
			continue
		}
		if item == "forum" {
			forum := models.Forum{}
			errUnexpected2 := tx.Get(&forum, `SELECT posts, threads, slug, title, author as user
				FROM forums WHERE lower(slug) = lower($1)`, post.Forum)

			if errUnexpected2 != nil {
				log.Println(errUnexpected2)
			}
			postFull.Forum = &forum

			continue
		}
		if item == "thread" {
			thread := models.Thread{}
			errUnexpected3 := tx.Get(&thread, `SELECT * FROM threads WHERE id = $1`, post.Thread)

			if errUnexpected3 != nil {
				log.Println(errUnexpected3)
			}
			postFull.Thread = &thread

			continue
		}
	}

	tx.Commit()
	return operations.NewPostGetOneOK().WithPayload(&postFull)
}

// PostUpdate TODO Затестить
func (dbManager ForumPgSQL) PostUpdate(params operations.PostUpdateParams) middleware.Responder {
	tx := dbManager.db.MustBegin()

	post := models.Post{}

	err := tx.Get(&post, "UPDATE posts SET is_edited = true, message = $1 WHERE id = $2 RETURNING *",
		params.Post.Message, params.ID)

	if err != nil {
		log.Println(err)
		tx.Rollback()
		return operations.NewUserUpdateConflict().WithPayload(&models.Error{Message: ERR_NOT_FOUND})
	}

	tx.Commit()
	return operations.NewPostUpdateOK().WithPayload(&post)
}

// PostsCreate OK OK
func (dbManager ForumPgSQL) PostsCreate(params operations.PostsCreateParams) middleware.Responder {
	tx := dbManager.db.MustBegin()

	thread := models.Thread{}
	posts := models.Posts{}

	postID := ID{}
	threadID := ID{}

	slug, id := SlugID(params.SlugOrID)

	err := tx.Get(&threadID, `SELECT id FROM threads WHERE lower(slug) = lower($1) OR id = $2`, slug, id)
	tx.Get(&thread, `SELECT slug, forum FROM threads WHERE id = $1`, threadID.ID)
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return operations.NewPostsCreateNotFound().WithPayload(&models.Error{Message: ERR_NOT_FOUND})
	}

	if len(params.Posts) == 0 {
		return operations.NewPostsCreateCreated().WithPayload(params.Posts)
	}

	checkQuery := "SELECT id FROM posts WHERE thread = $1 AND id = $2"

	insertQuery := "INSERT INTO posts (forum, thread, author, message, parent) VALUES "
	for idx, item := range params.Posts {
		log.Println(item.Parent)

		if item.Parent != 0 {
			errNotFound := tx.Get(&postID, checkQuery, threadID.ID, item.Parent)
			if errNotFound != nil {
				log.Println(errNotFound)
				tx.Rollback()
				return operations.NewPostsCreateConflict().WithPayload(&models.Error{Message: ERR_NOT_FOUND})
			}
		}

		if idx > 0 {
			insertQuery += ","
		}

		insertQuery += "('" + thread.Forum + "', " + strconv.FormatInt(int64(threadID.ID), 10) +
			", '" + item.Author + "', '" + item.Message + "', " + strconv.FormatInt(item.Parent, 10) + ")"
	}
	insertQuery += " RETURNING author, created, forum, id, is_edited as isEdited, message, thread, parent"

	errAlreadyExists := tx.Select(&posts, insertQuery)
	if errAlreadyExists != nil {
		log.Println(errAlreadyExists)
		tx.Rollback()
		return operations.NewPostsCreateConflict().WithPayload(&models.Error{Message: ERR_ALREADY_EXISTS})
	}

	tx.MustExec("UPDATE forums SET posts = posts + $1 WHERE slug = $2", len(params.Posts), thread.Forum)
	tx.Commit()
	return operations.NewPostsCreateCreated().WithPayload(posts)
}

// Status ... OK
func (dbManager ForumPgSQL) Status(params operations.StatusParams) middleware.Responder {
	tx := dbManager.db.MustBegin()
	defer tx.Rollback()

	status := models.Status{}

	err := tx.Get(&status, `SELECT (SELECT COUNT(forums.*) FROM forums) as forum,
	(SELECT COUNT(threads.*) FROM threads) as thread,
	(SELECT COUNT(posts.*) FROM posts) as post,
	(SELECT COUNT(users.*) FROM users) as user`)

	check(err)
	check(tx.Commit())

	return operations.NewStatusOK().WithPayload(&status)
}

// ThreadCreate ... OK OK
func (dbManager ForumPgSQL) ThreadCreate(params operations.ThreadCreateParams) middleware.Responder {
	tx := dbManager.db.MustBegin()

	thread := models.Thread{}
	forum := models.Forum{}
	user := models.User{}

	errNotFound := tx.Get(&forum, `SELECT slug FROM forums WHERE lower(slug) = lower($1)`, params.Slug)
	errNotFound2 := tx.Get(&user, `SELECT nickname FROM users WHERE lower(nickname) = lower($1)`, params.Thread.Author)
	if errNotFound != nil || errNotFound2 != nil {
		log.Println(errNotFound)
		log.Println(errNotFound2)
		tx.Rollback()
		return operations.NewThreadCreateNotFound().WithPayload(&models.Error{Message: ERR_NOT_FOUND})
	}

	if params.Thread.Slug != "" {
		errAlreadyExists := tx.Get(&thread, `SELECT * FROM threads WHERE lower(slug) = lower($1)`, params.Thread.Slug)
		if errAlreadyExists == nil {
			log.Println(errAlreadyExists)
			tx.Rollback()
			return operations.NewThreadCreateConflict().WithPayload(&thread)
		}
	}

	err := tx.Get(&thread, `INSERT INTO threads (forum, author, created, message, title, slug)
	VALUES ($1, $2, COALESCE($3, now()), $4, $5, $6) RETURNING *`,
		forum.Slug, user.Nickname, params.Thread.Created, params.Thread.Message, params.Thread.Title, params.Thread.Slug)
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return operations.NewThreadCreateNotFound().WithPayload(&models.Error{Message: ERR})
	}

	tx.Commit()
	return operations.NewThreadCreateCreated().WithPayload(&thread)
}

// ThreadGetOne ... OK
func (dbManager ForumPgSQL) ThreadGetOne(params operations.ThreadGetOneParams) middleware.Responder {
	tx := dbManager.db.MustBegin()

	thread := models.Thread{}

	slug, id := SlugID(params.SlugOrID)
	err := tx.Get(&thread, `SELECT * FROM threads WHERE lower(slug) = lower($1) OR id = $2`, slug, id)
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return operations.NewThreadGetOneNotFound().WithPayload(&models.Error{Message: ERR_NOT_FOUND})
	}

	tx.Commit()
	return operations.NewThreadGetOneOK().WithPayload(&thread)
}

// ThreadGetPosts ... OK
func (dbManager ForumPgSQL) ThreadGetPosts(params operations.ThreadGetPostsParams) middleware.Responder {
	tx := dbManager.db.MustBegin()

	threadID := ID{}
	posts := models.Posts{}

	slug, id := SlugID(params.SlugOrID)
	errNotFound := tx.Get(&threadID, `SELECT id FROM threads WHERE lower(slug) = lower($1) OR id = $2`, slug, id)
	if errNotFound != nil {
		log.Println(errNotFound)
		tx.Rollback()
		return operations.NewThreadGetPostsNotFound().WithPayload(&models.Error{Message: ERR_NOT_FOUND})
	}

	query := `SELECT id, forum, thread, author, created, is_edited as isEdited, message, parent from posts WHERE thread = $1`

	desc := params.Desc != nil && *params.Desc
	limit := strconv.FormatInt(int64(*params.Limit), 10)

	if *params.Sort == "flat" {
		if params.Since != nil {
			query += ` AND id `
			if desc {
				query += ` < $2`
			} else {
				query += ` > $2`
			}
		}

		if desc {
			query += ` ORDER BY created DESC, id DESC`
		} else {
			query += ` ORDER BY created, id`
		}

		if params.Limit != nil {
			query += ` LIMIT ` + limit
		}

		if params.Since != nil {
			err := tx.Select(&posts, query, threadID.ID, params.Since)
			if err != nil {
				log.Println(err)
				tx.Rollback()
				return operations.NewThreadGetPostsNotFound().WithPayload(&models.Error{Message: ERR})
			}
		} else {
			err := tx.Select(&posts, query, threadID.ID)
			if err != nil {
				log.Println(err)
				tx.Rollback()
				return operations.NewThreadGetPostsNotFound().WithPayload(&models.Error{Message: ERR})
			}
		}
	}
	if *params.Sort == "tree" {
		if params.Since != nil {
			query += ` AND path `
			if desc {
				query += ` < `
			} else {
				query += ` > `
			}
			query += `(SELECT path FROM posts WHERE id = $2) `
		}

		query += ` ORDER BY string_to_array(subltree(posts.path, 0, 1)::text,'.')::integer[]`
		if desc {
			query += ` DESC `
		}
		query += `, string_to_array(posts.path::text,'.')::integer[] `
		if desc {
			query += ` DESC `
		}

		if params.Limit != nil {
			query += ` LIMIT ` + limit
		}

		if params.Since != nil {
			err := tx.Select(&posts, query, threadID.ID, params.Since)
			if err != nil {
				log.Println(err)
				tx.Rollback()
				return operations.NewThreadGetPostsNotFound().WithPayload(&models.Error{Message: ERR})
			}
		} else {
			err := tx.Select(&posts, query, threadID.ID)
			if err != nil {
				log.Println(err)
				tx.Rollback()
				return operations.NewThreadGetPostsNotFound().WithPayload(&models.Error{Message: ERR})
			}
		}
	}
	if *params.Sort == "parent_tree" {
		query += ` AND subltree(path, 0, 1) IN (SELECT p1.path FROM posts p1
			WHERE nlevel(p1.path) = 1 AND p1.thread = $1 `
		if params.Since != nil {
			query += ` AND path `
			if desc {
				query += ` < `
			} else {
				query += ` > `
			}
			query += `(SELECT p2.path FROM posts p2 WHERE p2.id = $2) `
		}

		query += ` ORDER BY string_to_array(subltree(posts.path, 0, 1)::text,'.')::integer[]`
		if desc {
			query += ` DESC `
		}
		query += `, string_to_array(p1.path::text,'.')::integer[] `
		if desc {
			query += ` DESC `
		}

		if params.Limit != nil {
			query += ` LIMIT ` + limit
		}

		query += `) ORDER BY string_to_array(posts.path::text,'.')::integer[] `
		if desc {
			query += ` DESC `
		}
		query += `, string_to_array(posts.path::text,'.')::integer[] `
		if desc {
			query += ` DESC `
		}

		log.Println(query)

		if params.Since != nil {
			err := tx.Select(&posts, query, threadID.ID, params.Since)
			if err != nil {
				log.Println(err)
				tx.Rollback()
				return operations.NewThreadGetPostsNotFound().WithPayload(&models.Error{Message: ERR})
			}
		} else {
			err := tx.Select(&posts, query, threadID.ID)
			if err != nil {
				log.Println(err)
				tx.Rollback()
				return operations.NewThreadGetPostsNotFound().WithPayload(&models.Error{Message: ERR})
			}
		}
	}

	tx.Commit()
	return operations.NewThreadGetPostsOK().WithPayload(posts)
}

// ThreadUpdate ... OK
func (dbManager ForumPgSQL) ThreadUpdate(params operations.ThreadUpdateParams) middleware.Responder {
	tx := dbManager.db.MustBegin()

	threadID := ID{}
	thread := models.Thread{}

	slug, id := SlugID(params.SlugOrID)
	err := tx.Get(&threadID, `SELECT id FROM threads WHERE slug = $1 OR id = $2`, slug, id)
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return operations.NewThreadUpdateNotFound().WithPayload(&models.Error{Message: ERR_NOT_FOUND})
	}

	errNotFound := tx.Get(&thread, `UPDATE threads SET message = COALESCE($1, message), title = COALESCE($2, title)
	WHERE id = $3 RETURNING *`,
		params.Thread.Message, params.Thread.Title, threadID.ID)

	if errNotFound != nil {
		log.Println(errNotFound)
		tx.Rollback()
		return operations.NewUserUpdateConflict().WithPayload(&models.Error{Message: ERR_NOT_FOUND})
	}

	tx.Commit()
	return operations.NewThreadUpdateOK().WithPayload(&thread)
}

// ThreadVote ... OK
func (dbManager ForumPgSQL) ThreadVote(params operations.ThreadVoteParams) middleware.Responder {
	tx := dbManager.db.MustBegin()

	thread := models.Thread{}
	threadID := ID{}
	voteID := ID{}

	slug, id := SlugID(params.SlugOrID)
	err := tx.Get(&threadID, `SELECT id FROM threads WHERE slug = $1 OR id = $2`, slug, id)
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return operations.NewThreadVoteNotFound().WithPayload(&models.Error{Message: ERR_NOT_FOUND})
	}

	errExist := tx.Get(&voteID, `SELECT id FROM votes WHERE author = $1 AND thread = $2`, params.Vote.Nickname, threadID.ID)
	if errExist != nil {
		_, errAlreadyExists := tx.Exec(`INSERT INTO votes (voice, author, thread) VALUES ($1, $2, $3)`,
			params.Vote.Voice, params.Vote.Nickname, threadID.ID)

		if errAlreadyExists != nil {
			log.Println(errAlreadyExists)
		}
	} else {
		_, errNotFound := tx.Exec(`UPDATE votes SET voice = $1 WHERE author = $2 AND thread = $3`,
			params.Vote.Voice, params.Vote.Nickname, threadID.ID)

		if errNotFound != nil {
			log.Println(errNotFound)
		}
	}

	tx.Get(&thread, `UPDATE threads SET votes = (SELECT SUM(voice) FROM votes WHERE thread = $1)
                WHERE id = $1 RETURNING *`, threadID.ID)

	tx.Commit()
	return operations.NewThreadVoteOK().WithPayload(&thread)
}

//UserCreate ... OK OK
func (dbManager ForumPgSQL) UserCreate(params operations.UserCreateParams) middleware.Responder {
	tx := dbManager.db.MustBegin()

	user := models.User{}
	users := models.Users{}

	tx.Select(&users, "SELECT nickname, fullname, about, email FROM users WHERE lower(users.nickname) = lower($1) OR lower(users.email) = lower($2)", params.Nickname, params.Profile.Email)

	if len(users) != 0 {
		return operations.NewUserCreateConflict().WithPayload(users)
	}

	tx.Get(&user, "INSERT INTO users (nickname, fullname, about, email) VALUES ($1, $2, $3, $4) RETURNING nickname, fullname, about, email",
		params.Nickname, params.Profile.Fullname, params.Profile.About, params.Profile.Email)

	tx.Commit()
	return operations.NewUserCreateCreated().WithPayload(&user)
}

//UserGetOne ... OK OK
func (dbManager ForumPgSQL) UserGetOne(params operations.UserGetOneParams) middleware.Responder {
	tx := dbManager.db.MustBegin()
	defer tx.Rollback()

	users := []models.User{}
	check(tx.Select(&users, "SELECT nickname, fullname, about, email FROM users WHERE lower(users.nickname) = lower($1)", params.Nickname))
	check(tx.Commit())

	if len(users) == 0 {
		return operations.NewUserGetOneNotFound().WithPayload(&models.Error{Message: ERR_NOT_FOUND})
	}

	return operations.NewUserGetOneOK().WithPayload(&users[0])
}

//UserUpdate ... OK OK
func (dbManager ForumPgSQL) UserUpdate(params operations.UserUpdateParams) middleware.Responder {
	tx := dbManager.db.MustBegin()

	user := models.User{}
	users := models.Users{}

	check(tx.Select(&users, `SELECT nickname FROM users
		WHERE lower(users.nickname) = lower($1) OR lower(users.email) = COALESCE(lower($2), email)`, params.Nickname, params.Profile.Email))
	if len(users) == 0 {
		tx.Rollback()
		return operations.NewUserUpdateNotFound().WithPayload(&models.Error{Message: ERR_NOT_FOUND})
	}
	if len(users) > 1 {
		tx.Rollback()
		return operations.NewUserUpdateConflict().WithPayload(&models.Error{Message: ERR_ALREADY_EXISTS})
	}

	if params.Profile == nil {
		return operations.NewUserUpdateOK().WithPayload(users[0])
	}

	query := `UPDATE users SET nickname = nickname`
	if params.Profile.Fullname != "" {
		query += `, fullname = '` + params.Profile.Fullname + `'`
	}
	if params.Profile.Email != "" {
		query += `, email = '` + params.Profile.Email.String() + `'`
	}
	if params.Profile.About != "" {
		query += `, about = '` + params.Profile.About + `'`
	}
	query += ` WHERE lower(nickname) = lower($1) RETURNING about, email, fullname, nickname`

	tx.Get(&user, query, params.Nickname)

	tx.Commit()
	return operations.NewUserUpdateOK().WithPayload(&user)
}

func SlugID(slugOrID string) (string, int64) {
	id, err := strconv.ParseInt(slugOrID, 10, 64)
	slug := slugOrID
	if err != nil {
		id = -1
	}
	return slug, id
}
