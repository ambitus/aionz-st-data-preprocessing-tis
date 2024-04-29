package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"os"

	"github.com/gorilla/mux"
)

type responseBody struct {
	Data []string `json:"data"`
}

func ai_inference_triton(input_data []string) map[string]interface{} {
	// Construct triton input data
	triton_data := `{"inputs":[{"name":"IN0","shape":[1,10],"datatype":"BYTES","data":[[`
	for i := 0; i < len(input_data)-1; i++ {
		triton_data += `"` + input_data[i] + `",`
	}
	triton_data += `"` + input_data[len(input_data)-1] + `"]]}],"outputs":[{"name":"OUT0"}]}`

	log.Println(triton_data)

	requestBody := []byte(triton_data)

	resp, err := http.Post(os.Getenv("SCORING_URL") + ":" + os.Getenv("SCORING_PORT") + "/v2/models/rf_model/infer", "application/json", bytes.NewBuffer(requestBody))

	if err != nil {
		log.Println(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Println(err)
	}

	log.Println(string(body))

	var jsonMap map[string]interface{}
	json.Unmarshal([]byte(string(body)), &jsonMap)
	log.Println(jsonMap)

	return jsonMap
}

func preprocess(w http.ResponseWriter, r *http.Request) []string {
	w.Header().Set("Content-Type", "application/json")

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprint(w, err.Error())
	}
	log.Println(data)

	var Data responseBody

	json.Unmarshal(data, &Data)

	// Time
	log.Println(Data.Data[5])
	time := strings.ReplaceAll(Data.Data[5], ":", "")
	log.Println(time)
	Data.Data[5] = time

	// Amount
	log.Println(Data.Data[6])
	amount := strings.ReplaceAll(Data.Data[6], "$", "")
	log.Println(amount)
	Data.Data[6] = amount

	// Use Chip
	log.Println(Data.Data[7])
	use_chip_str := Data.Data[7]

	if use_chip_str == "Swipe Transaction" {
		use_chip_int := 0
		log.Println(use_chip_int)
	} else {
		use_chip_int := 1
		log.Println(use_chip_int)
	}

	return Data.Data
}

func ai_inference(w http.ResponseWriter, r *http.Request) {
	preprocessed_data := preprocess(w, r)
	triton_resp := ai_inference_triton(preprocessed_data)

	json.NewEncoder(w).Encode(triton_resp)
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/ai_inference", ai_inference).Methods("PUT")

	log.Fatal(http.ListenAndServe(":8080", r))
}
