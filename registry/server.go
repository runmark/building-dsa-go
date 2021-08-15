package registry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

const (
	ServerPort = ":50012"
	ServiceURL = "http://localhost" + ServerPort + "/services"
)

type registry struct {
	registrations []Registration
	mux           *sync.RWMutex
}

func (r *registry) add(reg Registration) error {
	r.mux.Lock()
	r.registrations = append(r.registrations, reg)
	r.mux.Unlock()
	err := r.sendRequiredServices(reg)
	r.notify(patch{
		Added: []patchEntry{
			{reg.ServiceName, reg.ServiceUrl},
		},
	})
	return err
}

func (r *registry) notify(fullPatch patch) {

	r.mux.RLock()
	defer r.mux.RUnlock()

	for _, reg := range r.registrations {

		go func(reg Registration) {

			p := patch{Added: []patchEntry{}, Removed: []patchEntry{}}

			for _, requiredServiceName := range reg.RequiredServices {
				for _, addEntry := range fullPatch.Added {
					if requiredServiceName == addEntry.ServiceName {
						p.Added = append(p.Added, addEntry)
					}
				}

				for _, removedEntry := range fullPatch.Removed {
					if requiredServiceName == removedEntry.ServiceName {
						p.Removed = append(p.Removed, removedEntry)
					}
				}
			}

			if len(p.Added) > 0 || len(p.Removed) > 0 {
				err := r.sendPatch(p, reg.UpdateServiceURL)
				if err != nil {
					log.Println(err)
					return
				}
			}
		}(reg)

	}

}

func (r registry) sendRequiredServices(reg Registration) error {
	r.mux.RLock()
	defer r.mux.RUnlock()

	var p patch
	for _, reqService := range reg.RequiredServices {
		for _, reg := range r.registrations {
			if reqService == reg.ServiceName {
				pe := patchEntry{ServiceName: reg.ServiceName, URL: reg.ServiceUrl}
				p.Added = append(p.Added, pe)
			}
		}
	}

	err := r.sendPatch(p, reg.UpdateServiceURL)

	return err
}

func (r registry) sendPatch(p patch, url string) error {
	data := new(bytes.Buffer)
	json.NewEncoder(data).Encode(p)

	_, err := http.Post(url, "application/json", data)
	return err
}

func (r *registry) remove(serviceUrl string) error {

	for i := range r.registrations {
		if serviceUrl == r.registrations[i].ServiceUrl {

			r.notify(patch{
				Removed: []patchEntry{
					{r.registrations[i].ServiceName, serviceUrl},
				},
			})

			r.mux.Lock()
			r.registrations = append(r.registrations[:i], r.registrations[i+1:]...)
			r.mux.Unlock()

			return nil
		}
	}

	return fmt.Errorf("Service at URL %v not found\n", serviceUrl)
}

var once sync.Once

func SetupRegistryService() {
	once.Do(func() {
		go reg.heartbeat(time.Second * 3)
	})
}

func (r *registry) heartbeat(freq time.Duration) {
	for {
		var waitGroup sync.WaitGroup
		for _, reg := range r.registrations {

			waitGroup.Add(1)
			go func(reg Registration) {

				defer waitGroup.Done()
				lastSuccess := true

				for i := 0; i < 3; i++ {

					resp, err := http.Get(fmt.Sprintf("%v/heartbeat", reg.ServiceUrl))
					if err != nil {
						log.Println(err)
					} else if resp.StatusCode == http.StatusOK {
						log.Printf("heartbeat check passed for %v\n", reg.ServiceName)
						if !lastSuccess {
							r.registrations = append(r.registrations, reg)
						}
						break
					}
					log.Printf("heartbeat check failed for %v\n", reg.ServiceName)
					if lastSuccess {
						lastSuccess = false
						r.remove(reg.ServiceUrl)
					}

					time.Sleep(3 * time.Second) // wait to try again, can Progressively try
				}
			}(reg)

		}

		waitGroup.Wait()
		time.Sleep(freq)
	}
}

var reg = registry{
	registrations: make([]Registration, 0),
	mux:           new(sync.RWMutex),
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
