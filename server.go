package desafioclientserverapi

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type Cambio struct {
	Bid string `json:"bid"`
}

func main() {
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log.Println("Request iniciada")
	defer log.Println("Request finalizada")
}

func BuscaCambio() (*Cambio, error) {
	res, err := http.Get("https://economia.awesomeapi.com.br/json/last/USD-BRL")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var c Cambio
	err = json.Unmarshal(body, &c)
	if err != nil {
		return nil, err
	}
	return &c, nil

}
