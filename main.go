package main

import (
	"github.com/rs/cors"
	"net/http"
)

func main() {
	InitDB()

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // Разрешить все домены
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})

	var srv http.Server
	srv.Handler = InitRoutes()
	srv.Handler = corsHandler.Handler(srv.Handler)
	//srv.TLSConfig = &tls.Config{
	//	MinVersion: tls.VersionTLS12,
	//}
	srv.Addr = "0.0.0.0:8887"
	srv.ListenAndServe()
}
