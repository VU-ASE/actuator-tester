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

// required json format:
// channel :  0 for steering servo, 1 for left motor, 2 for right motor
// value : from -1 to 1
type ChannelCommand struct {
    Channel int     `json:"channel"`
    Value   float64 `json:"value"`
}

// The main user space program
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

	// open a connection over a UDP port specified in service.yaml
	connection, err := net.ListenPacket("udp", port)
	if err != nil {
		return fmt.Errorf("failed to access port %s", port)
	}
	defer connection.Close()

	log.Info().Msgf("Listening on port: %s", port)

	for {
		// read incoming data from the connection
		buf := make([]byte, 1024)
		nBytesRead, _, err := connection.ReadFrom(buf)
		if err != nil {
			log.Error().Msgf("Error encountered while receiving a packet")
			continue
		}
		log.Info().Msgf("Received a message")

		// decode the raw data into a useable format
		var command ChannelCommand
		err = json.Unmarshal(buf[:nBytesRead], &command)
		if err != nil {
			log.Error().Msgf("Failed to unmarshal JSON: %v", err)
			continue
		}
		log.Info().Msgf("Unmarshalled json info: channel: %d, value: %f", command.Channel, command.Value)

		if command.Value > 1 {
			command.Value = 1
			log.Warn().Msgf("Read value greater than 1. Setting to 1")
		}

		if command.Value < -1 {
			command.Value = -1
			log.Warn().Msgf("Read value less than -1. Setting to -1")
		}

		// format the command as an output stream readable by other services
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
			result.SteeringAngle = float32(0)
			result.LeftThrottle = float32(0)
			result.RightThrottle = float32(0)
			log.Error().Msgf("Unrecognized value in the [channel] field. Expected: [0-2], got: %d", channel)
		}

		// Send it for the actuator (and others) to use
		err = writeStream.Write(
			&pb_outputs.SensorOutput{
				SensorId:  2,
				Timestamp: uint64(time.Now().UnixMilli()),
				SensorOutput: &pb_outputs.SensorOutput_ControllerOutput{
					ControllerOutput: &result,
				},
			},
		)
		if err != nil {
			log.Err(err).Msg("Failed to send tester output")
			continue
		}

		log.Debug().Msg("Sent controller output")
	}
}

// This function gets called when roverd wants to terminate the service
func onTerminate(sig os.Signal) error {
	log.Info().Str("signal", sig.String()).Msg("Terminating service")
	return nil
}

// This is just a wrapper to run the user program
func main() {
	roverlib.Run(run, onTerminate)
}
