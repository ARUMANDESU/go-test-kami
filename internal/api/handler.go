package api

import "net/http"

func (a API) ReserveRoom(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("reserve room"))
}

func (a API) GetRoomReservations(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("get room reservations"))
}
