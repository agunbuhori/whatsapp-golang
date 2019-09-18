package main

import (
	"os"
)

func main() {
	deleteSession()
}

func deleteSession() error {
	return os.Remove(os.TempDir() + "/whatsappSession.gob")
}
