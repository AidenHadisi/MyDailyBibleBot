/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"net/http"
	"os"
	"time"

	"github.com/AidenHadisi/MyDailyBibleBot/api"
	"github.com/AidenHadisi/MyDailyBibleBot/configs"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mydailybiblebot",
	Short: "Start my daily bible bot app.",
	Long:  "Start my daily bible bot app.",

	// Uncomment the following line if your bare application
	// has an action associated with it:
	RunE: func(cmd *cobra.Command, args []string) error {
		dev, err := cmd.Flags().GetBool("dev")
		if err != nil {
			return err
		}
		cfg := configs.LoadConfig(dev)

		httpClient := &http.Client{Timeout: time.Minute * 2}

		config := oauth1.NewConfig(cfg.ConsumerKey, cfg.ConsumerSecret)
		token := oauth1.NewToken(cfg.AccessToken, cfg.AccessSecret)
		httpClient := config.Client(oauth1.NoContext, token)
		goTwitter := twitter.NewClient(httpClient)

		bibleApi := api.BibleAPI()

		twitterClient := api.NewTwitterApi(goTwitter)

		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.MyDailyBibleBot.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
