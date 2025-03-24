package main

import (
	"fmt"
	"os"
	"time"

	"encoding/json"
	"net"

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

// If UDP tuning is enabled, then this service will accept values JSON objects over
// a UDP connection and send them to the actuator.
func fetchCommandOverUdp(connection net.PacketConn) (*ChannelCommand, error) {
	// read incoming data from the connection
	buf := make([]byte, 1024)
	nBytesRead, _, err := connection.ReadFrom(buf)
	if err != nil {
		log.Error().Msgf("Error encountered while receiving a packet")
		return nil, err
	}
	log.Info().Msgf("Received a message")

	var command ChannelCommand

	// decode the raw data into a useable format
	// var command ChannelCommand
	err = json.Unmarshal(buf[:nBytesRead], &command)
	if err != nil {
		log.Error().Msgf("Failed to unmarshal JSON: %v", err)
		return nil, err
	}
	log.Info().Msgf("Unmarshalled json info: channel: %d, value: %f", command.Channel, command.Value)

	return &command, nil
}

// Given an channelCommand and a write stream this function will perform necessary
// checks and write the channelCommand to the output stream.
func outputCommand(command *ChannelCommand, writeStream *roverlib.WriteStream) {
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

	log.Info().Msgf("Setting channel %d to %f", command.Channel, command.Value)

	// Send it for the actuator (and others) to use
	err := writeStream.Write(
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
	}

	log.Debug().Msg("Sent controller output")
}

func run(service roverlib.Service, configuration *roverlib.ServiceConfiguration) error {
	if configuration == nil {
		return fmt.Errorf("configuration cannot be accessed")
	}

	port, err := configuration.GetStringSafe("udp-port")
	if err != nil {
		return fmt.Errorf("failed to get configuration: %v", err)
	}
	log.Info().Msgf("Listening for JSON objects over UDP on port %s", port)
	

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

	// The current command publish
	var command *ChannelCommand

	for {
		command, err = fetchCommandOverUdp(connection)

		if err != nil {
			log.Error().Msgf("%v", err)
		} else {
			outputCommand(command, writeStream)
		}
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
