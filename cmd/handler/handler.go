package handler

import (
	"net/http"

	//"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"github.com/rs/cors"
	"go.uber.org/zap"
)

const webDir = "/app/web"

// StartHandler to provide http api
func StartHandler() error {
	router := mux.NewRouter()

	// register all the routes
	registerStatusRoutes(router)
	registerJobRoutes(router)
	registerGeneratorRoutes(router)
	registerQueryRoutes(router)

	// register static files
	fs := http.FileServer(http.Dir(webDir))
	//fsGzip := handlers.CompressHandler(fs)
	router.PathPrefix("/").Handler(fs)

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
		// Enable Debugging for testing, consider disabling in production
		Debug: false,
	})

	// Insert the middleware
	handler := c.Handler(router)

	addr := ":8080"
	zap.L().Info("Listening HTTP service on", zap.String("address", addr), zap.String("webDirectory", webDir))
	return http.ListenAndServe(addr, handler)
}
