package config

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
)

func (cfg Config) LoadAWSConfig() aws.Config {
	conf, err := awsConfig.LoadDefaultConfig(context.TODO(),
		awsConfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			cfg.R2.ApiKey, cfg.R2.ApiSecret, "",
		)), awsConfig.WithRegion("auto"), )

	if err != nil {
		log.Fatal().Msgf("unable to load aws config %v", err)
	}

	log.Info().Msg("Success Loaded AWS config")

	return conf
}