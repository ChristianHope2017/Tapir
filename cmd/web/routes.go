package main

import (
	"net/http"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("GET /{$}", app.home)
	mux.HandleFunc("GET /getfeedback", app.getfeedback)
	mux.HandleFunc("POST /feedback/new", app.createFeedback)
	mux.HandleFunc("GET /feedbacks", app.feedbackSuccess)

	mux.HandleFunc("GET /getjournal", app.getjournal)
	mux.HandleFunc("POST /journal/new", app.createJournal)
	mux.HandleFunc("GET /journals", app.journalSuccess)

	mux.HandleFunc("GET /gettodo", app.gettodo)
	mux.HandleFunc("POST /todo/new", app.createTodo)
	mux.HandleFunc("GET /todos", app.todoSuccess)

	return app.loggingMiddleware(mux)
}
