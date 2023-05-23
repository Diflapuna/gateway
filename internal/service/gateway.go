package service

import (
	"encoding/json"
	"gateway/internal/models"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

//const (
//	baseURL = "/api/v1/"
//)

type Gateway struct {
	Logger          *zap.SugaredLogger
	Router          *mux.Router
	ProtectedRouter *mux.Router
}

func NewGateway() *Gateway {
	gw := &Gateway{Router: mux.NewRouter()}
	gw.Logger = NewLogger()
	gw.registerHandlers()

	return gw
}

func (gw *Gateway) Start() error {
	gw.Logger.Info("Started gateway: on port 1337")
	err := http.ListenAndServe(":1337", gw.Router)
	if err != nil {
		gw.Logger.Fatal("Can't start gateway: ", err)
		return err
	}

	return nil
}

func (gw *Gateway) hello() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		greet := &models.Greeting{Greeting: "O hi Mark!"}
		if err := json.NewEncoder(w).Encode(greet); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		//w.WriteHeader(http.StatusOK)  superfluous response.WriteHeader call
	}
}

func (gw *Gateway) gethandlersList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		service := &models.Service{}
		if err := json.NewDecoder(r.Body).Decode(service); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			gw.Logger.Errorf("Failed to decode request: %w", err)

			return
		}
		defer r.Body.Close()

		registerService(service, gw, gw.registerDefaultService)
	}
}

func (gw *Gateway) registerHandlers() {
	gw.Router.Path("/hello").Handler(gw.hello()).Methods("GET")
	gw.Router.Path("/handlers").Handler(gw.gethandlersList()).Methods("POST")
	//gw.registerNewHandlers()
}

//Закоментил registerNewHandlers потому что теперь это работает без нее, но я не знаю будем ли мы
//возвращаться к похожей модели, если это долго тут пролежит без дела то можно будет смело сносить
/*
func (gw *Gateway) registerNewHandlers() {
	list := make(map[string]func(string, string) http.HandlerFunc)
	list["user"] = gw.registerDefaultService

	// Антоха и Леха это мок на список сервисов убогий который надо будет переделать мне сейчас просто лень пиздец
	//задачу обьясню попозже просто хочю понять что все работает
	services := make([]models.Service, 0)
	endpoints := make([]models.Endpoint, 0)
	methods := []string{"GET", "POST"}
	endpoints = append(endpoints, models.Endpoint{URL: "/users", Protected: false, Methods: methods})
	services = append(services, models.Service{
		Name:      "user",
		Port:      "6969",
		IP:        "localhost",
		Endpoints: endpoints,
	})
	for _, srv := range services {
		fn, ok := list[srv.Name]
		if ok {
			registerService(&srv, gw, fn)
		}
	}
}
*/

func (gw *Gateway) registerDefaultService(ip string, port string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		redirectURL, err := buildURLhandler(ip, port)
		if err != nil {
			gw.Logger.Fatal("Failed to register service ", err)
		}
		proxy := httputil.NewSingleHostReverseProxy(redirectURL)
		r.Header.Set("RedirectURL", redirectURL.String())
		proxy.ServeHTTP(w, r)
	}
}

func buildURLhandler(ip string, port string) (*url.URL, error) {
	redirectURL, err := url.Parse("http://" + ip + ":" + port)
	if err != nil {
		return nil, err
	}

	return redirectURL, nil
}

func registerService(srv *models.Service, gw *Gateway, fn func(string, string) http.HandlerFunc) {
	for _, endpoint := range srv.Endpoints {
		if endpoint.Protected {
			gw.ProtectedRouter.Path(
				endpoint.URL).Handler(fn(srv.IP, srv.Port)).Methods(endpoint.Methods...)
		} else {
			gw.Router.Path(
				endpoint.URL).Handler(fn(srv.IP, srv.Port)).Methods(endpoint.Methods...)
		}
	}

	gw.Logger.Infof("Registred service. Name: %s, IP: %s, Port: %s. Endpoints: %d", srv.Name, srv.IP, srv.Port, len(srv.Endpoints))
}
