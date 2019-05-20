package main

type Config struct {
	Address string `json:"Address"`
	Name    string `json:"Name"`
	Status  string `json:"Status"`
}

type Response struct {
	Output string
	Type   string
	Status string
}

type Request struct {
	Type  string
	Input map[string]interface{}
}

type Admin struct{}

type ECOM struct{}
