package main

import (
	"app/infrastructure/config"
	"app/infrastructure/log"
	"app/infrastructure/search"
	"app/service"

	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{Use: config.Get().Name}
	rootCmd.AddCommand(serveCommand(), sonicCommand())
	rootCmd.Execute()
}

func serveCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "serve",
		Short: "Run HTTP Server",
		Run: func(cmd *cobra.Command, args []string) {
			service.Run()
		},
	}
}

func sonicCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "sonic",
		Short: "Sonic Commands",
	}
	var syncCmd = &cobra.Command{
		Use:   "sync",
		Short: "Sync words & dicts to Sonic",
		Run: func(cmd *cobra.Command, args []string) {
			err := search.Sync()
			if err != nil {
				log.Error(err.Error())
			}
		},
	}
	cmd.AddCommand(syncCmd)
	return cmd
}
