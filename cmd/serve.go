package cmd

import (
	"github.com/Techassi/gomark/internal/app"
	"github.com/Techassi/gomark/internal/server"

	"github.com/spf13/cobra"
)

var configPath string

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "A brief description of your command",
	Long:  `A longer description that spans multiple lines`,
	Run: func(cmd *cobra.Command, args []string) {
		a := app.New(configPath)
		s := server.New(a)
		s.Run()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	serveCmd.PersistentFlags().StringVar(&configPath, "config", "", "Path to your config file")
	cobra.MarkFlagRequired(serveCmd.PersistentFlags(), "config")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
