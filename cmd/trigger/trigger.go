package httpTrigger

import (
	"github.com/ind9/rocks/cmd/ops"
	"github.com/spf13/cobra"
)

var trigger = &cobra.Command{
	Use:   "trigger",
	Short: "Triggers a backup on a remote system",
	Long:  `Triggers a backup on a remote system`,
}

func init() {
	ops.Rocks.AddCommand(trigger)
}