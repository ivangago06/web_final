package config

import (
	"encoding/json"
	"flag"
	"log"
	"os"
)

type Configuration struct {
	Db struct {
		Ip          string `json:"DbIp"`
		Port        string `json:"DbPort"`
		Name        string `json:"DbName"`
		User        string `json:"DbUser"`
		Password    string `json:"DbPassword"`
		AutoMigrate bool   `json:"DbAutoMigrate"`
	}

	Listen struct {
		Ip   string `json:"HttpIp"`
		Port string `json:"HttpPort"`
	}

	Template struct {
		BaseName string `json:"BaseTemplateName"`
	}
}

func LoadConfig() Configuration {
	c := flag.String("c", "env.json", "Ruta para el archivo JSON")
	flag.Parse()
	file, err := os.Open(*c)
	if err != nil {
		log.Fatal("Error al abrir el archivo JSON: ", err)
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal("Error al cerrar el archivo JSON: ", err)
		}
	}(file)

	decoder := json.NewDecoder(file)
	Config := Configuration{}
	err = decoder.Decode(&Config)
	if err != nil {
		log.Fatal("Error al decodificar el archivo JSON ", err)
	}

	return Config
}
