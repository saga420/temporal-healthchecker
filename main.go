package main

import (
	"github.com/saga420/temopral-healthchecker/healthchecker"
	"github.com/saga420/temopral-healthchecker/version"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"time"
)

var rootCmd = &cobra.Command{
	Use:   "health-checker",
	Short: "Health checker for temporal services",
	Long:  `This is a health checker for temporal services, it checks the status of FrontendService, HistoryService, and MatchingService.`,
	Run: func(cmd *cobra.Command, args []string) {
		runHealthCheck()
	},
}

func init() {
	rootCmd.PersistentFlags().String("config", "", "config file (default is $HOME/.config.yaml)")
	err := viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))
	if err != nil {
		log.Fatalf("unable to bind flag: %v", err)
	}
	viper.AutomaticEnv()
}

func runHealthCheck() {
	log.Printf("Health Checker GitRevision %s", version.GitRevision)
	var elapsed time.Duration
	start := time.Now()

	viper.SetConfigFile(viper.GetString("config"))
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Fatalf("Config file not found: %v", err)
		} else {
			log.Fatalf("Error reading config file: %v", err)
		}
	}

	var cfg healthchecker.HealthCheckConfig
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}

	hc, err := healthchecker.NewHealthChecker(cfg)
	if err != nil {
		log.Fatalf(healthchecker.FormatError(err))
	}
	defer hc.Close()
	elapsed = time.Since(start)
	log.Printf("NewHealthChecker took %s", elapsed)

	log.Printf("BasicCheck is starting...")
	err = hc.BasicCheck()
	if err != nil {
		log.Fatalf(healthchecker.FormatError(err))
	}

	log.Printf("BasicCheck is done.")
	elapsed = time.Since(start)
	log.Printf("BasicCheck took %s", elapsed)

	log.Printf("FullCheck is starting...")
	err = hc.FullCheck()
	if err != nil {
		log.Fatalf(healthchecker.FormatError(err))
	}
	log.Printf("FullCheck is done.")
	elapsed = time.Since(start)
	log.Printf("FullCheck took %s", elapsed)

	log.Println("All services are healthy.")
}

func main() {
	cobra.CheckErr(rootCmd.Execute())
}
