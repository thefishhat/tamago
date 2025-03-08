package main

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/thefishhat/tamago/cli/hotswapmodel"
	"github.com/thefishhat/tamago/cli/views/entities"
	"github.com/thefishhat/tamago/client"
	"github.com/thefishhat/tamago/config"
)

func main() {
	var f *os.File
	var err error
	if len(os.Args) > 1 && os.Args[1] == "--debug" {
		f, err = tea.LogToFile("debug.log", "debug")
		if err != nil {
			fmt.Println("[fatal] creating logger:", err)
			os.Exit(1)
		}
		defer f.Close()
	}

	cfg := config.LoadConfig()
	client := client.NewClient(cfg.Addr)

	entities := entities.NewEntitiesModel(client)
	hotswap := hotswapmodel.New(entities)
	p := tea.NewProgram(hotswap, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		log.Fatal("running program:", err)
	}
}
