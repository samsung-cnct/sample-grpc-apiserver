package configs

import (
	"github.com/spf13/viper"
)

var (
	envPrefix                = "grpc_hello"
	envServiceAddress        = "service_address"
	envServicePort           = "service_port"
	envGateWayServiceAddress = "gw_service_address"
	envGateWayPort           = "gw_port"
	envClientName            = "client_name"
	envGateWaySwaggerDir     = "gw_swagger_dir"
	addressDefault           = "localhost"
	portDefault              = 5300
	gwAddressDefault         = "localhost"
	gwPortDefault            = 6300
	gwSwaggerDirDefault      = "swagger"
	clientNameDefault        = "trident"
)

// InitEnvVars allows you to initiate gathering environment variables
func InitEnvVars() error {
	var err error
	viper.SetEnvPrefix(envPrefix)

	if err = viper.BindEnv(envServicePort); err != nil {
		return err
	}

	if err = viper.BindEnv(envServiceAddress); err != nil {
		return err
	}

	if err = viper.BindEnv(envGateWayServiceAddress); err != nil {
		return err
	}

	if err = viper.BindEnv(envGateWayPort); err != nil {
		return err
	}

	if err = viper.BindEnv(envClientName); err != nil {
		return err
	}

	err = viper.BindEnv(envGateWaySwaggerDir)

	return err
}

// ParseGWSwaggerEnvVars parses environment variables consumed by swagger server
func ParseGWSwaggerEnvVars() string {
	gwSwaggerDir := viper.GetString(envGateWaySwaggerDir)
	if gwSwaggerDir == "" {
		gwSwaggerDir = gwSwaggerDirDefault
	}
	return gwSwaggerDir
}

// ParseGateWayEnvVars parses environment variables consumed by the gateway service
func ParseGateWayEnvVars() (int, int, string, string) {
	gwPort := viper.GetInt(envGateWayPort)
	if gwPort == 0 {
		gwPort = gwPortDefault
	}

	port := viper.GetInt(envServicePort)
	if port == 0 {
		port = portDefault
	}

	gwServiceAddress := viper.GetString(envGateWayServiceAddress)
	if gwServiceAddress == "" {
		gwServiceAddress = gwAddressDefault
	}

	serviceAddress := viper.GetString(envServiceAddress)
	if serviceAddress == "" {
		serviceAddress = addressDefault
	}

	return gwPort, port, gwServiceAddress, serviceAddress
}

// ParseClientEnvVars parses environment variables consumed by clients
func ParseClientEnvVars() string {
	clientName := viper.GetString(envClientName)
	if clientName == "" {
		clientName = clientNameDefault
	}

	return clientName
}
