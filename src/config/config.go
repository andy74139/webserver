package config

import (
	"fmt"

	"github.com/joho/godotenv"
)

var envFile map[string]string

func init() {
	env, err := godotenv.Read("env")
	if err != nil {
		panic("error loading env file")
	}

	envFile = env

}

func GetLocalDSN() string {
	name := getEnvPanic("DB_NAME")
	port := getEnvPanic("DB_PORT")
	password := getEnvPanic("DB_PASSWORD")
	host := "localhost"
	return fmt.Sprintf("postgres://postgres:%s@%s:%s/%s?sslmode=disable", password, host, port, name)
}

func GetContainerDSN() string {
	name := getEnvPanic("DB_NAME")
	port := getEnvPanic("DB_PORT")
	password := getEnvPanic("DB_PASSWORD")
	host := getEnvPanic("DB_SERVICE_NAME")
	return fmt.Sprintf("postgres://postgres:%s@%s:%s/%s?sslmode=disable", password, host, port, name)
}

func GetCacheAddress() string {
	serviceName := getEnvPanic("CACHE_SERVICE_NAME")
	port := getEnvPanic("CACHE_PORT")
	return fmt.Sprintf("%s:%s", serviceName, port)
}

func GetWebServerPort() string {
	webPort := getEnvPanic("WEB_PORT")
	return webPort
}

func getEnv(arg string) string {
	val, ok := envFile[arg]
	if !ok {
		return ""
	}
	return val
}

func getEnvPanic(arg string) string {
	val, ok := envFile[arg]
	if !ok {
		panic(fmt.Errorf("env variable %s not set", arg))
	}
	return val
}
