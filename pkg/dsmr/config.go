package dsmr

import "github.com/tarm/serial"

type Config struct {
	Device   string          `yaml:"device" viper:"device" env:"DEVICE" default:"/dev/ttyUSB0"`
	Baud     int             `yaml:"baud" env:"BAUD" viper:"baud" default:"115200"`
	Bits     byte            `yaml:"bits" viper:"bits" ENV:"BITS" default:"8"`
	Parity   serial.Parity   `yaml:"parity" viper:"parity" ENV:"PARITY" default:"N"`
	StopBits serial.StopBits `yaml:"stopbits" viper:"stopbits" ENV:"STOPBITS" default:"1"`

	ValidateCrc bool `yaml:"validate-crc" viper:"validate-crc" ENV:"VALIDATE_CRC" default:"true"`
}
