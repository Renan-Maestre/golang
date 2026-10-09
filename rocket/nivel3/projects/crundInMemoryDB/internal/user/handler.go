package user

import (
	"encoding/json"
	"net/http"

	"crudMemoryDB/internal/model"
	"crudMemoryDB/internal/shared/response"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) FindAll(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.FindAll()
	if err != nil {
		response.JSON(
			w,
			response.ApiResponse{
				Error: "user not found",
			},
			http.StatusNotFound,
		)
		return
	}

	response.JSON(w, response.ApiResponse{Data: users}, http.StatusOK)
}

func (h *Handler) FindByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "idUser")
	user, err := h.service.FindByID(id)
	if err != nil {
		response.JSON(
			w,
			response.ApiResponse{
				Error: "user not found",
			},
			http.StatusNotFound,
		)
		return
	}

	response.JSON(w, response.ApiResponse{Data: user}, http.StatusOK)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var body model.User

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.JSON(w, response.ApiResponse{Error: "invalid body"}, http.StatusUnprocessableEntity)
		return
	}

	if err := h.service.Create(body); err != nil {
		response.JSON(w, response.ApiResponse{Error: err.Error()}, http.StatusBadRequest)
	}

	response.JSON(
		w,
		response.ApiResponse{
			Data: "Criado com Sucesso",
		},
		http.StatusCreated,
	)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	var body model.User

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.JSON(w, response.ApiResponse{Error: "invalid body"}, http.StatusUnprocessableEntity)
		return
	}

	if err := h.service.Update(body); err != nil {
		response.JSON(w, response.ApiResponse{Error: err.Error()}, http.StatusBadRequest)
	}

	response.JSON(
		w,
		response.ApiResponse{
			Data: "Atualizado com Sucesso",
		},
		http.StatusCreated,
	)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "idUser")
	if err := h.service.Delete(id); err != nil {
		response.JSON(
			w,
			response.ApiResponse{
				Error: "user not found",
			},
			http.StatusNotFound,
		)
		return
	}

	response.JSON(w, response.ApiResponse{Data: "Deletado com Sucesso"}, http.StatusOK)
}
