package handlers

import (
	"github.com/cingozr/yemek_sepeti/service"
	"net/http"
)

var (
	memoryStorageService = service.NewMemoryStorage(nil)
)

func SaveKey(w http.ResponseWriter, r *http.Request) {
	if err := memoryStorageService.SaveKey(r.FormValue("key"), r.FormValue("val")); err != nil {
		w.Write([]byte(err.Error()))
	}
}

func GetKey(w http.ResponseWriter, r *http.Request) {
	if res, err := memoryStorageService.GetKey(r.FormValue("key")); err != nil {
		w.Write([]byte(err.Error()))
	} else {
		w.Write([]byte(res))
	}
}

func FlushMemory(w http.ResponseWriter, r *http.Request) {
	memoryStorageService.FlushMemory()
	w.Write([]byte("Memory deleted"))
}
