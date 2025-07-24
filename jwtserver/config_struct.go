package main

type Config struct {
	// Configuration settings for the REST server.
	Server struct {
		// Port on which the REST service is available.
		Port int `yaml:"port"`
	}
	// token attributes
	Token struct {
		Audience     string `yaml:"audience"`
		Issuer       string `yaml:"issuer"`
		ValidMinutes int    `yaml:"valid_minutes"`
	}
}
