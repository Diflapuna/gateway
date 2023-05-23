package models

type Services struct {
	Services []Service `json:"services"`
}
type Service struct {
	Name      string     `json:"name"`
	Port      string     `json:"port"`
	IP        string     `json:"ip"`
	Endpoints []Endpoint `json:"endpoints"`
}

type Endpoint struct {
	URL       string   `json:"url"`
	Protected bool     `json:"protected"`
	Methods   []string `json:"methods"`
}

type Greeting struct {
	Greeting string `json:"greeting"`
}
