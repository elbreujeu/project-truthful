package client

import "os"

func GetServerVersion() string {
	return os.Getenv("SERVER_VERSION")
}
