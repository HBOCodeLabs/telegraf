package appdynamics

import (
	"fmt"
	"net/http"
	"testing"
	_ "encoding/json"
	"encoding/json"
	"io/ioutil"
)

func TestHTTPController(t *testing.T) {
	url := "https://hbo-go.saas.appdynamics.com/controller/rest/applications/hurley.staging/tiers/comet?output=JSON"
	client := &http.Client{}

	/* Auth */
	req, err := http.NewRequest("GET", url, nil)
	req.SetBasicAuth("apirouser@hbo-go","HBORocks2!")

	res, err := client.Do(req)
	if err != nil {
		panic(err.Error)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err.Error())
	}

	var tiers []struct{
		Id int64 `json:"id"`
	}

	err = json.Unmarshal(body, &tiers)
	if err != nil {
		fmt.Println("Error: ", err)
	}
	if (len(tiers) != 1){
		fmt.Println("Invalid reply: ", tiers)
	}
	fmt.Println("Tier Id: ", tiers[0].Id)

}
