package main

import (
	"net/http"

	"github.com/ChristianHope2017/di/internal/data"

	"github.com/ChristianHope2017/di/internal/validator"

	"fmt"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	data := NewTemplateData()
	data.Title = "Welcome"
	data.HeaderText = "We are here to help"
	err := app.render(w, http.StatusOK, "home.tmpl", data)
	if err != nil {
		app.logger.Error("failed to render home page", "template", "home.tmpl", "error", err, "url", r.URL.Path, "method", r.Method)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (app *application) getfeedback(w http.ResponseWriter, r *http.Request) {
	data := NewTemplateData()
	data.Title = "Welcome"
	data.HeaderText = "We are here to help"
	err := app.render(w, http.StatusOK, "getfeedback.tmpl", data)
	if err != nil {
		app.logger.Error("failed to render getfeedback page", "template", "getfeedback.tmpl", "error", err, "url", r.URL.Path, "method", r.Method)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (app *application) getjournal(w http.ResponseWriter, r *http.Request) {
	data := NewTemplateData()
	data.Title = "Welcome"
	data.HeaderText = "We are here to help"
	err := app.render(w, http.StatusOK, "getjournal.tmpl", data)
	if err != nil {
		app.logger.Error("failed to render getjournal page", "template", "getjournal.tmpl", "error", err, "url", r.URL.Path, "method", r.Method)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (app *application) gettodo(w http.ResponseWriter, r *http.Request) {
	data := NewTemplateData()
	data.Title = "Welcome"
	data.HeaderText = "We are here to help"
	err := app.render(w, http.StatusOK, "gettodo.tmpl", data)
	if err != nil {
		app.logger.Error("failed to render gettodo page", "template", "gettodo.tmpl", "error", err, "url", r.URL.Path, "method", r.Method)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (app *application) createFeedback(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		app.logger.Error("failed to parse form", "error", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	name := r.PostForm.Get("name")
	email := r.PostForm.Get("email")
	subject := r.PostForm.Get("subject")
	message := r.PostForm.Get("message")

	feedback := &data.Feedback{
		Fullname: name,
		Email:    email,
		Subject:  subject,
		Message:  message,
	}

	// validate data
	v := validator.NewValidator()
	data.ValidateFeedback(v, feedback)
	// Check for validation errors
	if !v.ValidData() {
		data := NewTemplateData()
		data.Title = "Welcome"
		data.HeaderText = "We are here to help"
		data.FormErrors = v.Errors
		data.FormData = map[string]string{
			"name":    name,
			"email":   email,
			"subject": subject,
			"message": message,
		}

		err := app.render(w, http.StatusUnprocessableEntity, "getfeedback.tmpl", data)
		if err != nil {
			app.logger.Error("failed to render getfeedback page", "template", "getfeedback.tmpl", "error", err, "url", r.URL.Path, "method", r.Method)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		return
	}

	err = app.feedback.Insert(feedback)
	if err != nil {
		app.logger.Error("failed to insert feedback", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/feedbacks", http.StatusSeeOther)
}

func (app *application) feedbackSuccess(w http.ResponseWriter, r *http.Request) {
	// Fetch feedback data from database
	feedbacks, err := app.feedback.GetAll()
	if err != nil {
		app.logger.Error("failed to fetch feedback from database", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Convert database records to formatted strings
	var feedbackStrings []string
	for _, Feedback := range feedbacks {
		feedbackStr := fmt.Sprintf(
			"ID: %v | Created At: %s | Name: %s | Email: %s | Subject: %s | Message: %s",
			Feedback.ID,
			Feedback.CreatedAt.Format("2006-01-02 15:04:05"),
			Feedback.Fullname,
			Feedback.Email,
			Feedback.Subject,
			Feedback.Message,
		)
		feedbackStrings = append(feedbackStrings, feedbackStr)
	}

	// Prepare template data
	data := NewTemplateData()
	data.Title = "Feedback Submitted"
	data.HeaderText = "Thank You for Your Feedback!"
	data.Feedback = feedbackStrings

	err = app.render(w, http.StatusOK, "feedbacks.tmpl", data)
	if err != nil {
		app.logger.Error("failed to render feedback success page", "template", "feedbacks.tmpl", "error", err, "url", r.URL.Path, "method", r.Method)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (app *application) createJournal(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.logger.Error("failed to parse form", "error", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")

	journal := &data.Journal{
		Title:   title,
		Content: content,
	}

	// validate data
	v := validator.NewValidator()
	data.ValidateJournal(v, journal)
	// Check for validation errors
	if !v.ValidData() {
		data := NewTemplateData()
		data.Title = "Welcome"
		data.HeaderText = "We are here to help"
		data.FormErrors = v.Errors
		data.FormData = map[string]string{
			"title":   title,
			"content": content,
		}

		err := app.render(w, http.StatusUnprocessableEntity, "getjournal.tmpl", data)
		if err != nil {
			app.logger.Error("failed to render getjournal page", "template", "getjournal.tmpl", "error", err, "url", r.URL.Path, "method", r.Method)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		return
	}

	err = app.journal.Insert(journal)
	if err != nil {
		app.logger.Error("failed to insert journal entry", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/journals", http.StatusSeeOther)
}

func (app *application) journalSuccess(w http.ResponseWriter, r *http.Request) {
	// Fetch journal data from database
	journals, err := app.journal.GetAll()
	if err != nil {
		app.logger.Error("failed to fetch journal from database", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Convert database records to formatted strings
	var journalStrings []string
	for _, Journal := range journals {
		journalStr := fmt.Sprintf(
			"ID: %v | Created At: %s | Title: %s | Content: %s",
			Journal.ID,
			Journal.CreatedAt.Format("2006-01-02 15:04:05"),
			Journal.Title,
			Journal.Content,
		)
		journalStrings = append(journalStrings, journalStr)
	}

	// Prepare template data
	data := NewTemplateData()
	data.Title = "Journal Entry Submitted"
	data.HeaderText = "Thank You for Your Journal Entry!"
	data.Journal = journalStrings

	err = app.render(w, http.StatusOK, "journals.tmpl", data)
	if err != nil {
		app.logger.Error("failed to render journal success page", "template", "journals.tmpl", "error", err, "url", r.URL.Path, "method", r.Method)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (app *application) createTodo(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.logger.Error("failed to parse form", "error", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	title := r.PostForm.Get("title")
	task := r.PostForm.Get("task")

	todo := &data.Todo{
		Title: title,
		Task:  task,
	}

	// validate data
	v := validator.NewValidator()
	data.ValidateTodo(v, todo)
	// Check for validation errors
	if !v.ValidData() {
		data := NewTemplateData()
		data.Title = "Welcome"
		data.HeaderText = "We are here to help"
		data.FormErrors = v.Errors
		data.FormData = map[string]string{
			"title": title,
			"task":  task,
		}

		err := app.render(w, http.StatusUnprocessableEntity, "gettodo.tmpl", data)
		if err != nil {
			app.logger.Error("failed to render gettodo page", "template", "gettodo.tmpl", "error", err, "url", r.URL.Path, "method", r.Method)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		return
	}

	err = app.todo.Insert(todo)
	if err != nil {
		app.logger.Error("failed to insert todo item", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/todos", http.StatusSeeOther)
}

func (app *application) todoSuccess(w http.ResponseWriter, r *http.Request) {
	// Fetch todo data from database
	todos, err := app.todo.Getall()
	if err != nil {
		app.logger.Error("failed to fetch todo from database", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Convert database records to formatted strings
	var todoStrings []string
	for _, Todo := range todos {
		todoStr := fmt.Sprintf(
			"ID: %v | Created At: %s | Title: %s | Task: %s",
			Todo.ID,
			Todo.CreatedAt.Format("2006-01-02 15:04:05"),
			Todo.Title,
			Todo.Task,
		)
		todoStrings = append(todoStrings, todoStr)
	}

	// Prepare template data
	data := NewTemplateData()
	data.Title = "Todo Item Submitted"
	data.HeaderText = "Thank You for Your Todo Item!"
	data.Todo = todoStrings

	err = app.render(w, http.StatusOK, "todos.tmpl", data)
	if err != nil {
		app.logger.Error("failed to render todo success page", "template", "todos.tmpl", "error", err, "url", r.URL.Path, "method", r.Method)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
