package http

import (
	logger "github.com/sirupsen/logrus"
	"github.com/ugorji/go/codec"
	"net/http"
)

var jsonHandle codec.JsonHandle

func EncodeResponse(w http.ResponseWriter, statusCode int, res interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := codec.NewEncoder(w, &jsonHandle).Encode(res); err != nil {
		logger.Error(err)
	}
}

func EncodeByteResponse(w http.ResponseWriter, statusCode int, res []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, err := w.Write(res)
	if err != nil {
		logger.Error(err)
	}
}

func JsonToStruct(r *http.Request, target interface{}) error {
	err := codec.NewDecoder(r.Body, &jsonHandle).Decode(target)
	if err != nil {
		return err
	}

	return nil
}
