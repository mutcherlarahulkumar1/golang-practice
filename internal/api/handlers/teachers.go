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
)

func TeachersHandler(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("Hello Teachers Route")

	switch request.Method {
	case http.MethodGet:
		getTeachersHandler(writer, request)

	case http.MethodPost:
		addTeachersHandler(writer, request)

	case http.MethodPut:
		updateTeacherHandler(writer, request)

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

	teacherID := extractTeacherID(request)
	if teacherID == -1 {
		teachersList := make([]models.Teacher, 0)

		var args []any

		query := "SELECT id,firstName,lastName,email,class,subject FROM teachers WHERE 1=1"

		// Adding filters and sort params
		query, args = addFilters(request, query, args)
		query = addSortParams(request, query)

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
		return
	}

	var teacher models.Teacher
	err = db.QueryRow("SELECT id,firstName,lastName,email,class,subject FROM teachers WHERE id=$1", teacherID).Scan(&teacher.ID, &teacher.FirstName, &teacher.LastName, &teacher.Email, &teacher.Class, &teacher.Subject)
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

func extractTeacherID(request *http.Request) int {
	path := strings.TrimPrefix(request.URL.Path, "/teachers/")
	teacherID := strings.TrimSuffix(path, "/")
	id, err := strconv.Atoi(teacherID)
	if err != nil {
		return -1
	}
	return id
}

func addSortParams(request *http.Request, query string) string {
	sortParams := request.URL.Query()["sortby"]
	if len(sortParams) > 0 {
		query += " ORDER BY "
		for i, param := range sortParams {
			parts := strings.Split(param, ":")

			field, order := parts[0], parts[1]
			if i > 0 {
				query += ","
			}
			query += " " + field + " " + order
		}

	}
	return query
}

func addFilters(request *http.Request, query string, args []any) (string, []any) {
	params := map[string]string{
		"firstname": "firstName",
		"lastname":  "lastName",
		"email":     "email",
		"class":     "class",
		"subject":   "subject",
	}

	i := 1
	for param, dbField := range params {
		value := request.URL.Query().Get(param)
		if value != "" {
			query += fmt.Sprintf(" AND %s = $%d", dbField, i)
			args = append(args, value)
			i++
		}
	}
	return query, args
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

// PUT METHOD ->Updating the complete entry
func updateTeacherHandler(writer http.ResponseWriter, request *http.Request) {
	teacherID := extractTeacherID(request)
	if teacherID == -1 {
		http.Error(writer, "Invalid ID", http.StatusBadRequest)
		return
	}
	var updatedTeacher models.Teacher

	err := json.NewDecoder(request.Body).Decode(&updatedTeacher)
	if err != nil {
		http.Error(writer, "Invalid BOdy Format", http.StatusBadRequest)
		return
	}

	db, err := sqlconnect.ConnectToDB()
	if err != nil {
		http.Error(writer, "DB Error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var existingTeacher models.Teacher

	err = db.QueryRow("SELECT id,firstName,lastName,email,class,subject FROM teachers WHERE id=$1", teacherID).Scan(&existingTeacher.ID, &existingTeacher.FirstName, &existingTeacher.LastName, &existingTeacher.Email, &existingTeacher.Class, &existingTeacher.Subject)
	if err == sql.ErrNoRows {
		http.Error(writer, "No Data Found for requested ID", http.StatusFound)
		return
	}
	if err != nil {
		http.Error(writer, "Cannot Get Items From DB", http.StatusNotFound)
		return
	}

	updatedTeacher.ID = existingTeacher.ID

	_, err = db.Exec("UPDATE teachers SET firstName = $1, lastName = $2, email = $3, class = $4, subject = $5 where id = $6", updatedTeacher.FirstName, updatedTeacher.LastName, updatedTeacher.Email, updatedTeacher.Class, updatedTeacher.Subject, teacherID)
	if err != nil {
		http.Error(writer, "Error From DB", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(writer).Encode(updatedTeacher)
}
