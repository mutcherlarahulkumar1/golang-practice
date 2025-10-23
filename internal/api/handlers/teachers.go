package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"golang/models"
	"golang/sqlconnect"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

var (
	teachers = make(map[int]models.Teacher)
	mutex    = &sync.Mutex{}
	nextID   = 1
)

func init() {
	teachers[nextID] = models.Teacher{
		ID:        nextID,
		FirstName: "John",
		LastName:  "Doe",
		Class:     "2",
		Subject:   "Chemistry",
	}
	nextID++
	teachers[nextID] = models.Teacher{
		ID:        nextID,
		FirstName: "Jane",
		LastName:  "Smith",
		Class:     "10 A",
		Subject:   "Algebra",
	}
}

func TeachersHandler(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("Hello Teachers Route")

	switch request.Method {
	case http.MethodGet:
		getTeachersHandler(writer, request)

	case http.MethodPost:
		addTeachersHandler(writer, request)

	case http.MethodPut:
		writer.Write([]byte("THis PUT Call for teachers"))

	case http.MethodPatch:
		writer.Write([]byte("THis PATCH Call for teachers"))

	case http.MethodDelete:
		writer.Write([]byte("THis DELETE Call for teachers"))
	}

}

func getTeachersHandler(writer http.ResponseWriter, request *http.Request) {

	db, err := sqlconnect.ConnectToDB()
	if err != nil {
		http.Error(writer, "Error Connecting DB", http.StatusBadGateway)
		return
	}
	defer db.Close()

	path := strings.TrimPrefix(request.URL.Path, "/teachers/")
	teacherID := strings.TrimSuffix(path, "/")
	if teacherID == "" {
		teachersList := make([]models.Teacher, 0)
		lastName := request.URL.Query().Get("lastname")

		var args []any

		query := "SELECT id,firstName,lastName,email,class,subject FROM teachers WHERE 1=1"

		if lastName != "" {
			query += " AND lastName = $1"
			args = append(args, lastName)
		}

		var rows, err = db.Query(query, args...)
		if err != nil {
			http.Error(writer, "Cannot Get Items From DB", http.StatusNotFound)
			return
		}
		defer rows.Close()

		for rows.Next() {
			var teacher models.Teacher
			err := rows.Scan(&teacher.ID, &teacher.FirstName, &teacher.LastName, &teacher.Email, &teacher.Class, &teacher.Subject)
			if err != nil {
				http.Error(writer, "Cannot Get Items From DB", http.StatusNotFound)
				return
			}
			teachersList = append(teachersList, teacher)
		}

		response := struct {
			Status string           `json:"status"`
			Count  int              `json:"count"`
			Data   []models.Teacher `json:"data"`
		}{
			Status: "Success",
			Count:  len(teachersList),
			Data:   teachersList,
		}
		json.NewEncoder(writer).Encode(response)
	}
	id, err := strconv.Atoi(teacherID)
	if err != nil {
		http.Error(writer, "Invalid ID Format", http.StatusNotFound)
		return
	}

	var teacher models.Teacher
	err = db.QueryRow("SELECT id,firstName,lastName,email,class,subject FROM teachers WHERE id=$1", id).Scan(&teacher.ID, &teacher.FirstName, &teacher.LastName, &teacher.Email, &teacher.Class, &teacher.Subject)
	if err == sql.ErrNoRows {
		http.Error(writer, "No Data Found for requested ID", http.StatusFound)
		return
	}
	if err != nil {
		http.Error(writer, "Cannot Get Items From DB", http.StatusNotFound)
		return
	}
	json.NewEncoder(writer).Encode(teacher)

}

func addTeachersHandler(writer http.ResponseWriter, request *http.Request) {
	db, err := sqlconnect.ConnectToDB()
	if err != nil {
		http.Error(writer, "Error Connecting DB", http.StatusBadGateway)
		return
	}
	defer db.Close()

	// check if its an valid req
	var newTeachers []models.Teacher

	err = json.NewDecoder(request.Body).Decode(&newTeachers)
	if err != nil {
		http.Error(writer, "Not an Valid Request", http.StatusBadRequest)
		return
	}

	stmt, err := db.Prepare("INSERT INTO teachers (firstName,lastName,email,class,subject) VALUES ($1, $2, $3, $4, $5) RETURNING id")
	if err != nil {
		http.Error(writer, "Not able to prepare stmt", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	addedTeachers := make([]models.Teacher, 0, len(newTeachers))
	for _, newTeacher := range newTeachers {
		var id int
		err := stmt.QueryRow(newTeacher.FirstName, newTeacher.LastName, newTeacher.Email, newTeacher.Class, newTeacher.Subject).Scan(&id)
		if err != nil {
			http.Error(writer, "Error posting db", http.StatusInternalServerError)
			return
		}
		newTeacher.ID = id
		addedTeachers = append(addedTeachers, newTeacher)
	}

	resp := struct {
		Status string           `json:"status"`
		Count  int              `json:"count"`
		Data   []models.Teacher `json:"data"`
	}{
		Status: "Success",
		Count:  len(addedTeachers),
		Data:   addedTeachers,
	}
	json.NewEncoder(writer).Encode(resp)
	return
}
