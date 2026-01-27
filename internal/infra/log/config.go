package log

type Config struct {
	Level      string `mapstructure:"level"`       // debug / info / warn / error
	Encoding   string `mapstructure:"encoding"`    // json / console
	Output     string `mapstructure:"output"`      // stdout / stderr / file
	FilePath   string `mapstructure:"file_path"`   // 当 output=file
	TimeFormat string `mapstructure:"time_format"` // RFC3339 / custom
}
