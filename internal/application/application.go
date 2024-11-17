package application

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/DobryySoul/yandex_repo/pkg/calculation"
)

type Config struct {
	Addr string
}

func ConfigFromEnv() *Config {
	config := new(Config)
	config.Addr = os.Getenv("PORT")
	if config.Addr == "" {
		config.Addr = "8080"
	}

	return config
}

type Application struct {
	config *Config
}

func New() *Application {
	return &Application{
		config: ConfigFromEnv(),
	}
}

func (a *Application) Run() error {
	for {
		log.Println("input expression")
		reader := bufio.NewReader(os.Stdin)
		text, err := reader.ReadString('\n')
		if err != nil {
			log.Println("failed to read expred to read from console")
		}

		text = strings.TrimSpace(text)
		if text == "exit" {
			log.Println("application was successfully closed")
			return nil
		}

		result, err := calculation.Calc(text)
		if err != nil {
			log.Println("failed to calculate expression with error: ", err)
		} else {
			log.Printf("%s = %f", text, result)
		}
	}
}

type Expression struct {
	Expression string `json:"expression"`
}

func CalcHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	expressionRequest := new(Expression)
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(expressionRequest)
	if err != nil {
		http.Error(w, "unexpected end of JSON input", http.StatusBadRequest)
		return
	}
	result, err := calculation.Calc(expressionRequest.Expression)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Fprint(w, result)
}

func (a *Application) RunServer() error {
	http.HandleFunc("/", CalcHandler)
	log.Printf("server started on port :%s", a.config.Addr)
	if err := http.ListenAndServe(":"+a.config.Addr, nil); err != nil {
		log.Fatalf("server is drop with error: %s", err)
	}
	return nil
}
