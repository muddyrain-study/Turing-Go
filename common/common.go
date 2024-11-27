package common

import (
	"Turing-Go/config"
	"Turing-Go/models"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"sync"
)

var Template models.HtmlTemplate

func LoadTemplate() {
	w := sync.WaitGroup{}
	w.Add(1)
	go func() {
		var err error
		Template, err = models.InitTemplate(config.Cfg.System.CurrentDir + "/template/")
		if err != nil {
			panic(err)
		}
		w.Done()
	}()
	w.Wait()
}

func Success(w http.ResponseWriter, data interface{}) {
	result := models.Result{
		Code:  200,
		Error: "",
		Data:  data,
	}
	resultJson, _ := json.Marshal(result)
	w.Header().Set("Content-Type", "application/json")
	_, err := w.Write(resultJson)
	if err != nil {
		log.Println("返回数据失败", err)
	}
}

func GetRequestJsonParams(r *http.Request) map[string]interface{} {
	var params map[string]interface{}
	body, _ := io.ReadAll(r.Body)
	_ = json.Unmarshal(body, &params)
	return params
}

func Error(w http.ResponseWriter, err error) {
	result := models.Result{
		Code:  500,
		Error: err.Error(),
	}
	resultJson, _ := json.Marshal(result)
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(resultJson)
	if err != nil {
		log.Println("返回数据失败", err)
	}
}
