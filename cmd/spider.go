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
	spiderCmd.PersistentFlags().Bool("json", false, "output json format")
	spiderCmd.PersistentFlags().Bool("xml", false, "output xml format")
	spiderCmd.AddCommand(spider.NewAdocCmd())
	return spiderCmd
}
