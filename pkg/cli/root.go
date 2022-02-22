package cli

import (
	"log"
	"mime"
	"os"
	"path"
	"path/filepath"

	"github.com/gleanerio/nabu/internal/objects"
	"github.com/gleanerio/nabu/pkg/config"
	"github.com/minio/minio-go/v7"
	"github.com/spf13/cobra"

	"github.com/spf13/viper"
)

var cfgFile, cfgName, cfgPath, nabuConfName string
var minioVal, portVal, accessVal, secretVal, bucketVal string
var sslVal bool
var viperVal *viper.Viper
var mc *minio.Client
var prefixVal string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "nabu",
	Short: "nabu ",
	Long: `nabu
`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	log.SetFlags(log.Lshortfile)

	mime.AddExtensionType(".jsonld", "application/ld+json")

	log.Println("EarthCube Nabu")
	akey := os.Getenv("MINIO_ACCESS_KEY")
	skey := os.Getenv("MINIO_SECRET_KEY")
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&prefixVal, "prefix", "", "prefix to run")

	rootCmd.PersistentFlags().StringVar(&cfgPath, "cfgPath", "configs", "base location for config files (default is configs/)")
	rootCmd.PersistentFlags().StringVar(&cfgName, "cfgName", "local", "config file (default is local so configs/local)")
	rootCmd.PersistentFlags().StringVar(&nabuConfName, "nabuConfName", "nabu", "config file (default is local so configs/local)")
	rootCmd.PersistentFlags().StringVar(&cfgFile, "cfg", "", "compatibility/overload: full path to config file (default location gleaner in configs/local)")

	// minio env variables
	rootCmd.PersistentFlags().StringVar(&minioVal, "address", "localhost", "FQDN for server")
	rootCmd.PersistentFlags().StringVar(&portVal, "port", "9000", "Port for minio server, default 9000")
	rootCmd.PersistentFlags().StringVar(&accessVal, "access", akey, "Access Key ID")
	rootCmd.PersistentFlags().StringVar(&secretVal, "secret", skey, "Secret access key")
	rootCmd.PersistentFlags().StringVar(&bucketVal, "bucket", "gleaner", "The configuration bucket")

	rootCmd.PersistentFlags().BoolVar(&sslVal, "ssl", false, "Use SSL boolean")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	var err error
	//viperVal := viper.New()
	if cfgFile != "" {
		// Use config file from the flag.
		//viperVal.SetConfigFile(cfgFile)
		viperVal, err = config.ReadNabuConfig(filepath.Base(cfgFile), filepath.Dir(cfgFile))
		if err != nil {
			log.Fatal("cannot read config %s", err)
		}
	} else {
		// Find home directory.
		//home, err := os.UserHomeDir()
		//cobra.CheckErr(err)
		//
		//// Search config in home directory with name "nabu" (without extension).
		//viperVal.AddConfigPath(home)
		//viperVal.AddConfigPath(path.Join(cfgPath, cfgName))
		//viperVal.SetConfigType("yaml")
		//viperVal.SetConfigName("nabu")
		viperVal, err = config.ReadNabuConfig(nabuConfName, path.Join(cfgPath, cfgName))
		if err != nil {
			log.Fatal("cannot read config %s", err)
		}
	}

	//viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.

	// Set up some logging approaches
	f, err := os.OpenFile("naburun.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)
	log.SetFlags(log.Lshortfile)
	// log.SetOutput(ioutil.Discard) // turn off all logging
	//wrt := io.MultiWriter(os.Stdout, f)
	//log.SetOutput(wrt)

	mc, err = objects.MinioConnection(viperVal)
	if err != nil {
		log.Fatal("cannot connect to minio: %s", err)
	}

	bucketVal, err = config.GetBucketName(viperVal)
	if err != nil {
		log.Println("cannot read bucketname from : %s ", err)
	}
	// Override prefix in config if flag set
	//if isFlagPassed("prefix") {
	//	out := viperVal.GetStringMapString("objects")
	//	b := out["bucket"]
	//	p := prefixVal
	//	// r := out["region"]
	//	// v1.Set("objects", map[string]string{"bucket": b, "prefix": NEWPREFIX, "region": r})
	//	viperVal.Set("objects", map[string]string{"bucket": b, "prefix": p})
	//}
	if prefixVal != "" {
		out := viperVal.GetStringMapString("objects")
		b := out["bucket"]
		p := prefixVal
		// r := out["region"]
		// v1.Set("objects", map[string]string{"bucket": b, "prefix": NEWPREFIX, "region": r})
		viperVal.Set("objects", map[string]string{"bucket": b, "prefix": p})
	}

}
