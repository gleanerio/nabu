package cli

import (
	"fmt"
	"github.com/gleanerio/nabu/pkg"
	log "github.com/sirupsen/logrus"
	"os"

	"github.com/spf13/cobra"
)

// checkCmd represents the check command
var tikaCmd = &cobra.Command{
	Use:   "tika",
	Short: "nabu tika command",
	Long:  `(not implemented)This will read the configs/{cfgPath}/gleaner file, and try to connect to the minio server`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("tike called")
		err := pkg.Tika(viperVal, mc)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		os.Exit(0)
	},
}

func init() {
	rootCmd.AddCommand(tikaCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// checkCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// checkCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
