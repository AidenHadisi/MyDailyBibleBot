/*
Copyright © 2022 AIDEN HADISI

*/
package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/AidenHadisi/MyDailyBibleBot/configs"
	"github.com/AidenHadisi/MyDailyBibleBot/pkg/bible"
	"github.com/AidenHadisi/MyDailyBibleBot/pkg/bot"
	"github.com/AidenHadisi/MyDailyBibleBot/pkg/cache"
	"github.com/AidenHadisi/MyDailyBibleBot/pkg/cron"
	"github.com/AidenHadisi/MyDailyBibleBot/pkg/image"
	"github.com/AidenHadisi/MyDailyBibleBot/pkg/twitter"
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
		twitterApi := twitter.NewTwitterApi(cfg)
		memoryCache := cache.NewMemoryCache()
		bibleAPI := bible.NewBibleAPI(&http.Client{Timeout: time.Minute}, memoryCache)
		simpleCron := cron.NewSimpleCron()
		imageProcessor := image.NewImageProcessor(&http.Client{Timeout: time.Minute * 2})
		bot := bot.NewBot(cfg, twitterApi, bibleAPI, simpleCron, imageProcessor)

		err = bot.Init()
		if err != nil {
			return err
		}

		// twitterClient := api.NewTwitterApi(goTwitter)

		sc := make(chan os.Signal, 1)
		signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM)
		<-sc

		bot.Shutdown()

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
	rootCmd.Flags().BoolP("dev", "d", false, "run dev mode")
}
