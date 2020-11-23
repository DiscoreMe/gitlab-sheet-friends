package main

import (
	"fmt"
	"log"

	"github.com/DiscoreMe/gitlab-sheets-friends/config"
	"github.com/DiscoreMe/gitlab-sheets-friends/service"
	"github.com/DiscoreMe/gitlab-sheets-friends/sheets"
	"github.com/DiscoreMe/gitlab-sheets-friends/storage"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func main() {
	log.SetFlags(0)

	conf, err := config.LoadConfig("config.yaml")
	if err != nil {
		panic(err)
	}
	stor, err := storage.NewStorage(conf.DB)
	if err != nil {
		panic(err)
	}

	serv, err := service.NewService(conf, stor)
	if err != nil {
		if e, ok := err.(sheets.ErrTokenNotFoundOrOutdated); ok {
			logError(e)
			log.Println("please add a token by entering a command in the terminal: google add")
			log.Println()
		} else {
			log.Fatal(err)
		}
	}

	var rootCmd = &cobra.Command{}
	rootCmd.AddCommand(&cobra.Command{
		Use:   "run",
		Short: "Run",
		Run: func(cmd *cobra.Command, args []string) {
			if err := serv.Run(); err != nil {
				log.Fatal(err)
			}
		},
	})
	rootCmd.AddCommand(&cobra.Command{
		Use:   "ping",
		Short: "Pings to all git projects",
		Run: func(cmd *cobra.Command, args []string) {
			serv.Ping()
		},
	})

	var googleCmd = &cobra.Command{
		Use:   "google",
		Short: "Commands for working with the Google API token",
		Args:  cobra.MinimumNArgs(1),
	}
	googleCmd.AddCommand(&cobra.Command{
		Use:   "update",
		Short: "update Google API token",
		Run: func(cmd *cobra.Command, args []string) {
			if err := serv.UpdateToken(); err != nil {
				log.Fatalln(err)
			}
			log.Println("token was updated successfully")
		},
	})
	googleCmd.AddCommand(&cobra.Command{
		Use:   "columns",
		Short: "update sheets columns",
		Run: func(cmd *cobra.Command, args []string) {
			if err := serv.UpdateColumns(); err != nil {
				log.Fatalln(err)
			}
		},
	})

	rootCmd.AddCommand(googleCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func logError(err error) {
	_, _ = color.New(color.FgRed).Print("err:")
	fmt.Printf(" %s\n", err)
}
