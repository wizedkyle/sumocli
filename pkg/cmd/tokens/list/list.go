package list

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/wizedkyle/sumocli/api"
	"github.com/wizedkyle/sumocli/pkg/cmd/factory"
	"github.com/wizedkyle/sumocli/pkg/logging"
	"io"
)

func NewCmdTokensList() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Get a list of all tokens in the token library.",
		Run: func(cmd *cobra.Command, args []string) {
			listTokens()
		},
	}
	return cmd
}

func listTokens() {
	var tokenResponse api.ListTokenResponse
	log := logging.GetConsoleLogger()
	requestUrl := "v1/tokens"
	client, request := factory.NewHttpRequest("GET", requestUrl)
	response, err := client.Do(request)
	if err != nil {
		log.Error().Err(err).Msg("failed to make http request to " + requestUrl)
	}

	defer response.Body.Close()
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		log.Error().Err(err).Msg("failed to read response body")
	}

	err = json.Unmarshal(responseBody, &tokenResponse)
	if err != nil {
		log.Error().Err(err).Msg("failed to unmarshal response body")
	}

	tokenResponseJson, err := json.MarshalIndent(tokenResponse, "", "    ")
	if err != nil {
		log.Error().Err(err).Msg("failed to marshal lookupTableResponse")
	}

	if response.StatusCode != 200 {
		factory.HttpError(response.StatusCode, responseBody, log)
	} else {
		fmt.Println(string(tokenResponseJson))
	}
}
