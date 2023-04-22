package cmd

import (
	"fmt"
	"log"
	"vssd/api"

	"github.com/spf13/cobra"
)

var (
	apiHost    string
	apiPort    int
	apiName    string
	apiSize    string
	apiVerbose bool
)

var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "Create a http server for VSSD",
	Run: func(cmd *cobra.Command, args []string) {
		size := formatSizeString(serverSize)
		if size == -1 {
			fmt.Println("Invalid size of SSD, please use a valid size (e.g. 1KB, 1MB, 1GB, 1B, 1)")
			return
		}

		a := api.New(apiHost, apiPort, apiName, size, apiVerbose)
		if err := a.Start(); err != nil {
			log.Panicln(err)
		}
	},
}

func init() {
	apiCmd.Flags().StringVarP(&apiHost, "host", "l", "localhost", "host of the SSD")
	apiCmd.Flags().IntVarP(&apiPort, "port", "p", 8080, "port of the SSD")
	apiCmd.Flags().StringVarP(&apiName, "name", "n", "default", "name of the SSD")
	apiCmd.Flags().StringVarP(&apiSize, "size", "s", "1KB", "size of the SSD (e.g. 1KB, 1MB, 1GB, 1B, 1)")
	apiCmd.Flags().BoolVarP(&apiVerbose, "verbose", "v", false, "verbose mode")

	rootCmd.AddCommand(apiCmd)
}
