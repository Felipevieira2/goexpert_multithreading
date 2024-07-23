package main

import (
	"encoding/json"
	"fmt"
	"time"

	"net/http"
)

type Message struct {
	Msg string
}

type Formattable interface {
	Format() string
}

type ViaCEPResponse struct {
	Cep        string `json:"cep"`
	Logradouro string `json:"logradouro"`
	Bairro     string `json:"bairro"`
	Localidade string `json:"localidade"`
	Uf         string `json:"uf"`
}

type BrasilAPIResponse struct {
	Cep        string `json:"cep"`
	Logradouro string `json:"street"`
	Bairro     string `json:"neighborhood"`
	Localidade string `json:"city"`
	Uf         string `json:"state"`
}

func (b *BrasilAPIResponse) Format() string {
	//print all fields and value
	return fmt.Sprintf(
		"CEP: %s\nLogradouro: %s\nBairro: %s\nLocalidade: %s\nUF: %s\n Consulta feita por: Brasil APi\n",
		b.Cep, b.Logradouro, b.Bairro, b.Localidade, b.Uf)

}

func (v *ViaCEPResponse) Format() string {
	//print all fields and value
	return fmt.Sprintf(
		"CEP: %s\nLogradouro: %s\nBairro: %s\nLocalidade: %s\nUF: %s\n Consulta feita por: Via CEP\n",
		v.Cep, v.Logradouro, v.Bairro, v.Localidade, v.Uf)
}

func main() {
	c1 := make(chan Message)

	go requestApi("https://brasilapi.com.br/api/cep/v1/01153000", &BrasilAPIResponse{}, c1)
	go requestApi("http://viacep.com.br/ws/01153000/json", &ViaCEPResponse{}, c1)

	select {
	case msg := <-c1:
		println(msg.Msg)
	case <-time.After(time.Second * 1):
		println("timeout")
		// default:
		// 	println("default")
	}
}

func requestApi(url string, responseStruct Formattable, c1 chan Message) {
	//guardar o tempo de resposta em uma variável

	timeStart := time.Now()
	resp, err := http.Get(url)

	if err != nil {
		c1 <- Message{"Failed to make request to " + url}
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c1 <- Message{"Non-OK HTTP status: " + resp.Status}
		return
	}

	if err := json.NewDecoder(resp.Body).Decode(responseStruct); err != nil {
		c1 <- Message{"Failed to decode response from " + url}
		return
	}

	timeEnd := time.Now()

	fmt.Printf("\nRequest to %s took %s\n", url, timeEnd.Sub(timeStart))
	// time.Sleep(1 * time.Second) simuta a resposta de 1 segundo para timeout

	c1 <- Message{responseStruct.Format()}

}
