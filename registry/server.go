package registry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

const ServerPort = ":3000"
const ServiceURL = "http://localhost" + ServerPort + "/services"

type registry struct{
	registrations []Registration	
	mutex *sync.RWMutex
}

func (r *registry) add(reg Registration) error{
	r.mutex.Lock()
	r.registrations = append(r.registrations, reg)
	r.mutex.Unlock()
	err := r.sendRequiredServices(reg)
	r.notify(patch{
		Added: []patchEntry{
			patchEntry{
				Name: reg.ServiceName,
				URL: reg.ServiceURL,
			},
		},
	})	
	return err
}

func (r registry) notify(fullPath patch) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	for _, reg := range r.registrations {
		go func(reg Registration) {
			for _, reqservice := range reg.RequiredServices {
				p := patch{Added: []patchEntry{}, Removed: []patchEntry{}}
				sendUpdate := false
				for _, added := range fullPath.Added {
					if added.Name == reqservice {
						p.Added = append(p.Added, added)
						sendUpdate = true
					}
				}
				for _, removed := range fullPath.Removed {
					if removed.Name == reqservice {
						p.Removed = append(p.Removed, removed)
						sendUpdate = true
					}
				}
				if sendUpdate {
					err := r.sendPatch(p, reg.ServiceUpdateURL)
					if err != nil {
						log.Panicln(err)
						return
					}
				}
			}
		}(reg)
	}
}

func (r registry) sendRequiredServices(reg Registration) error{
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var p patch
	for _, serviceReg := range r.registrations {
		for _, reqService := range reg.RequiredServices {
			if serviceReg.ServiceName == reqService {
				p.Added = append(p.Added, patchEntry{
					Name: serviceReg.ServiceName,
					URL: serviceReg.ServiceURL,
				})
			}
		}
	}
	err := r.sendPatch(p, reg.ServiceUpdateURL)
	if err != nil {
		return err
	}
	return nil
}

func (r registry) sendPatch(p patch, url string) error {
	d, err := json.Marshal(p)
	if err != nil {
		return err
	}
	_, err = http.Post(url, "application/json",	 bytes.NewBuffer(d))
	if  err != nil {
		return err
	}
	return nil
}

func (r * registry) remove(url string) error {
	for i := range reg.registrations {
		if reg.registrations[i].ServiceURL == url {
			r.notify(patch{
				Removed: []patchEntry{
					patchEntry{
						Name: r.registrations[i].ServiceName,
						URL: r.registrations[i].ServiceURL,
					},
				},
			})
			r.mutex.Lock()
			reg.registrations = append(reg.registrations[:i], reg.registrations[i+1:]...)
			r.mutex.Unlock()
			return nil
		}
	}
	return fmt.Errorf("service at URL %s not found", url)
}

func (r *registry) heartbeat(freq time.Duration) {
	for {
		var wg sync.WaitGroup
		for _, reg := range r.registrations {
			wg.Add(1)
			go func(reg Registration) {
				defer wg.Done()
				success := true
				for attemps := 0; attemps < 3; attemps++ {
					res,err := http.Get(reg.HeartBeatURL)
					if err != nil {
						log.Println(err)
					} else if res.StatusCode == http.StatusOK {
						log.Printf("HeartBeat check passed for %v", reg.ServiceName)
						if !success {
							r.add(reg)
						}
						break;
					}
					log.Printf("HeartBeat check faild for %v", reg.ServiceName)
					if success {
						success = false
						r.remove(reg.ServiceURL)
					}
					time.Sleep(1 * time.Second)
				}
			}(reg)
			wg.Wait()
			time.Sleep(freq)
		}
	}
}

var once sync.Once

func SetupRegistryService() {
	once.Do(func() {
		go reg.heartbeat(3 * time.Second)
	})
}

var reg = registry{
	registrations: make([]Registration, 0),
	mutex: new(sync.RWMutex),
}

type RegistryService struct{}	

func (s RegistryService) ServeHTTP(w http.ResponseWriter, r *http.Request){
	log.Println	("Request received")
	switch r.Method {
	case http.MethodPost:
		dec := json.NewDecoder(r.Body)
		var r Registration
		err := dec.Decode(&r)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		log.Printf("Adding service: %v with URL: %s\n", r.ServiceName, r.ServiceURL)
		err = reg.add(r)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	case http.MethodDelete:
		payload, err := io.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		url := string(payload)
		log.Printf("Removing service at URL: %s\n", url)
		err = reg.remove(url)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}