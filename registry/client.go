package registry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

type providers struct {
	services map[ServiceName][]string
	mux      *sync.RWMutex
}

var provs = providers{
	services: make(map[ServiceName][]string),
	mux:      new(sync.RWMutex),
}

func (ps *providers) Update(p patch) {
	ps.mux.Lock()
	defer ps.mux.Unlock()

	for _, entry := range p.Added {
		_, ok := ps.services[entry.ServiceName]
		if !ok {
			ps.services[entry.ServiceName] = make([]string, 0)
		}

		ps.services[entry.ServiceName] = append(ps.services[entry.ServiceName], entry.URL)
	}

	for _, entry := range p.Removed {
		providerURLs, ok := ps.services[entry.ServiceName]
		if ok {
			for i := range providerURLs {
				if providerURLs[i] == entry.URL {
					ps.services[entry.ServiceName] = append(providerURLs[:i], providerURLs[i+1:]...)
				}
			}
		}
	}
}

func (ps *providers) get(name ServiceName) (string, error) {
	prov, ok := ps.services[name]
	if !ok {
		return "", fmt.Errorf("No providers available for service %v", name)
	}

	idx := int(rand.Float32() * float32(len(prov)))
	return prov[idx], nil
}

func GetProvider(name ServiceName) (string, error) {
	return provs.get(name)
}

func RegisterService(reg Registration) error {

	hearbeatURL, err := url.Parse(reg.HearbeatURL)
	if err != nil {
		log.Println(err)
	}

	http.HandleFunc(hearbeatURL.Path, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})


	serviceUpdateURL, err := url.Parse(reg.UpdateServiceURL)
	if err != nil {
		return err
	}

	http.Handle(serviceUpdateURL.Path, &serviceUpdateHandler{})

	buf := new(bytes.Buffer)
	err = json.NewEncoder(buf).Encode(reg)
	if err != nil {
		return err
	}

	res, err := http.Post(ServiceURL, "application/json", buf)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("Failed to register service. Register service respond with code %v", res.StatusCode)
	}

	return nil
}

func ShutdownService(serviceURL string) error {
	req, err := http.NewRequest(http.MethodDelete, ServiceURL, strings.NewReader(serviceURL))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "plain/text")
	res, err := http.DefaultClient.Do(req)
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to deregister service. Register service respond with cdoe %v", res.StatusCode)
	}

	return err
}

type serviceUpdateHandler struct{}

func (h serviceUpdateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "hasn't implement method: %v for url: %v\n", r.Method, r.URL)
		return
	}

	var p patch
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		return
	}

	fmt.Printf("Updated received: %+v\n", p)
	provs.Update(p)
}
