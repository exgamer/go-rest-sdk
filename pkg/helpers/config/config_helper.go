package config

import (
	"encoding/json"
	"fmt"
	"github.com/exgamer/go-rest-sdk/pkg/config/structures"
	"github.com/spf13/viper"
	"log"
	"os"
	"strings"
)

func InitBaseConfig() (*structures.AppConfig, *structures.DbConfig, error) {
	appConfig := &structures.AppConfig{}
	dbConfig := &structures.DbConfig{}

	workingDir, err := os.Getwd()

	if err != nil {
		log.Fatalf("Some error occurred. Err: %s", err)
	}

	if _, err := os.Stat(workingDir + "/.env"); err == nil {
		viper.SetConfigFile(workingDir + "/.env")
		viper.SetConfigType("env")
		viper.AutomaticEnv()
		//Find and read the config file
		err := viper.ReadInConfig()

		if err != nil {
			log.Fatalf("Some error occured. Err: %s", err)
		}

		viper.Unmarshal(&appConfig)
		viper.Unmarshal(&dbConfig)
	}

	envKeys := append(appConfig.GetFieldsAsUpperSnake(), dbConfig.GetFieldsAsUpperSnake()...)

	osEnvMap := make(map[string]string)

	for _, key := range envKeys {
		if value, exists := os.LookupEnv(key); exists {
			key = strings.ToLower(key)
			osEnvMap[key] = fmt.Sprint(value)
		}
	}

	//	// Convert the map to JSON
	jsonData, _ := json.Marshal(osEnvMap)
	// Convert the JSON to a struct
	json.Unmarshal(jsonData, &appConfig)
	json.Unmarshal(jsonData, &dbConfig)

	return appConfig, dbConfig, nil
}
