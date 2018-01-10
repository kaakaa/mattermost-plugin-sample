package main

import (
	"fmt"
)

type Configuration struct {
	Hogeru string
}

func (c *Configuration) IsValid() error {
	if len(c.Hogeru) < 2 {
		return fmt.Errorf("Minimum length is 2")
	}
	return nil
}