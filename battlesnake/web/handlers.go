package web

import (
	"encoding/json"
	"fmt"
	"github.com/kristian-d/distributed-minimax/battlesnake/game"
	"io/ioutil"
	"net/http"
)

func unmarshal(req *http.Request) (game.Update, error) {
	info := &game.Update{}
	b, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		return *info, err
	}

	err = json.Unmarshal(b, info)
	return *info, err
}

func (h *handler) index(w http.ResponseWriter, r *http.Request, data interface{}) {
	_, err := fmt.Fprintf(w, "sssSssSsssssSs I am alive! sssSssSsssssSs"); if err != nil {
		h.logger.Errorf("error writing response from / err=%v", err)
	}
}

func (h *handler) ping(w http.ResponseWriter, r *http.Request, data interface{}) {
	_, err := fmt.Fprintf(w, "sssSssSsssssSs I am alive! sssSssSsssssSs"); if err != nil {
		h.logger.Errorf("error writing response from /ping err=%v", err)
	}
}

func (h *handler) start(w http.ResponseWriter, r *http.Request, data interface{}) {
	state, err := unmarshal(r)
	if err != nil {
		http.Error(w, err.Error(), 400)
		h.logger.Errorf("error unmarshalling request to /start err=%v", err)
		return
	}

	game.CreateGame(state)
	fmt.Printf("game created id=%s\n", state.Game.Id)

	res, err := json.Marshal(struct {
		Color    string `json:"color"`
		HeadType string `json:"headType"`
		TailType string `json:"tailType"`
	}{
		Color:    "#ff00ff",
		HeadType: "bendr",
		TailType: "pixel",
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		h.logger.Errorf("error marshalling response err=%v", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(res)
	if err != nil {
		h.logger.Errorf("error writing response from /start err=%v", err)
	}
}

func (h *handler) move(w http.ResponseWriter, r *http.Request, data interface{}) {
	state, err := unmarshal(r)
	if err != nil {
		http.Error(w, err.Error(), 400)
		h.logger.Errorf("error unmarshalling request to /move err=%v", err)
		return
	}

	err = game.UpdateGame(state)
	if err != nil {
		http.Error(w, err.Error(), 500)
		h.logger.Infof("cannot update game id=%s\n", state.Game.Id)
		return
	}
	h.logger.Infof("game updated id=%s\n", state.Game.Id)

	result := h.engine.ComputeMove(game.Games[state.Game.Id].Board, 1000) // process the move for x ms, leaving (500 - x) ms for the network
	res, err := json.Marshal(struct {
		Move  string `json:"move"`
		Shout string `json:"shout"`
	}{
		Move:  string(result),
		Shout: "shouting!",
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(res)
	if err != nil {
		h.logger.Errorf("error writing response from /move err=%v", err)
	}
}

func (h *handler) end(w http.ResponseWriter, r *http.Request, data interface{}) {
	state, err := unmarshal(r)
	if err != nil {
		http.Error(w, err.Error(), 400)
		h.logger.Errorf("error unmarshalling request to /end err=%v", err)
		return
	}

	err = game.DeleteGame(state); if err != nil {
		http.Error(w, err.Error(), 404)
		h.logger.Infof("error deleting game id=%s\n", state.Game.Id)
		return
	}
	h.logger.Infof("game deleted id=%s\n", state.Game.Id)
}
