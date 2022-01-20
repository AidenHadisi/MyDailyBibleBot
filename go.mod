module github.com/AidenHadisi/MyDailyBibleBot

// +heroku goVersion go1.17      <--add to

go 1.17

require (
	github.com/dghubble/go-twitter v0.0.0-20201011215211-4b180d0cc78d
	github.com/dghubble/oauth1 v0.7.0
	github.com/fogleman/gg v1.3.0
	github.com/go-co-op/gocron v1.11.0
	github.com/mitchellh/mapstructure v1.4.3
	github.com/nfnt/resize v0.0.0-20180221191011-83c6a9932646
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/spf13/cobra v1.3.0
	github.com/stretchr/testify v1.7.0
)

require (
	github.com/cenkalti/backoff v2.2.1+incompatible // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dghubble/sling v1.3.0 // indirect
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/robfig/cron/v3 v3.0.1 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/stretchr/objx v0.1.1 // indirect
	golang.org/x/image v0.0.0-20211028202545-6944b10bf410 // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)

replace github.com/dghubble/go-twitter => github.com/AidenHadisi/go-twitter v0.0.0-20190721142740-110a39637298
