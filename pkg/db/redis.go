package db

import (
	"github.com/redis/rueidis"
	"github.com/spf13/viper"
)

func newrdb(dbn int) (rueidis.Client, error) {
	return rueidis.NewClient(rueidis.ClientOption{
		InitAddress: viper.GetStringSlice("redis.hosts"),
		Username:    viper.GetString("redis.username"),
		Password:    viper.GetString("redis.password"),
		SelectDB:    dbn,
	})
}

type OneTimeTokenStore rueidis.Client

func NewOneTimeTokenStore() (OneTimeTokenStore, error) {
	return newrdb(0)
}
