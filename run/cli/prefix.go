package cli

import (
	"fmt"
	"github.com/gleanerio/nabu/internal/flows"
	"github.com/minio/minio-go/v7"
	"github.com/spf13/viper"
	"log"
	"os"

	"github.com/spf13/cobra"
)

// checkCmd represents the check command
var PrefixCmd = &cobra.Command{
	Use:   "prefix ",
	Short: "nabu prefix command",
	Long:  `Load graphs from prefix to triplestore`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("prefix called")
		err := Prefix(viperVal, mc)
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
	// checkCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func Prefix(v1 *viper.Viper, mc *minio.Client) error {

	log.Println("Load graphs from prefix to triplestore")
	err := flows.ObjectAssembly(v1, mc)

	if err != nil {
		log.Println(err)
	}
	return err
}
