package cmd

import (
	"github.com/spf13/cobra"
)


func init(){
	root.AddCommand(FileSync)

}

var root = &cobra.Command{
	Use: "sync",
	Short: "file sync what else",
	Long: "Say longer thing about project",
}

func Execute() error{
	return root.Execute()
}
