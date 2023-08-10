package cosmos_governance_bot

import (
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

type Config struct {
	Production      bool             `yaml:"production"`
	EnableLogs      bool             `yaml:"enableLogs"`
	IntervalMinutes int64            `yaml:"interval_minutes"`
	Chains          map[string]Chain `yaml:"chains"`
	BotToken        string           `yaml:"bot_token"`
	Discord         *Discord         `yaml:"discord"`

	DiscordThreads struct {
		EnableThreadsAndReactions bool  `yaml:"enable_threads_and_reactions"`
		Archivethreads            bool  `yaml:"archive_threads"`
		ThreadArchiveMinutes      int16 `yaml:"thread_archive_minutes"`
	} `yaml:"discord_threads"`
}

type Chain struct {
	GrpcUrl        string   `yaml:"grpc_url"`
	ExplorerGovUrl string   `yaml:"explorer_gov_url"`
	Discord        *Discord `yaml:"discord"`
}

type Discord struct {
	WebhookId    string `yaml:"webhook_id"`
	WebhookToken string `yaml:"webhook_token"`
	AvatarUrl    string `yaml:"avatar_url"`
	HexColor     string `yaml:"hex_color"`
	Tags         string `yaml:"tags"`
}

func (ctf Config) Load() Config {
	f, err := os.Open("config.yaml")

	if err != nil {
		log.Fatal("Can't read config.yaml file, error", err)
	}

	defer f.Close()

	var cfg Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		log.Fatal("Could not decode config.yaml file", err)
	}

	return cfg
}
