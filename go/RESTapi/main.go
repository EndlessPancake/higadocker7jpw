package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/ant0ine/go-json-rest/rest"
)

func main() {
	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	router, err := rest.MakeRouter(
		rest.Get("/services", GetAllCountries),
		rest.Post("/services", PostService),
		rest.Get("/services/:code", GetService),
		rest.Delete("/services/:code", DeleteService),
	)
	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(router)
	log.Fatal(http.ListenAndServe(":8080", api.MakeHandler()))
}

type Service struct {
	ServiceCode string
	IPv4        string
}

var store = map[string]*Service{}

var lock = sync.RWMutex{}

func GetService(w rest.ResponseWriter, r *rest.Request) {
	code := r.PathParam("code")

	lock.RLock()
	var service *Service
	if store[code] != nil {
		service = &Service{}
		*service = *store[code]
	}
	lock.RUnlock()

	if service == nil {
		rest.NotFound(w, r)
		return
	}
	w.WriteJson(service)
}

func GetAllCountries(w rest.ResponseWriter, r *rest.Request) {
	lock.RLock()
	services := make([]Service, len(store))
	i := 0
	for _, service := range store {
		services[i] = *service
		i++
	}
	lock.RUnlock()
	w.WriteJson(&services)
}

func PostService(w rest.ResponseWriter, r *rest.Request) {
	service := Service{}
	err := r.DecodeJsonPayload(&service)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if service.ServiceCode == "" {
		rest.Error(w, "Service code required", 400)
		return
	}
	if service.IPv4 == "" {
		rest.Error(w, "Service name required", 400)
		return
	}
	lock.Lock()
	store[service.ServiceCode] = &service
	lock.Unlock()
	w.WriteJson(&service)
}

func DeleteService(w rest.ResponseWriter, r *rest.Request) {
	code := r.PathParam("code")
	lock.Lock()
	delete(store, code)
	lock.Unlock()
	w.WriteHeader(http.StatusOK)
}
