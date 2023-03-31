// Handling the settings for the bot.
package settings

// Settings for the Bot
type Settings struct {
	AllowedUsers string `yaml:"allowed_users" bson:"allowed_users"`
	BotToken     string `yaml:"bot_token" bson:"bot_token"`
	StartMessage string `yaml:"start_message" bson:"start_text"`
	StopMessage  string `yaml:"start_message" bson:"stop_text"`
	HelpText     string `yaml:"start_message" bson:"help_text"`
}
