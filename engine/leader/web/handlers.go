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

func (h *handler) postFollowers(w http.ResponseWriter, r *http.Request, data interface{}) {
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
			h.logger.Errorf("error when attempting connection at /followers err=%v", err)
		} else {
			confirmations = append(confirmations, addr)
		}
	}

	res, err := json.Marshal(struct {
		Successful  []string `json:"successful"`
		Errors []error `json:"errors"`
	}{
		Successful: confirmations,
		Errors: errors,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(res)
	if err != nil {
		h.logger.Errorf("error writing response to POST from /follower err=%v", err)
	}
}

func (h *handler) getFollowers(w http.ResponseWriter, r *http.Request, data interface{}) {
	addresses := h.pools.GetFollowerAddresses()
	res, err := json.Marshal(struct {
		Addresses []string `json:"addresses"`
	}{
		Addresses: addresses,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(res)
	if err != nil {
		h.logger.Errorf("error writing response to GET from /follower err=%v", err)
	}
}
