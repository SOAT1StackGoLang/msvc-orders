package helpers

import (
	"fmt"
	"os"
	"strconv"
)

var (
	postgresUsername string
	postgresPassword string
	postgresHost     string
	postgresPort     int
	postgresDbName   string
	postgresURI      string
)

func ReadPgxConnEnvs() {
	postgresURI = os.Getenv("DB_URI")
	postgresUsername = os.Getenv("DB_USER")
	postgresPassword = os.Getenv("DB_PASSWORD")
	postgresHost = os.Getenv("DB_HOST")
	port, _ := strconv.ParseInt(os.Getenv("DB_PORT"), 10, 64)
	postgresPort = int(port)
	postgresDbName = os.Getenv("DB_NAME")
}

func GetConnectionParams() string {
	if postgresURI != "" {
		return postgresURI
	}
	return ToDsnWithDbName()
}

func ToDsnWithDbName() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable", postgresHost, postgresUsername, postgresPassword, postgresDbName, postgresPort)
}

func ToDsn() string {
	return fmt.Sprintf("host=%s user=%s password=%s port=%d sslmode=disable", postgresHost, postgresUsername, postgresPassword, postgresPort)
}
