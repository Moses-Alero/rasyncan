package cmd

import (
	"context"
	"io"
	"log"

	"github.com/spf13/cobra"
)

type SyncConfig struct{}

func (s *SyncConfig) Initialize() {
	root.AddCommand(FileSync)
}

var root = &cobra.Command{
	Use:   "sync",
	Short: "file sync what else",
	Long:  "Say longer thing about project",
}

func Execute(ctx context.Context, config SyncConfig, out io.Writer) error {
	//set logger to Stdout
	log.SetOutput(out)
	config.Initialize()
	//simple for loop running on a different thread
	// this is the daemon process
	go func() {
	LOOP:

		for {
			break LOOP
		}

	}()

	return root.Execute()
}
