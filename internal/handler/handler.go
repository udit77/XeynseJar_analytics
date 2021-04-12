package handler

import (
	"encoding/json"
	"io/ioutil"
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

	router.POST("/xeynseJar/analytics/stats/home", httpHandler(handler.fetchAllJarStats))
	router.POST("/xeynseJar/analytics/stats/home/jars", httpHandler(handler.fetchJarStatsByJarID))

	router.POST("/xeynseJar/analytics/consumption/calories", httpHandler(handler.getCaloriesConsumptiion))
}

func httpHandler(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		h(w, r, ps)
	}
}

func (h handler) fetchAllJarStats(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utilErr := errUtil.BadError("", "error reading create mode request body: "+err.Error())
		response := entity.Response{Status: utilErr.HttpCode, Success: false, Message: "Fail", Error: utilErr}
		responseBytes, _ := json.Marshal(response)
		w.WriteHeader(utilErr.HttpCode)
		w.Write(responseBytes)
		return

	}
	statusRequest := &entity.AllJarStatusRequest{}
	err = json.Unmarshal(requestBody, &statusRequest)
	if err != nil {
		utilErr := errUtil.BadError("", "error unmarshalling create mode request body")
		response := entity.Response{Status: utilErr.HttpCode, Success: false, Message: "Fail", Error: utilErr}
		responseBytes, _ := json.Marshal(response)
		w.WriteHeader(utilErr.HttpCode)
		w.Write(responseBytes)
		return
	}

	if statusRequest.HomeID == "" {
		utilErr := errUtil.BadError("", "home ID cannot be empty")
		response := entity.Response{Status: utilErr.HttpCode, Success: false, Message: "Fail", Error: utilErr}
		responseBytes, _ := json.Marshal(response)
		w.WriteHeader(utilErr.HttpCode)
		w.Write(responseBytes)
		return
	}

	jarStatusResponse, err := h.jarUseCase.GetAllJarStats(statusRequest.HomeID)
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
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utilErr := errUtil.BadError("", "error reading create mode request body: "+err.Error())
		response := entity.Response{Status: utilErr.HttpCode, Success: false, Message: "Fail", Error: utilErr}
		responseBytes, _ := json.Marshal(response)
		w.WriteHeader(utilErr.HttpCode)
		w.Write(responseBytes)
		return

	}
	statusRequest := &entity.JarStatusRequest{}
	err = json.Unmarshal(requestBody, &statusRequest)
	if err != nil {
		utilErr := errUtil.BadError("", "error unmarshalling create mode request body")
		response := entity.Response{Status: utilErr.HttpCode, Success: false, Message: "Fail", Error: utilErr}
		responseBytes, _ := json.Marshal(response)
		w.WriteHeader(utilErr.HttpCode)
		w.Write(responseBytes)
		return
	}

	if statusRequest.HomeID == "" {
		utilErr := errUtil.BadError("", "home ID cannot be empty")
		response := entity.Response{Status: utilErr.HttpCode, Success: false, Message: "Fail", Error: utilErr}
		responseBytes, _ := json.Marshal(response)
		w.WriteHeader(utilErr.HttpCode)
		w.Write(responseBytes)
		return
	}

	if statusRequest.JarID == "" {
		utilErr := errUtil.BadError("", "jar ID cannot be empty")
		response := entity.Response{Status: utilErr.HttpCode, Success: false, Message: "Fail", Error: utilErr}
		responseBytes, _ := json.Marshal(response)
		w.WriteHeader(utilErr.HttpCode)
		w.Write(responseBytes)
		return
	}

	jarStatusResponse, err := h.jarUseCase.GetJarStatByJarID(statusRequest)
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

func (h handler) getCaloriesConsumptiion(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utilErr := errUtil.BadError("", "error reading create mode request body: "+err.Error())
		response := entity.Response{Status: utilErr.HttpCode, Success: false, Message: "Fail", Error: utilErr}
		responseBytes, _ := json.Marshal(response)
		w.WriteHeader(utilErr.HttpCode)
		w.Write(responseBytes)
		return

	}
	statusRequest := &entity.ConsumptionRequest{}
	err = json.Unmarshal(requestBody, &statusRequest)
	if err != nil {
		utilErr := errUtil.BadError("", "error unmarshalling create mode request body")
		response := entity.Response{Status: utilErr.HttpCode, Success: false, Message: "Fail", Error: utilErr}
		responseBytes, _ := json.Marshal(response)
		w.WriteHeader(utilErr.HttpCode)
		w.Write(responseBytes)
		return
	}

	if statusRequest.HomeID == "" {
		utilErr := errUtil.BadError("", "home ID cannot be empty")
		response := entity.Response{Status: utilErr.HttpCode, Success: false, Message: "Fail", Error: utilErr}
		responseBytes, _ := json.Marshal(response)
		w.WriteHeader(utilErr.HttpCode)
		w.Write(responseBytes)
		return
	}

	if statusRequest.JarID == "" {
		utilErr := errUtil.BadError("", "jar ID cannot be empty")
		response := entity.Response{Status: utilErr.HttpCode, Success: false, Message: "Fail", Error: utilErr}
		responseBytes, _ := json.Marshal(response)
		w.WriteHeader(utilErr.HttpCode)
		w.Write(responseBytes)
		return
	}

	jarStatusResponse, err := h.jarUseCase.GetJarCalorieConsumption(statusRequest)
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
