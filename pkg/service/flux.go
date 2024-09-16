package service

import (
	"fmt"

	fluxpb "github.com/Winens/flux/proto/gen/go"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

func NewFluxImageProcessingService() (flux fluxpb.FluxImageClient, err error) {
	var creds credentials.TransportCredentials

	tlsEnabled := viper.GetBool("flux.tls.enabled")
	if tlsEnabled {
		creds, err = credentials.NewClientTLSFromFile(viper.GetString("flux.tls.certFile"), "")
		if err != nil {
			return nil, err
		}

	} else {
		creds = insecure.NewCredentials()
	}

	address := fmt.Sprintf("%s:%d", viper.GetString("flux.host"), viper.GetInt("flux.port"))
	fluxConn, err := grpc.NewClient(address, grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to flux service")
	}

	return fluxpb.NewFluxImageClient(fluxConn), nil
}
