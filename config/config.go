package config

type Config struct {
    DeepseekAPIKey string
    GRPCPort       string
}

func LoadConfig() *Config {
    return &Config{
        DeepseekAPIKey: "your-deepseek-api-key", // Replace with your actual API key
        GRPCPort:       ":8080",
    }
}