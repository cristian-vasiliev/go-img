package config

type Config struct {
	Port               int
	StaticDir          string
	DefaultQuality     int
	DefaultImageFormat string
}

func NewConfig() *Config {
	return &Config{
		// Port is the default port for the server to listen on
		Port: 8080,
		// StaticDir is the default directory for the static files
		StaticDir: "static",
		// DefaultQuality is the default quality for the image from 0 to 100
		DefaultQuality: 75,
		// DefaultImageFormat is the default image format for the image (jpeg, png, webp, avif)
		DefaultImageFormat: "jpeg",
	}
}
