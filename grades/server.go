package grades

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func RegisterHandler() {
	handler := new(studentsHandler)
	http.Handle("/students/", handler)
	http.Handle("/students", handler)
}

type studentsHandler struct{}

// /students
// /students/{:id}
// /students/{:id}/grades
func (handler studentsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	pathSegments := strings.Split(r.URL.Path, "/")
	switch len(pathSegments) {
	case 2:
		handler.getAll(w, r)
	case 3:
		id, err := strconv.Atoi(pathSegments[2])
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		handler.getOne(w, r, id)
	case 4:
		id, err := strconv.Atoi(pathSegments[2])
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if strings.ToLower(pathSegments[3]) != "grades" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		handler.addGrade(w, r, id)
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

func (handler studentsHandler) getAll(w http.ResponseWriter, r *http.Request) {
	gradeMutex.Lock()
	defer gradeMutex.Unlock()

	data, err := handler.toJSON(students)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write(data)
}

func (handler studentsHandler) getOne(w http.ResponseWriter, r *http.Request, id int) {
	gradeMutex.Lock()
	defer gradeMutex.Unlock()

	student, err := students.GetByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Println(err)
		return
	}

	data, err := handler.toJSON(student)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(data)
}

func (handler studentsHandler) addGrade(w http.ResponseWriter, r *http.Request, id int) {
	gradeMutex.Lock()
	defer gradeMutex.Unlock()

	var grade Grade
	err := json.NewDecoder(r.Body).Decode(&grade)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("cannot seirialize data: %v", err)
		return
	}

	student, err := students.GetByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Println(err)
		return
	}

	student.Grades = append(student.Grades, grade)
	w.WriteHeader(http.StatusCreated)

	data, err := handler.toJSON(grade)
	if err != nil {
		log.Println(err)
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(data)
}

func (hander studentsHandler) toJSON(v interface{}) ([]byte, error) {

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(v)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize students: %w", err)
	}

	return buf.Bytes(), nil
}
