package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/iccolo/rma/analyzer"
	"github.com/iccolo/rma/gui/server/analyze"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	//statikFS, err := fs.New()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//http.Handle("/", http.FileServer(statikFS))
	http.HandleFunc("/api/rma/get_instance_list", GetInstanceList)
	http.HandleFunc("/api/rma/start_analyze", StartAnalyze)
	http.HandleFunc("/api/rma/get_key_type", GetKeyType)
	http.HandleFunc("/api/rma/expand", Expand)
	http.HandleFunc("/api/rma/get_key_info", GetKeyInfo)
	if err := http.ListenAndServe(":8090", nil); err != nil {
		log.Fatal(err)
	}
}

var h = analyze.NewHandler()

func GetInstanceList(response http.ResponseWriter, request *http.Request) {
	if intercept(response, request, nil) {
		return
	}
	list := h.GetInstanceList()
	out, _ := json.Marshal(list)
	log.Println(string(out))
	response.Write(out)
}

func StartAnalyze(response http.ResponseWriter, request *http.Request) {
	a := &analyzer.Analyzer{}
	if intercept(response, request, a) {
		return
	}
	log.Printf("%+v\n", a)
	if a.Separators == "" {
		response.WriteHeader(http.StatusBadRequest)
	}
	h.StartAnalyze(a)
}

func GetKeyType(response http.ResponseWriter, request *http.Request) {
	type In struct {
		Host string `json:"host"`
	}
	in := &In{}
	if intercept(response, request, in) {
		return
	}
	keyTypes, err := h.GetKeyTypes(in.Host)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		return
	}
	out, _ := json.Marshal(keyTypes)
	log.Println(string(out))
	response.Write(out)
}

func Expand(response http.ResponseWriter, request *http.Request) {
	type In struct {
		Host      string `json:"host"`
		KeyType   string `json:"key_type"`
		KeyPrefix string `json:"key_prefix"`
		NumLimit  int64  `json:"num_limit"`
		SortVar   int32  `json:"sort_var"`
	}
	in := &In{}
	if intercept(response, request, in) {
		return
	}
	nodeList, err := h.Expand(in.Host, in.KeyType, in.KeyPrefix, in.NumLimit, analyze.SortVar(in.SortVar))
	if err != nil {
		log.Printf("Expand err:%v, in:%+v", err, in)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	out, _ := json.Marshal(nodeList)
	log.Println("out:", string(out))
	response.Write(out)
	return
}

func GetKeyInfo(response http.ResponseWriter, request *http.Request) {
	type In struct {
		Host string `json:"host"`
		Key  string `json:"key"`
	}
	in := &In{}
	if intercept(response, request, in) {
		return
	}
	keyInfo, err := h.GetKeyInfo(in.Host, in.Key, 10)
	if err != nil {
		log.Printf("GetKeyInfo err:%v, in:%+v", err, in)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	out, _ := json.Marshal(keyInfo)
	log.Println("out:", string(out))
	response.Write(out)
	return
}

// intercept unmarshal body and return if intercept process logic
func intercept(response http.ResponseWriter, request *http.Request, body interface{}) bool {
	response.Header().Set("Access-Control-Allow-Origin", "*")
	response.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	response.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization,X-CSRF-Token")
	if request.Method == "OPTIONS" {
		return true
	}
	binaryBody, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Printf("ReadAll request.Body err:%v\n", err)
		response.WriteHeader(http.StatusBadRequest)
		return true
	}
	if body != nil {
		log.Println("req body:", string(binaryBody))
		if err = json.Unmarshal(binaryBody, body); err != nil {
			log.Printf("json.Unmarshal request.Body:%v, err:%v\n", string(binaryBody), err)
			response.WriteHeader(http.StatusBadRequest)
			return true
		}
	}
	return false
}
