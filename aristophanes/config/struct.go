package config

import "github.com/odysseia-greek/aristoteles"

type Config struct {
	PodName   string
	Namespace string
	Index     string
	Elastic   aristoteles.Client
}
