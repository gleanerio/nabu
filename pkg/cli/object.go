package cli

import (
	"fmt"
	"github.com/gleanerio/nabu/pkg"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var bucket, objectVal string

// checkCmd represents the check command
var objectCmd = &cobra.Command{
	Use:   "objectVal",
	Short: "nabu objectVal command",
	Long:  `Load graph objectVal to triplestore`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("objectVal called")
		if objectVal == "" {

		}
		err := pkg.Object(viperVal, mc, bucketVal, objectVal)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		os.Exit(0)

	},
}

func init() {
	rootCmd.AddCommand(objectCmd)

	// Here you will define your flags and configuration settings.
	// bucketVal is available at top level
	objectCmd.Flags().StringVar(&objectVal, "objectVal", "", "objectVal to load")
	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// checkCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// checkCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
