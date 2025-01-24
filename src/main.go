package main

import (
	"fmt"
	"os"
	"time"

	"net"
	"encoding/json"

	pb_outputs "github.com/VU-ASE/rovercom/packages/go/outputs"
	roverlib "github.com/VU-ASE/roverlib-go/src"

	"github.com/rs/zerolog/log"
)

type UDPObject struct {
    Channel int     `json:"channel"`
    Value   float64 `json:"value"`
}

// The main user space program
// this program has all you need from roverlib: service identity, reading, writing and configuration
func run(service roverlib.Service, configuration *roverlib.ServiceConfiguration) error {
	if configuration == nil {
		return fmt.Errorf("configuration cannot be accessed")
	}


	port, err := configuration.GetStringSafe("udp-port")
	if err != nil {
		return fmt.Errorf("failed to get configuration: %v", err)
	}
	log.Info().Msgf("Fetched runtime configuration UDP port number: %s", port)


	writeStream := service.GetWriteStream("decision")
	if writeStream == nil {
		return fmt.Errorf("failed to get write stream")
	}


	pc, err := net.ListenPacket("udp", port)
	if err != nil {
		return fmt.Errorf("failed to access port %s", port)
	}
	defer pc.Close()

	log.Info().Msgf("Listening on port: %s", port)

	for {
		
		buf := make([]byte, 1024)

		n, _, err := pc.ReadFrom(buf)
		if err != nil {
			log.Error().Msgf("Error encountered while receiving a packet")
			continue
		}

		log.Info().Msgf("Received a message")

		var command UDPObject
		// json.Unmarshal(buf[:n], &command)
		err = json.Unmarshal(buf[:n], &command)
		if err != nil {
			log.Error().Msgf("Failed to unmarshal JSON: %v", err)
			continue
		}

		log.Info().Msgf("Unmarshalled jason info: channel: %d, value: %f", command.Channel, command.Value)


		var result pb_outputs.ControllerOutput
		result.FrontLights = false

		switch channel := command.Channel; channel {
		case 0:
			result.SteeringAngle = float32(command.Value)
		case 1:
			result.LeftThrottle = float32(command.Value)
		case 2:
			result.RightThrottle = float32(command.Value)
		default:
			log.Error().Msgf("Unrecognized value in the [channel] field. Expected: [0-2], got: %d", channel)
		}


		err = writeStream.Write(
			&pb_outputs.SensorOutput{
				SensorId:  2,
				Timestamp: uint64(time.Now().UnixMilli()),
				SensorOutput: &pb_outputs.SensorOutput_ControllerOutput{
					ControllerOutput: &result,
				},
			},
		)
		// Send it for the actuator (and others) to use
		if err != nil {
			log.Err(err).Msg("Failed to send controller output")
			continue
		}

		log.Debug().Msg("Sent controller output")
	}
}

// This function gets called when roverd wants to terminate the service
func onTerminate(sig os.Signal) error {
	log.Info().Str("signal", sig.String()).Msg("Terminating service")

	//
	// ...
	// Any clean up logic here
	// ...
	//

	return nil
}

// This is just a wrapper to run the user program
// it is not recommended to put any other logic here
func main() {
	roverlib.Run(run, onTerminate)
}
