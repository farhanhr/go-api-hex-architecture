package app

import (
	"gonews/config"
	"gonews/lib/auth"
	"gonews/lib/middleware"
	"gonews/lib/pagination"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/rs/zerolog/log"
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

	_ = auth.NewJwt(cfg)
	_ = middleware.NewMiddleware(cfg)
	_ = pagination.NewPagination()
}