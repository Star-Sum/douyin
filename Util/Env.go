package Util

import "os"

func EnvTransfer(env string) string {
	return os.Getenv(env)
}
