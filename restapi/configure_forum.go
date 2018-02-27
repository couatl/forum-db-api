package restapi

import (
	"crypto/tls"
	"log"
	"net/http"

	"github.com/dre1080/recover"
	errors "github.com/go-openapi/errors"
	runtime "github.com/go-openapi/runtime"
	graceful "github.com/tylerb/graceful"

	"github.com/couatl/forum-db-api/modules/service"
	"github.com/couatl/forum-db-api/restapi/operations"
	"github.com/go-openapi/swag"
)

//go:generate swagger generate server --target .. --name forum --spec ../swagger.yml
//go:generate go-bindata -pkg assets_db -o ../modules/assets/assets_db/assets_db.go -prefix ../modules/assets/ ../modules/assets/...

type DatabaseFlags struct {
	Database string `long:"database" description:"database connection parameters"`
}

var dbFlags DatabaseFlags

func configureFlags(api *operations.ForumAPI) {
	api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{
		{"database", "database connection parameters", &dbFlags},
	}
}

func configureAPI(api *operations.ForumAPI) http.Handler {
	api.ServeError = errors.ServeError

	api.JSONConsumer = runtime.JSONConsumer()

	api.BinConsumer = runtime.ByteStreamConsumer()

	api.JSONProducer = runtime.JSONProducer()

	var handler service.ForumHandler = service.NewForum(dbFlags.Database)

	api.ClearHandler = operations.ClearHandlerFunc(handler.Clear)
	api.StatusHandler = operations.StatusHandlerFunc(handler.Status)

	api.ForumCreateHandler = operations.ForumCreateHandlerFunc(handler.ForumCreate)
	api.ForumGetOneHandler = operations.ForumGetOneHandlerFunc(handler.ForumGetOne)
	api.ForumGetThreadsHandler = operations.ForumGetThreadsHandlerFunc(handler.ForumGetThreads)
	api.ForumGetUsersHandler = operations.ForumGetUsersHandlerFunc(handler.ForumGetUsers)

	api.PostGetOneHandler = operations.PostGetOneHandlerFunc(handler.PostGetOne)
	api.PostUpdateHandler = operations.PostUpdateHandlerFunc(handler.PostUpdate)
	api.PostsCreateHandler = operations.PostsCreateHandlerFunc(handler.PostsCreate)

	api.ThreadCreateHandler = operations.ThreadCreateHandlerFunc(handler.ThreadCreate)
	api.ThreadGetOneHandler = operations.ThreadGetOneHandlerFunc(handler.ThreadGetOne)
	api.ThreadGetPostsHandler = operations.ThreadGetPostsHandlerFunc(handler.ThreadGetPosts)
	api.ThreadUpdateHandler = operations.ThreadUpdateHandlerFunc(handler.ThreadUpdate)
	api.ThreadVoteHandler = operations.ThreadVoteHandlerFunc(handler.ThreadVote)

	api.UserCreateHandler = operations.UserCreateHandlerFunc(handler.UserCreate)
	api.UserGetOneHandler = operations.UserGetOneHandlerFunc(handler.UserGetOne)
	api.UserUpdateHandler = operations.UserUpdateHandlerFunc(handler.UserUpdate)

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
}

func configureServer(s *graceful.Server, scheme, addr string) {
}

func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

func setupGlobalMiddleware(handler http.Handler) http.Handler {
	recovery := recover.New(&recover.Options{
		Log: log.Print,
	})
	return recovery(uiMiddleware(handler))
}

func uiMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/swagger.json" {
			handler.ServeHTTP(w, r)
			return
		}
		// Swagger UI
		if r.URL.Path == "/api/" {
			r.URL.Path = "/api"
		}
		handler.ServeHTTP(w, r)
	})
}
