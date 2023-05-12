package models

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
