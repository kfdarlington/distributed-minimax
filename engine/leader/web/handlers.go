package web

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type followerConnections struct {
	Addresses []string `json:"addresses"`
}

func unmarshal(req *http.Request) (followerConnections, error) {
	info := &followerConnections{}
	b, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		return *info, err
	}

	err = json.Unmarshal(b, info)
	return *info, err
}

func (h *handler) index(w http.ResponseWriter, r *http.Request, data interface{}) {
	_, err := fmt.Fprintf(w, "running"); if err != nil {
		h.logger.Errorf("error writing response from / err=%v", err)
	}
}

func (h *handler) ping(w http.ResponseWriter, r *http.Request, data interface{}) {
	_, err := fmt.Fprintf(w, "pong"); if err != nil {
		h.logger.Errorf("error writing response from /ping err=%v", err)
	}
}

func (h *handler) followers(w http.ResponseWriter, r *http.Request, data interface{}) {
	followerConns, err := unmarshal(r)
	if err != nil {
		http.Error(w, err.Error(), 400)
		h.logger.Errorf("error unmarshalling request to /followers err=%v", err)
		return
	}

	confirmations := make([]string, 0)
	errors := make([]error, 0)
	for _, addr := range followerConns.Addresses {
		if err = h.pools.AddFollower(addr); err != nil {
			errors = append(errors, err)
		} else {
			confirmations = append(confirmations, addr)
		}
	}

	if _, err = fmt.Fprintf(w, "successful: %v\nerrors: %v", confirmations, errors); err != nil {
		h.logger.Errorf("error writing response from /followers err=%v", err)
	}
}
