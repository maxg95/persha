package main
;
import (
	"net/http"

	"persha.maxg95/assets"

	"github.com/gorilla/mux"
)

func (app *application) routes() http.Handler {
	mux := mux.NewRouter()
	mux.NotFoundHandler = http.HandlerFunc(app.notFound)

	mux.Use(app.recoverPanic)
	mux.Use(app.securityHeaders)

	fileServer := http.FileServer(http.FS(assets.EmbeddedFiles))
	mux.PathPrefix("/static/").Handler(fileServer)

	mux.HandleFunc("/", app.home).Methods("GET")
	mux.HandleFunc("/", app.insertRequest).Methods("POST")
	mux.HandleFunc("/pomynky", app.pomynky).Methods("GET")
	mux.HandleFunc("/pomynalni_obidy", app.pomynalni_obidy).Methods("GET")
	mux.HandleFunc("/pomynalni_obidy_lutsk", app.pomynalni_obidy_lutsk).Methods("GET")
	mux.HandleFunc("/vesillia", app.vesillia).Methods("GET")
	mux.HandleFunc("/vesillia_cafe", app.vesillia_cafe).Methods("GET")
	mux.HandleFunc("/vesillia_lutsk", app.vesillia_lutsk).Methods("GET")
	mux.HandleFunc("/keiterinh", app.keiterinh).Methods("GET")
	mux.HandleFunc("/keiterinh_cafe", app.keiterinh_cafe).Methods("GET")
	mux.HandleFunc("/keiterinh_lutsk", app.keiterinh_lutsk).Methods("GET")

	return mux
}
