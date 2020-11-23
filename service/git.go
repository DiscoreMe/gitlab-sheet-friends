package service

import (
	"fmt"
	"github.com/fatih/color"
)

func (s *Service) Ping() {
	for name, g := range s.gits {
		fmt.Printf("Pinging %s...   ", name)
		if err := g.Ping(); err != nil {
			color.Red("FAILED")
			fmt.Println("error: ", err)
			fmt.Println()
		} else {
			color.Green("SUCCESS")
		}
	}
}
