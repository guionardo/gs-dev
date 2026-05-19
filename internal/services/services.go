package services

type (
	Service        any
	NewServiceFunc func(configurationFolder string) Service
)
