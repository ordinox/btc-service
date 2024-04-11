package inscriptions

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/alexellis/go-execute/v2"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/ordinox/btc-service/config"
	"github.com/rs/zerolog/log"
)

var ordNonZeroExitCodeErr = fmt.Errorf("non-zero-exit-code")

// Retry 10 times if the ord throws a non-zero exit code err (Usually when ord can't find a db lock)
func Inscribe(inscription, destination string, feeRate uint64, config config.BtcConfig) (res *InscriptionResultRaw, err error) {
	count := 0
	for {
		res, err = inscribe(inscription, destination, feeRate, config)

		if res != nil || !errors.Is(err, ordNonZeroExitCodeErr) {
			return res, err
		}
		// Log retry attempt if the specific retryable error occurred
		if count >= 9 { // Adjusted to 9 to correctly handle up to 10 retries
			return nil, fmt.Errorf("exceeded retry attempts due to persistent error: %w", err)
		}

		log.Debug().Msg("Retrying inscribing")
		time.Sleep(time.Second * 1)
		count = count + 1
	}
}

func inscribe(inscription, destination string, feeRate uint64, config config.BtcConfig) (*InscriptionResultRaw, error) {
	file, err := os.CreateTemp("", "*.txt")
	if err != nil {
		log.Error().Msgf("Error writing inscription to the temp file - %s", err)
		return nil, err
	}
	_, err = file.WriteString(inscription)
	if err != nil {
		log.Error().Msgf("Error writing inscription to the temp file - %s", err)
		return nil, err
	}
	args := []string{config.GetOrdChainConfigFlag(), "--bitcoin-data-dir", config.BitcoinDataDir, "--data-dir", config.OrdDataDir, "wallet", "inscribe", "--fee-rate", fmt.Sprintf("%d", feeRate), "--destination", destination, "--file", file.Name(), "--postage", "546sat"}
	if config.GetChainConfigParams().Name == chaincfg.MainNetParams.Name {
		args = args[1:]
	}
	cmd := execute.ExecTask{
		Command:      strings.TrimRight(config.OrdPath, "/") + "/ord",
		Args:         args,
		StreamStdio:  false,
		PrintCommand: true,
	}
	time.Sleep(time.Second * 1)
	res, err := cmd.Execute(context.Background())
	if err != nil {
		log.Error().Msgf("Error sending inscription %s %v with err: %s", cmd.Command, cmd.Args, res.Stderr)
		return nil, err
	}
	if res.ExitCode != 0 {
		log.Error().Msgf("Non 0 exit code while sending inscription %s | %s", res.Stderr, res.Stdout)
		return nil, ordNonZeroExitCodeErr
	}

	data := &InscriptionResultRaw{}
	err = json.Unmarshal([]byte(res.Stdout), data)
	if err != nil {
		log.Err(fmt.Errorf(res.Stderr)).Msgf("Error packing raw inscription output")
		return nil, err
	}
	return data, nil
}
