//go:build wireinject
// +build wireinject

package main

import (
	"github.com/AidenHadisi/MyDailyBibleBot/api"
	"github.com/AidenHadisi/MyDailyBibleBot/bot"
	"github.com/AidenHadisi/MyDailyBibleBot/config"
	"github.com/google/wire"
)

func InitializeBot(cfg *config.Config) *bot.Bot {

	wire.Build(
		bot.NewBot,
		api.BibleApiWire,
		api.TwitterWire,
	)
	return &bot.Bot{}
}
