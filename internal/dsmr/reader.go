package dsmr

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/hierynomus/iot-monitor/pkg/iot"
	"github.com/hierynomus/iot-monitor/pkg/logging"
	"github.com/hierynomus/p1-monitor/internal/crc16"
	"github.com/tarm/serial"
)

var _ iot.Scraper = (*Reader)(nil)

type Reader struct {
	config Config
	port   *serial.Port
	ch     chan iot.RawMessage
	wg     *sync.WaitGroup
}

func NewDsmrReader(config Config) (*Reader, error) {
	ch := make(chan iot.RawMessage)

	port, err := serial.OpenPort(&serial.Config{
		Name:     config.Device,
		Baud:     config.Baud,
		Size:     config.Bits,
		Parity:   config.Parity,
		StopBits: config.StopBits,
	})
	if err != nil {
		return nil, err
	}

	return &Reader{
		config: config,
		port:   port,
		ch:     ch,
		wg:     &sync.WaitGroup{},
	}, nil
}

func (r *Reader) Start(ctx context.Context) error {
	r.wg.Add(1)

	go r.run(ctx)

	return nil
}

func (r *Reader) Stop() error {
	r.port.Close()

	return nil
}

func (r *Reader) Wait() {
	r.wg.Wait()
}

func (r *Reader) Output() <-chan iot.RawMessage {
	return r.ch
}

func (r *Reader) run(ctx context.Context) {
	reader := bufio.NewReader(r.port)
	logger := logging.LoggerFor(ctx, "dsmr")

	defer r.wg.Done()
	defer close(r.ch)

	for {
		if err := r.FindStartTelegram(reader); err != nil {
			if errors.Is(err, io.EOF) {
				logger.Info().Msg("Port closed")
				return
			}
		}

		rawTelegram, err := r.ReadRawTelegram(reader)
		if err != nil && errors.Is(err, io.EOF) {
			logger.Info().Msg("Port closed")
			return
		} else if err != nil {
			logger.Error().Err(err).Msg("Error reading raw telegram")
			continue
		}

		if r.config.ValidateCrc {
			crc, err := r.ReadCrc(reader)
			if err != nil && errors.Is(err, io.EOF) {
				logger.Info().Msg("Port closed")
				return
			} else if err != nil {
				logger.Error().Err(err).Msg("Error reading crc")
				continue
			}

			actualCrc := fmt.Sprintf("%04X", crc16.Checksum(rawTelegram))
			expectedCrc := strings.ToUpper(strings.TrimSpace(string(crc)))
			if expectedCrc != actualCrc {
				logger.Error().Str("expected-crc", expectedCrc).Str("actual-crc", actualCrc).Msg("CRC mismatch")
				continue
			}
		}

		r.ch <- iot.RawMessage(rawTelegram)
	}
}

func (r *Reader) FindStartTelegram(reader *bufio.Reader) error {
	for {
		// Peek at the next byte, and look for the start of the telegram
		if peek, err := reader.Peek(1); err == nil {
			// The telegram starts with a '/' character keep reading
			// bytes until the start of the telegram is found
			if string(peek) == "/" {
				return nil
			}

			reader.ReadByte() //nolint:errcheck
		} else if errors.Is(err, io.EOF) {
			return err
		}
	}
}

func (r *Reader) ReadRawTelegram(reader *bufio.Reader) ([]byte, error) {
	return reader.ReadBytes('!')
}

func (r *Reader) ReadCrc(reader *bufio.Reader) ([]byte, error) {
	return reader.ReadBytes('\n')
}
