package app

import (
	"context"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"

	"github.com/Southclaws/scribble.rs/src/game"
)

// Config represents environment variable configuration parameters
type Config struct {
	ListenAddr string `default:"0.0.0.0:8080"          split_words:"true"`
}

// App stores root application state
type App struct {
	config Config
	server http.Server
	ctx    context.Context
	cancel context.CancelFunc
}

func Initialise(root context.Context) (app *App, err error) {
	app = &App{config: Config{}}
	if err = envconfig.Process("", &app.config); err != nil {
		return
	}

	app.ctx, app.cancel = WithSignal(root)

	g := game.New()

	router := mux.NewRouter()
	mount(router, "/api/game", g.Routes())

	router.HandleFunc("/{rest:[a-zA-Z0-9=\\-\\/]+}", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("no module found for that route"))
	})

	app.server = http.Server{
		Handler: router,
		Addr:    "0.0.0.0:8080",
	}

	return
}

// Start starts the application and blocks until fatal error
// The server will shut down if the root context is cancelled
func (app *App) Start() error {
	go func() { app.server.ListenAndServe() }()

	<-app.ctx.Done()

	zap.L().Error("server context cancelled", zap.Error(app.ctx.Err()))

	return app.server.Shutdown(context.Background())
}

func mount(r *mux.Router, path string, handler http.Handler) {
	r.PathPrefix(path).Handler(
		http.StripPrefix(strings.TrimSuffix(path, "/"), handler),
	)
}
