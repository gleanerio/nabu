package cli

import (
	"fmt"
	"github.com/gleanerio/nabu/internal/prune"
	"github.com/minio/minio-go/v7"
	"github.com/spf13/viper"
	"log"
	"os"

	"github.com/spf13/cobra"
)

// checkCmd represents the check command
var pruneCmd = &cobra.Command{
	Use:   "prune",
	Short: "nabu prune command",
	Long:  `(not implemented)This will read the configs/{cfgPath}/gleaner file, and try to connect to the minio server`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("prune called")
		err := Prune(viperVal, mc)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		os.Exit(0)
	},
}

func init() {
	rootCmd.AddCommand(pruneCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// checkCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// checkCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func Prune(v1 *viper.Viper, mc *minio.Client) error {
	fmt.Println("Prune graphs in triplestore not in object store")
	err := prune.Snip(v1, mc)
	if err != nil {
		log.Println(err)
	}
	return err
}
