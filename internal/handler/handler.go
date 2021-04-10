package handler

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/xeynse/XeynseJar_analytics/internal/entity"
	"github.com/xeynse/XeynseJar_analytics/internal/usecase/jar"
	errUtil "github.com/xeynse/XeynseJar_analytics/internal/util/error"
)

type handler struct {
	jarUseCase jar.Resource
}

func New(router *httprouter.Router, jarUseCase jar.Resource) {
	handler := &handler{
		jarUseCase: jarUseCase,
	}
	router.GET("/xeynseJar/analytics/stats/home/:homeID/jars", httpHandler(handler.fetchAllJarStats))
	router.GET("/xeynseJar/analytics/stats/home/:homeID/jars/:jarID", httpHandler(handler.fetchJarStatsByJarID))
}

func httpHandler(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		h(w, r, ps)
	}
}

func (h handler) fetchAllJarStats(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	homeID := ps.ByName("homeID")
	if homeID == "" {
		utilErr := errUtil.BadError("", "home id can not be empty")
		response := entity.Response{Status: utilErr.HttpCode, Success: false, Message: "Fail", Error: utilErr}
		responseBytes, _ := json.Marshal(response)
		w.WriteHeader(utilErr.HttpCode)
		w.Write(responseBytes)
		return
	}

	jarStatusResponse, err := h.jarUseCase.GetAllJarStats(homeID)
	if err != nil {
		utilErr := errUtil.InternalServerError("", "error occurred in fetching home analytics: "+err.Error())
		response := entity.Response{Status: utilErr.HttpCode, Success: false, Message: "Fail", Error: utilErr}
		responseBytes, _ := json.Marshal(response)
		w.WriteHeader(utilErr.HttpCode)
		w.Write(responseBytes)
		return
	}

	response := entity.Response{
		Status:  200,
		Success: true,
		Message: "Success",
		Error:   nil,
		Data:    jarStatusResponse,
	}

	responseBytes, err := json.Marshal(response)
	if err != nil {
		utilErr := errUtil.InternalServerError("", "error marshalling response")
		response := entity.Response{Status: utilErr.HttpCode, Success: false, Message: "Fail", Error: utilErr}
		responseBytes, _ := json.Marshal(response)
		w.WriteHeader(utilErr.HttpCode)
		w.Write(responseBytes)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(responseBytes)
	return
}

func (h handler) fetchJarStatsByJarID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	homeID := ps.ByName("homeID")
	jarID := ps.ByName("jarID")

	if homeID == "" {
		utilErr := errUtil.BadError("", "home id can not be empty")
		response := entity.Response{Status: utilErr.HttpCode, Success: false, Message: "Fail", Error: utilErr}
		responseBytes, _ := json.Marshal(response)
		w.WriteHeader(utilErr.HttpCode)
		w.Write(responseBytes)
		return
	}

	if jarID == "" {
		utilErr := errUtil.BadError("", "jar id can not be empty")
		response := entity.Response{Status: utilErr.HttpCode, Success: false, Message: "Fail", Error: utilErr}
		responseBytes, _ := json.Marshal(response)
		w.WriteHeader(utilErr.HttpCode)
		w.Write(responseBytes)
		return
	}

	jarStatusResponse, err := h.jarUseCase.GetJarStatByJarID(homeID, jarID)
	if err != nil {
		utilErr := errUtil.InternalServerError("", "error occurred in fetching home analytics: "+err.Error())
		response := entity.Response{Status: utilErr.HttpCode, Success: false, Message: "Fail", Error: utilErr}
		responseBytes, _ := json.Marshal(response)
		w.WriteHeader(utilErr.HttpCode)
		w.Write(responseBytes)
		return
	}

	response := entity.Response{
		Status:  200,
		Success: true,
		Message: "Success",
		Error:   nil,
		Data:    jarStatusResponse,
	}

	responseBytes, err := json.Marshal(response)
	if err != nil {
		utilErr := errUtil.InternalServerError("", "error marshalling response")
		response := entity.Response{Status: utilErr.HttpCode, Success: false, Message: "Fail", Error: utilErr}
		responseBytes, _ := json.Marshal(response)
		w.WriteHeader(utilErr.HttpCode)
		w.Write(responseBytes)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(responseBytes)
	return
}
