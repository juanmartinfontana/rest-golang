package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

/*task structura de una task*/
type task struct {
	ID      int    `json:ID`
	Name    string `json:Name`
	Content string `json:Content`
}

/*allTasks lista de tasks */
type allTasks []task

var tasks = allTasks{
	{
		ID:      1,
		Name:    "Task One",
		Content: "Some Content",
	},
}

/*createTask funcion que crea una Task */
func createTask(w http.ResponseWriter, r *http.Request) {
	var newTask task
	reqBody, err := ioutil.ReadAll(r.Body)

	if err != nil {
		fmt.Fprintf(w, "Insert a Valid task")
	}

	json.Unmarshal(reqBody, &newTask)
	newTask.ID = len(tasks) + 1
	tasks = append(tasks, newTask)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTask)

}

/*getTaskById obtine una task mediante un Id */
func getTaskById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskID, err := strconv.Atoi(vars["id"])

	if err != nil {
		fmt.Fprintf(w, "Invalid ID")
		return
	}

	for _, task := range tasks {
		if task.ID == taskID {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(task)
		}
	}
}

/*getTasks funcion que devuelve todas las tasks */
func getTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

/*deleteTask funcion que elimina una task mediante un id */
func deleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskID, err := strconv.Atoi(vars["id"])

	if err != nil {
		fmt.Fprintf(w, "Invalid ID")
		return
	}

	for i, t := range tasks {
		if t.ID == taskID {
			tasks = append(tasks[:i], tasks[i+1:]...)
			fmt.Fprintf(w, "The trask with ID %v has been romeve succesfully", taskID)
		}
	}
}

/*updateTask funcion que actualiza una task */
func updateTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskID, err := strconv.Atoi(vars["id"])
	var updatedTask task

	if err != nil {
		fmt.Fprintf(w, "Invalid ID")
		return
	}

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Please enter valid data")
	}
	json.Unmarshal(reqBody, &updatedTask)

	for i, t := range tasks {
		if t.ID == taskID {
			tasks = append(tasks[:i], tasks[i+1:]...)
			updatedTask.ID = taskID
			tasks = append(tasks, updatedTask)

			fmt.Fprintf(w, "The task with ID %v has been updated successfully", taskID)
		}
	}
}

func indexRoute(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to my API")

}

func main() {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", indexRoute)
	router.HandleFunc("/tasks", getTasks).Methods("GET")
	router.HandleFunc("/tasks", createTask).Methods("POST")
	router.HandleFunc("/tasks/{id}", getTaskById).Methods("GET")
	router.HandleFunc("/tasks/{id}", deleteTask).Methods("DELETE")
	router.HandleFunc("/tasks/{id}", updateTask).Methods("PUT")

	log.Fatal(http.ListenAndServe(":3000", router))
}
