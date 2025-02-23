package api

import (
	"fmt"
	"workwavebot/server"
)

func HandleDeleteAdmin(id string) error {
	err := server.DeleteAdmin(id)
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}
	return nil
}
