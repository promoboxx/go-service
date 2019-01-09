package service

import (
	"encoding/json"
	"net/http"

	"github.com/Sirupsen/logrus"

	"github.com/promoboxx/go-glitch/glitch"
)

// ReturnProblem will return a json http problem response
func ReturnProblem(w http.ResponseWriter, detail, code string, status int) (int, []byte) {
	prob := glitch.HTTPProblem{
		Title:  http.StatusText(status),
		Detail: detail,
		Code:   code,
		Status: status,
	}
	by, _ := json.Marshal(prob)
	if w != nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
	}
	return status, by
}

// LogErrorAndReturnProblem logs the given error and returns a json http problem response
func LogErrorAndReturnProblem(w http.ResponseWriter, detail string, code string, status int, err error, log *logrus.Entry) (int, []byte) {
	if status >= 500 {
		log.Error(err)
	} else {
		log.Warn(err)
	}

	return ReturnProblem(w, detail, code, status)
}

// WriteProblem will write a json http problem response
func WriteProblem(w http.ResponseWriter, detail, code string, status int) error {
	prob := glitch.HTTPProblem{
		Title:  http.StatusText(status),
		Detail: detail,
		Code:   code,
		Status: status,
	}
	by, err := json.Marshal(prob)
	if err != nil {
		return err
	}
	if w != nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(status)
		_, err = w.Write(by)
	}
	return err
}

// WriteJSONResponse will write a json response to the htt.ResponseWriter
func WriteJSONResponse(w http.ResponseWriter, status int, data interface{}) error {
	var by []byte
	var err error
	if data != nil {
		by, err = json.Marshal(data)
		if err != nil {
			return err
		}
	}
	if w != nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(status)
		_, err = w.Write(by)
	}
	return err
}
