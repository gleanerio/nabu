package cli

import (
	"os"

	"github.com/gleanerio/nabu/pkg"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var Source string

// checkCmd represents the check command
var PrefixCmd = &cobra.Command{
	Use:   "prefix ",
	Short: "nabu prefix command",
	Long:  `Load graphs from prefix to triplestore`,
	Run: func(cmd *cobra.Command, args []string) {

		if Source != "" {
			m := viperVal.GetStringMap("objects")
			m["prefix"] = Source
			viperVal.Set("objects", m)
		}

		err := pkg.Prefix(viperVal, mc)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		os.Exit(0)
	},
}

func init() {
	rootCmd.AddCommand(PrefixCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// checkCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	//checkCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	PrefixCmd.Flags().StringVarP(&Source, "source", "s", "", "Source prefix to load")
}
