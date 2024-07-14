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

func RegisterHandlers () {
	handler := new(studentHandler)
	http.Handle("/students", handler)
	http.Handle("/students/", handler)
}

type studentHandler struct {}

func (sh studentHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	pathSegments := strings.Split(r.URL.Path, "/")
	switch len(pathSegments) {
		case 2:
			sh.getAll(w,r)
		case 3:
			id, err := strconv.Atoi(pathSegments[2]) // 将字符串表示的数字转换为对应的整数
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			sh.getOne(w, r, id)
		case 4: 
			id, err := strconv.Atoi(pathSegments[2])
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			sh.addGrade(w, r, id)
		default:
			w.WriteHeader(http.StatusNotFound)
	}
}

func (sh studentHandler) getAll(w http.ResponseWriter, r *http.Request) {
	studentMutex.Lock()
	defer studentMutex.Unlock()

	data, err := sh.toJSON(students)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Panicln(err)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write(data)
}

func (sh studentHandler) getOne(w http.ResponseWriter, r *http.Request, id int) {
	studentMutex.Lock()
	defer studentMutex.Unlock()

	student, err := students.GetByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Panicln(err)
		return
	}

	data, err := sh.toJSON(student)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Panicf("Failed to serialize student: %q", err)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write(data)
}

func (sh studentHandler) addGrade(w http.ResponseWriter, r *http.Request, id int) {
	studentMutex.Lock()
	defer studentMutex.Unlock()

	student, err := students.GetByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Panicln(err)
		return
	}
	var g Grade
	dec := json.NewDecoder(r.Body)
	err = dec.Decode(&g)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Panicln(err)
		return
	}

	student.Grades = append(student.Grades, g)
	w.WriteHeader(http.StatusCreated)
	data, err := sh.toJSON(student)
	if err != nil {
		log.Panicln(err)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write(data)
}

func (sh studentHandler) toJSON(obj interface{}) ([]byte, error) {
	var b bytes.Buffer
	enc := json.NewEncoder(&b)
	err := enc.Encode(obj)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize students: %q", err)
	}
	return b.Bytes(), nil
}