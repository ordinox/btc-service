package inscriptions

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/alexellis/go-execute/v2"
	"github.com/ordinox/btc-service/config"
	"github.com/rs/zerolog/log"
)

func Inscribe(inscription, destination string, config config.BtcConfig) (*InscriptionResultRaw, error) {
	fee, err := getAvgFee()
	if err != nil {
		return nil, err
	}
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
	args := []string{config.GetOrdChainConfigFlag(), "wallet", "inscribe", "--fee-rate", fmt.Sprintf("%d", fee), "--destination", destination, "--file", file.Name()}
	cmd := execute.ExecTask{
		Command:     "ord",
		Args:        args,
		StreamStdio: false,
	}
	time.Sleep(time.Second * 1)
	res, err := cmd.Execute(context.Background())
	if err != nil {
		log.Error().Msgf("Error sending inscription %s %v with err: %s", cmd.Command, cmd.Args, res.Stderr)
		return nil, err
	}
	if res.ExitCode != 0 {
		log.Error().Msgf("Non 0 exit code while sending inscription %s %v with err: %s", cmd.Command, cmd.Args, res.Stderr)
		return nil, err
	}

	data := &InscriptionResultRaw{}
	err = json.Unmarshal([]byte(res.Stdout), data)
	if err != nil {
		log.Error().Msgf("Error packing raw inscription output %s %v with err: %s", cmd.Command, cmd.Args, res.Stderr)
		return nil, err
	}
	return data, nil
}
