package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"golang/internal/api/helpers"
	"golang/models"
	"golang/sqlconnect"
	"net/http"
	"strconv"
)

func AddStudentHandler(writer http.ResponseWriter, request *http.Request) {
	db, err := sqlconnect.ConnectToDB()
	if err != nil {
		http.Error(writer, "Error Connecting DB", http.StatusBadGateway)
		return
	}
	defer db.Close()

	// check if its an valid req
	var newStudent models.Student

	err = json.NewDecoder(request.Body).Decode(&newStudent)
	if err != nil {
		http.Error(writer, "Not an Valid Request", http.StatusBadRequest)
		return
	}

	stmt, err := db.Prepare("INSERT INTO students (firstName,lastName,email,class) VALUES ($1, $2, $3, $4) RETURNING id")
	if err != nil {
		http.Error(writer, "Not able to consruct SQL Stmt stmt", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()
	var id int
	err = stmt.QueryRow(newStudent.FirstName, newStudent.LastName, newStudent.Email, newStudent.Class).Scan(&id)
	if err != nil {
		http.Error(writer, "Error inserting in db", http.StatusInternalServerError)
		return
	}

	resp := struct {
		Status string `json:"status"`
		ID     int    `json:"id"`
	}{
		Status: "Success",
		ID:     id,
	}
	json.NewEncoder(writer).Encode(resp)
	return
}

func GetStudentHandler(writer http.ResponseWriter, request *http.Request) {
	studentID := request.PathValue("id")
	fmt.Println("Student ID : ", studentID)

	stdID, err := strconv.Atoi(studentID)
	if err != nil {
		http.Error(writer, "Invalid ID", http.StatusBadRequest)
		return
	}

	db, err := sqlconnect.ConnectToDB()
	if err != nil {
		http.Error(writer, "Error Reaching DB", http.StatusInternalServerError)
		return
	}

	defer db.Close()

	query := "SELECT id,firstName,lastName,email,class FROM students WHERE id=$1"

	rows, err := db.Query(query, stdID)
	if err != nil {
		http.Error(writer, "Cannot Get Items From DB", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var student models.Student
	err = db.QueryRow(query, stdID).Scan(&student.ID, &student.FirstName, &student.LastName, &student.Email, &student.Class)
	if err == sql.ErrNoRows {
		http.Error(writer, "Student not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(writer, "Error fetching student", http.StatusInternalServerError)
		return
	}

	response := struct {
		Status string         `json:"status"`
		Data   models.Student `json:"data"`
	}{
		Status: "Success",
		Data:   student,
	}
	json.NewEncoder(writer).Encode(response)
	return

}

func GetAllStudentsHandler(writer http.ResponseWriter, request *http.Request) {
	db, err := sqlconnect.ConnectToDB()
	if err != nil {
		http.Error(writer, "Error Reaching DB", http.StatusInternalServerError)
		return
	}

	defer db.Close()

	studentsList := make([]models.Student, 0)

	var args []any

	query := "SELECT id,firstName,lastName,email,class FROM students WHERE 1=1"

	query, args = helpers.AddFilters(request, query, args)
	query = helpers.AddSortParams(request, query)

	rows, err := db.Query(query, args...)
	if err != nil {
		http.Error(writer, "Cannot Get Items From DB", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var student models.Student
		err := rows.Scan(&student.ID, &student.FirstName, &student.LastName, &student.Email, &student.Class)
		if err != nil {
			http.Error(writer, "Cannot Get Items From DB", http.StatusNotFound)
			return
		}
		studentsList = append(studentsList, student)
	}

	response := struct {
		Status string           `json:"status"`
		Count  int              `json:"count"`
		Data   []models.Student `json:"data"`
	}{
		Status: "Success",
		Count:  len(studentsList),
		Data:   studentsList,
	}
	json.NewEncoder(writer).Encode(response)
	return

}

func UpdateStudentHandler(writer http.ResponseWriter, request *http.Request) {
	studentID, err := strconv.Atoi(request.PathValue("id"))
	if err != nil {
		http.Error(writer, "Invalid ID", http.StatusBadRequest)
		return
	}

	fmt.Println("Student ID : ", studentID)
	var updatedStudent models.StudentPatch

	err = json.NewDecoder(request.Body).Decode(&updatedStudent)
	if err != nil {
		http.Error(writer, "Invalid Body Format", http.StatusBadRequest)
		return
	}

	db, err := sqlconnect.ConnectToDB()
	if err != nil {
		http.Error(writer, "Error Reaching DB", http.StatusInternalServerError)
		return
	}

	defer db.Close()

	query := "SELECT id,firstName,lastName,email,class FROM students WHERE id=$1"

	rows, err := db.Query(query, studentID)
	if err != nil {
		http.Error(writer, "Cannot Get Items From DB", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var student models.Student
	err = db.QueryRow(query, studentID).Scan(&student.ID, &student.FirstName, &student.LastName, &student.Email, &student.Class)
	if err == sql.ErrNoRows {
		http.Error(writer, "Student not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(writer, "Error fetching student", http.StatusInternalServerError)
		return
	}

	if updatedStudent.FirstName != "" {
		student.FirstName = updatedStudent.FirstName
	}

	if updatedStudent.LastName != "" {
		student.LastName = updatedStudent.LastName
	}

	if updatedStudent.Email != "" {
		student.Email = updatedStudent.Email
	}

	if updatedStudent.Class != "" {
		// Check if class exists in teachers table before updating (foreign key validation)
		var classExists bool
		err = db.QueryRow("SELECT EXISTS (SELECT 1 FROM teachers WHERE class = $1)", updatedStudent.Class).Scan(&classExists)
		if err != nil {
			http.Error(writer, "Error checking class existence", http.StatusInternalServerError)
			return
		}
		if !classExists {
			http.Error(writer, "Class does not exist", http.StatusBadRequest)
			return
		}
		student.Class = updatedStudent.Class
	}

	// Perform the update query
	updateQuery := `
		UPDATE students 
		SET firstName=$1, lastName=$2, email=$3, class=$4
		WHERE id=$5
	`

	_, err = db.Exec(updateQuery, student.FirstName, student.LastName, student.Email, student.Class, student.ID)
	if err != nil {
		http.Error(writer, "Error updating student", http.StatusInternalServerError)
		return
	}

	response := struct {
		Status string         `json:"status"`
		Data   models.Student `json:"data"`
	}{
		Status: "Updated Successfully",
		Data:   student,
	}

	json.NewEncoder(writer).Encode(response)
}
