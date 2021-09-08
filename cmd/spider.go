package cmd

import (
	"github.com/spf13/cobra"

	"github.com/la0wan9/ark/cmd/spider"
)

// NewSpiderCmd creates a new spider command
func NewSpiderCmd() *cobra.Command {
	var spiderCmd = &cobra.Command{
		Use: "spider",
	}
	spiderCmd.Flags().SortFlags = false
	spiderCmd.AddCommand(spider.NewAdocCmd())
	return spiderCmd
}
