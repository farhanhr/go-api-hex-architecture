package app

import (
	"gonews/config"

	"github.com/rs/zerolog/log"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func RunServer() {
	cfg := config.NewConfig()
	_ , err := cfg.ConnectionPostgres()

	if err != nil {
		log.Fatal().Msgf("Error connecting to database: %v", err)
	}

	// Cloudflare R2

	cdfR2 := cfg.LoadAWSConfig()
	_ = s3.NewFromConfig(cdfR2)
}