package cmd

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"
)

// Current is to identify a rocksdb store.
const Current = "CURRENT"

// LatestBackup is used to find the backup location
const LatestBackup = "LATEST_BACKUP"

// Rocks is the entry point command in the application
var Rocks = &cobra.Command{
	Use:   "rocks",
	Short: "RocksDB Ops CLI",
	Long: `Perform common ops related tasks on one or many RocksDB instances.

Find more details at https://github.com/ind9/rocks`,
}

// CommandHandler is the wrapper interface that all commands to be implement as part of their "Run"
type CommandHandler func(args []string) error

// AttachHandler is a wrapper method for all commands that needs to be exposed
func AttachHandler(handler CommandHandler) func(*cobra.Command, []string) {
	return func(cmd *cobra.Command, args []string) {
		start := time.Now()
		err := handler(args)
		elapsed := time.Since(start).Seconds()
		fmt.Printf("This took  %f seconds\n", elapsed)
		if err != nil {
			log.Printf("[Error] %s", err.Error())
			os.Exit(1)
		}
	}
}
