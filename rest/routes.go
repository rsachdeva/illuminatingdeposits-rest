package rest

import (
	"log"
	"net/http"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/rsachdeva/illuminatingdeposits/database"
	"github.com/rsachdeva/illuminatingdeposits/invest"
	"github.com/rsachdeva/illuminatingdeposits/middleware"
	"github.com/rsachdeva/illuminatingdeposits/service"
	"github.com/rsachdeva/illuminatingdeposits/user"
)

func RegisterRoutesHandlers(server *http.Server, log *log.Logger, db *sqlx.DB, shutdownCh chan os.Signal) {
	// Construct the web.App which holds all routes as well as common Middleware.
	app := service.NewApp(shutdownCh, log, middleware.Logger(log), middleware.Errors(log), middleware.Metrics(), middleware.Panics(log))

	database.RegisterCheckHandler(db, app)

	user.RegisterUserHandler(db, app)

	invest.RegisterInvestHandler(log, app)

	server.Handler = app
}
