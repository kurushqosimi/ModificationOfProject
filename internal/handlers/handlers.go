package handlers

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"log"
	"main/internal/services"
	"main/pkg/models"
	"net/http"
	"strconv"
	"time"
)

type Handler struct {
	Service *services.Service
}

func NewHandler(service *services.Service) *Handler {
	return &Handler{Service: service}
}
func (h *Handler) UserCheck(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		var user models.User
		user.Username = request.Header.Get("login")
		user.Password = request.Header.Get("password")
		user.Active = true
		log.Println(user)
		//id := h.Service.Repository.CheckUsers(&user)
		id, err := h.Service.HashingThePassword(&user)
		if err != nil {
			log.Println(err)
			response.WriteHeader(http.StatusBadRequest)
			return
		}
		if id == 0 {
			log.Println("id не приходит")
			return
		}
		//log.Println(id)
		ctx := context.WithValue(request.Context(), "id", id)
		request = request.WithContext(ctx)
		next.ServeHTTP(response, request)
		response.Header().Set("Custom Header", user.Username)
	})
}
func (h *Handler) Create(response http.ResponseWriter, request *http.Request) { //todo
	var note models.Notes
	bytes, err := io.ReadAll(request.Body)
	if err != nil {
		log.Println(err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(bytes, &note)
	if err != nil {
		log.Println(err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	note.Active = true
	//var user models.User
	//user.Username = request.Header.Get("login")
	//user.Password = request.Header.Get("password")
	//user.Active = true
	//note.UserID = h.Service.Repository.CheckUsers(&user)
	UserID, ok := request.Context().Value("id").(int)
	log.Println(UserID)
	if !ok {
		log.Println("userID not found or has the wrong type in context")
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	note.UserID = UserID
	err = h.Service.Create(&note)
	if err != nil {
		log.Println(err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
}
func (h *Handler) Read(response http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println(err)
		response.WriteHeader(http.StatusBadRequest)
		return
	}
	//var user models.User
	//user.Username = request.Header.Get("login")
	//user.Password = request.Header.Get("password")
	//user.Active = true
	//user.ID = h.Service.Repository.CheckUsers(&user)
	userID, ok := request.Context().Value("id").(int)
	if !ok {
		log.Println("userID not found or has the wrong type in context")
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	note, err := h.Service.Read(&id)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	if userID != note.UserID || note.Active == false {
		response.WriteHeader(http.StatusBadRequest)
		return
	}
	jsonData, err := json.Marshal(struct {
		Content string    `json:"content"`
		Date    time.Time `json:"date"`
	}{
		Content: note.Content,
		Date:    note.Date,
	})
	if err != nil {
		log.Println(err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	response.Write(jsonData) //todo сделать структуру ответов
}
func (h *Handler) ReadAll(response http.ResponseWriter, request *http.Request) {
	//notes, err := h.Service.ReadAll()
	//if err != nil {
	//	response.WriteHeader(http.StatusInternalServerError)
	//	return
	//}
	//jsonData, err := json.Marshal(struct {
	//	Content string    `json:"content"`
	//	Date    time.Time `json:"date"`
	//}{
	//	Content: notes[].Content,
	//	Date:    notes[].Date,
	//})
	//response.Write(jsonData)
}
func (h *Handler) Update(response http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println(err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	note, err := h.Service.Read(&id)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	//var user models.User
	//user.Username = request.Header.Get("login")
	//user.Password = request.Header.Get("password")
	//user.Active = true
	//user.ID = h.Service.Repository.CheckUsers(&user)
	userID, ok := request.Context().Value("id").(int)
	if !ok {
		log.Println("userID not found or has the wrong type in context")
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	if userID != note.UserID || note.Active == false {
		response.WriteHeader(http.StatusBadRequest)
		return
	}
	bytes, err := io.ReadAll(request.Body)
	if err != nil {
		log.Println(err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	var updatedNote models.Notes
	err = json.Unmarshal(bytes, &updatedNote)
	if err != nil {
		log.Println(err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	updatedNote.ID = id
	h.Service.Update(&updatedNote)
}
func (h *Handler) Delete(response http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println(err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	note, err := h.Service.Read(&id)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	//var user models.User
	//user.Username = request.Header.Get("login")
	//user.Password = request.Header.Get("password")
	//user.Active = true
	//user.ID = h.Service.Repository.CheckUsers(&user)
	userID, ok := request.Context().Value("id").(int)
	if !ok {
		log.Println("userID not found or has the wrong type in context")
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	if userID != note.UserID || note.Active == false {
		response.WriteHeader(http.StatusBadRequest)
		return
	}
	h.Service.Delete(&id)
}
func (h *Handler) UserRegistration(response http.ResponseWriter, request *http.Request) {
	var user models.User
	bytes, err := io.ReadAll(request.Body)
	if err != nil {
		log.Println(err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(bytes, &user)
	if err != nil {
		log.Println(err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	user.Active = true
	if user.Username == "" {
		response.WriteHeader(http.StatusBadRequest)
	}
	if user.Password == "" {
		response.WriteHeader(http.StatusBadRequest)
	}
	err = h.Service.UserRegistration(&user)
	if err != nil {
		log.Println(err)
		response.WriteHeader(http.StatusBadRequest)
	}
}
