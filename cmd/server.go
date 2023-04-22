package cmd

import (
	"fmt"
	"log"
	"vssd/server"

	"github.com/spf13/cobra"
)

var (
	serverHost    string
	serverPort    int
	serverName    string
	serverSize    string
	serverVerbose bool
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Create a server server for VSSD",
	Run: func(cmd *cobra.Command, args []string) {
		size := formatSizeString(serverSize)
		if size == -1 {
			fmt.Println("Invalid size of SSD, please use a valid size (e.g. 1KB, 1MB, 1GB, 1B, 1)")
			return
		}

		s := server.New(serverHost, serverPort, serverName, size, serverVerbose)
		if err := s.Start(); err != nil {
			log.Panicln(err)
		}
	},
}

func init() {
	serverCmd.Flags().StringVarP(&serverHost, "host", "l", "localhost", "host of the SSD")
	serverCmd.Flags().IntVarP(&serverPort, "port", "p", 1337, "port of the SSD")
	serverCmd.Flags().StringVarP(&serverName, "name", "n", "default", "name of the SSD")
	serverCmd.Flags().StringVarP(&serverSize, "size", "s", "1KB", "size of the SSD (e.g. 1KB, 1MB, 1GB, 1B, 1)")
	serverCmd.Flags().BoolVarP(&serverVerbose, "verbose", "v", false, "verbose mode")

	rootCmd.AddCommand(serverCmd)
}
