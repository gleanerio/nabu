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
var objectCmd = &cobra.Command{
	Use:   "object",
	Short: "nabu object command",
	Long:  `(not implemented)This will read the configs/{cfgPath}/gleaner file, and try to connect to the minio server`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("object called")
		err := Object(viperVal, mc)
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

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// checkCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// checkCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func Object(v1 *viper.Viper, mc *minio.Client) error {
	fmt.Println("Load graph object to triplestore")
	spql := v1.GetStringMapString("sparql")
	s, err := flows.PipeLoad(v1, mc, "bucket", "object", spql["endpoint"])
	if err != nil {
		log.Println(err)
	}
	log.Println(string(s))
	return err
}
