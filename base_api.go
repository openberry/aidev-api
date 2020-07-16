package aidev

import (
	"encoding/json"
	"errors"
	"fmt"
)

const (
	errInvalidResponse = "invalid response"
	errInvalidToken    = "invalid token"
)

type baseAPI struct {
	cl *Client
}

type PersonInput struct {
	Name   string `json:"nombre"`
	CURP   string `json:"curp"`
	Age    string `json:"edad"`
	Gender string `json:"genero"` // 1 male, 2 female
}

type DateInput struct {
	StudyID   string `json:"idEstudio"`
	CabinID   string `json:"idCabina"`
	PersonID  string `json:"idPersona"`
	Date      string `json:"fechaCita"`
	HourStart string `json:"horaInicio"`
	HourEnd   string `json:"horaTermino"`
}

type DateResponse struct {
	CabinID       string `json:"idCabina"`
	PersonID      string `json:"idPersona"`
	AppointmentID string `json:"idCita"`
	StudyID       string `json:"idEstudio"`
	Token         string `json:"token"`
	Name          string `json:"nombre"`
	Code          string `json:"codigo"`
	Age           string `json:"edad"`
	Gender        string `json:"genero"` // 1 male, 2 female
	CURP          string `json:"curp"`
	Date          string `json:"fechaCita"`
	HourStart     string `json:"horaInicio"`
	HourEnd       string `json:"horaFin"`
	Status        string `json:"estatus"` // 0 => Cancelada, 1 => Asignada, 2 => Realizada
}

type CabinResponse struct {
	CabinID     string `json:"idCabina"`
	Code        string `json:"codigo"`
	Description string `json:"descripcion"`
	Address     string `json:"direccion"`
	Lat         string `json:"latitud"`
	Lng         string `json:"longitud"`
	OpenAt      string `json:"horaApertura"`
	CloseAt     string `json:"horaCierre"`
}

type StudyResponse struct {
	StudyID    string `json:"idEstudio"`
	Title      string `json:"titulo"`
	Code       string `json:"codigo"`
	Registered string `json:"fechaRegistro"`
	Updated    string `json:"fechaActulizacion"`
}

type AvailableTimeResponse struct {
	CabinID string   `json:"idCabina"`
	Slots   []string `json:"horarios_disponibles"`
}

type AddAppointmentResponse struct {
	AppointmentID string `json:"idCita"`
	Token         string `json:"token"`
}

type ResultsResponse struct {
	AppointmentID    string `json:"idCita"`
	CabinID          string `json:"idCabina"`
	StudyID          string `json:"idEstudio"`
	PersonID         string `json:"idPersona"`
	Title            string `json:"title"`
	Performed        string `json:"realizado"`
	TotalQuestions   string `json:"totalPreguntas"`
	CorrectQuestions string `json:"respuestaCorrectas"`
	Approved         string `json:"aprobado"` // 1 = yes, 2 = no
	Accuracy         string `json:"veracidad"`
	Result           string `json:"resultado"`
	Tendency         string `json:"tendencia"`
	DataDescription  string `json:"datos"`
	Instructions     string `json:"instrucciones"`
}

func (h *baseAPI) GetToken(user, pass string) (string, error) {
	response, err := h.cl.request(&requestOptions{
		method:   "POST",
		endpoint: "/getToken",
		data: map[string]interface{}{
			"nick": user,
			"psw":  pass,
		},
	})
	if err != nil {
		return "", err
	}
	token, ok := response["token"]
	if !ok {
		return "", errors.New(errInvalidResponse)
	}
	return token.(string), nil
}

func (h *baseAPI) GetDates(date string) ([]DateResponse, error) {
	response, err := h.cl.request(&requestOptions{
		method:   "POST",
		endpoint: "/getDates",
		data: map[string]interface{}{
			"token": h.cl.token,
			"dia":   date,
		},
	})
	if err != nil {
		return nil, err
	}
	var data []DateResponse
	if err := parseResponse(response, &data); err != nil {
		return nil, err
	}
	return data, nil
}

func (h *baseAPI) GetCabins() ([]CabinResponse, error) {
	response, err := h.cl.request(&requestOptions{
		method:   "POST",
		endpoint: "/getCabins",
		data: map[string]interface{}{
			"token": h.cl.token,
		},
	})
	if err != nil {
		return nil, err
	}
	var data []CabinResponse
	if err := parseResponse(response, &data); err != nil {
		return nil, err
	}
	return data, nil
}

func (h *baseAPI) GetStudies() ([]StudyResponse, error) {
	response, err := h.cl.request(&requestOptions{
		method:   "POST",
		endpoint: "/getStudies",
		data: map[string]interface{}{
			"token": h.cl.token,
		},
	})
	if err != nil {
		return nil, err
	}
	var data []StudyResponse
	if err := parseResponse(response, &data); err != nil {
		return nil, err
	}
	return data, nil
}

func (h *baseAPI) GetAvailableTime(cabinID string, date string) ([]string, error) {
	response, err := h.cl.request(&requestOptions{
		method:   "POST",
		endpoint: "/getAvailableTimeCabin",
		data: map[string]interface{}{
			"token":    h.cl.token,
			"idCabina": cabinID,
			"dia":      date,
		},
	})
	if err != nil {
		return nil, err
	}
	var data []string
	if err := parseResponse(response, &data); err != nil {
		return nil, err
	}
	return data, nil
}

func (h *baseAPI) AddPerson(pr PersonInput) (string, error) {
	response, err := h.cl.request(&requestOptions{
		method:   "POST",
		endpoint: "/addPerson",
		data:     getData(pr, h.cl.token),
	})
	if err != nil {
		return "", err
	}
	id, ok := response["idPersona"]
	if !ok {
		return "", errors.New(errInvalidResponse)
	}
	return fmt.Sprintf("%v", id), nil
}

func (h *baseAPI) AddAppointment(dr DateInput) (AddAppointmentResponse, error) {
	data := AddAppointmentResponse{}
	response, err := h.cl.request(&requestOptions{
		method:   "POST",
		endpoint: "/addDate",
		data:     getData(dr, h.cl.token),
	})
	if err != nil {
		return data, err
	}
	if err := parseResponse(response, &data); err != nil {
		return data, err
	}
	return data, nil
}

func (h *baseAPI) GetResults(appointmentID int) (ResultsResponse, error) {
	data := ResultsResponse{}
	response, err := h.cl.request(&requestOptions{
		method:   "POST",
		endpoint: "/getResults",
		data: map[string]interface{}{
			"token":  h.cl.token,
			"idCita": fmt.Sprintf("%d", appointmentID),
		},
	})
	if err != nil {
		return data, err
	}
	if err := parseResponse(response, &data); err != nil {
		return data, err
	}
	return data, nil
}

func getData(v interface{}, token string) map[string]interface{} {
	js, _ := json.Marshal(v)
	data := make(map[string]interface{})
	_ = json.Unmarshal(js, &data)
	data["token"] = token
	return data
}

func parseResponse(input map[string]interface{}, output interface{}) error {
	data, ok := input["data"]
	if !ok {
		return errors.New(errInvalidResponse)
	}
	js, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return json.Unmarshal(js, output)
}
