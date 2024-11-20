package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Tasks
type Task struct {
	ID           string   `json:"id"`
	Description  string   `json:"description"`
	Note         string   `json:"note"`
	Applications []string `json:"applications"`
}

var tasks = map[string]Task{
	"1": {
		ID:          "1",
		Description: "Сделать финальное задание темы REST API",
		Note:        "Если сегодня сделаю, то завтра будет свободный день. Ура!",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
		},
	},
	"2": {
		ID:          "2",
		Description: "Протестировать финальное задание с помощью Postmen",
		Note:        "Лучше это делать в процессе разработки, каждый раз, когда запускаешь сервер и проверяешь хендлер",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
			"Postman",
		},
	},
}

// Обработчик для получения всех задач
func getAllTasksHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(tasks); err != nil {
		http.Error(w, "Ошибка при получении задач", http.StatusInternalServerError)
	}
}

// Обработчик для отправки задачи на сервер
func createTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var task Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "Неверный формат задачи", http.StatusBadRequest)
		return
	}

	if _, exists := tasks[task.ID]; exists {
		http.Error(w, "Задача с таким ID уже существует", http.StatusBadRequest)
		return
	}

	tasks[task.ID] = task
	w.WriteHeader(http.StatusCreated)
}

// Обработчик для получения задачи по ID
func getTaskByIDHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := chi.URLParam(r, "id")
	task, exists := tasks[id]
	if !exists {
		http.Error(w, "Задача не найдена", http.StatusBadRequest)
		return
	}

	if err := json.NewEncoder(w).Encode(task); err != nil {
		http.Error(w, "Ошибка при отправке задачи", http.StatusInternalServerError)
	}
}

// Обработчик для удаления задачи по ID
func deleteTaskByIDHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := chi.URLParam(r, "id")
	if _, exists := tasks[id]; !exists {
		http.Error(w, "Задача не найдена", http.StatusBadRequest)
		return
	}

	delete(tasks, id)
	w.WriteHeader(http.StatusOK)
}

func main() {
	r := chi.NewRouter()

	// Регистрация обработчиков
	r.Get("/tasks", getAllTasksHandler)
	r.Post("/tasks", createTaskHandler)
	r.Get("/tasks/{id}", getTaskByIDHandler)
	r.Delete("/tasks/{id}", deleteTaskByIDHandler)

	fmt.Println("Сервер запущен на порту 8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
	}
}
