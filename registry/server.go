package registry

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

const (
	ServerPort = ":50012"
	ServiceURL = "http://localhost" + ServerPort + "/services"
)

type registry struct {
	registrations []Registration
	mux           *sync.Mutex
}

func (r *registry) add(reg Registration) error {
	r.mux.Lock()
	r.registrations = append(r.registrations, reg)
	r.mux.Unlock()
	return nil
}

func (r *registry) remove(serviceUrl string) error {

	for i := range r.registrations {
		if serviceUrl == r.registrations[i].ServiceUrl {
			r.mux.Lock()
			r.registrations = append(r.registrations[:i], r.registrations[i+1:]...)
			r.mux.Unlock()
			return nil
		}
	}

	return fmt.Errorf("Service at URL %v not found\n", serviceUrl)
}

var reg = registry{
	registrations: make([]Registration, 0),
	mux:           new(sync.Mutex),
}

type RegistryService struct{}

func (rs RegistryService) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	log.Println("Request received")

	switch req.Method {

	case http.MethodPost:

		var newRegistration Registration
		err := json.NewDecoder(req.Body).Decode(&newRegistration)
		if err != nil {
			log.Println(err)
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		log.Printf("Adding service: %v with url: %v\n", newRegistration.ServiceName, newRegistration.ServiceUrl)
		err = reg.add(newRegistration)
		if err != nil {
			log.Println(err)
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

	case http.MethodDelete:
		payload, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Println(err)
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		url := string(payload)
		log.Printf("Deleting service at url: %v\n", url)
		err = reg.remove(url)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
		}

	default:
		rw.WriteHeader(http.StatusMethodNotAllowed)
	}
}
