package get

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/wizedkyle/sumocli/api"
	"github.com/wizedkyle/sumocli/pkg/cmd/factory"
	"github.com/wizedkyle/sumocli/pkg/logging"
	"io"
)

func NewCmdGetUser() *cobra.Command {
	var id string

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Gets a Sumo Logic user",
		Run: func(cmd *cobra.Command, args []string) {
			GetUser(id)
		},
	}

	cmd.Flags().StringVar(&id, "id", "", "Specify the id of the user to get")
	cmd.MarkFlagRequired("id")
	return cmd
}

func GetUser(id string) {
	var userInfo api.UserResponse
	log := logging.GetConsoleLogger()
	requestUrl := "v1/users/" + id
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

	err = json.Unmarshal(responseBody, &userInfo)
	if err != nil {
		log.Error().Err(err).Msg("failed to unmarshal response body")
	}

	userInfoJson, err := json.MarshalIndent(userInfo, "", "    ")
	if err != nil {
		log.Error().Err(err).Msg("failed to marshal response")
	}

	if response.StatusCode != 200 {
		factory.HttpError(response.StatusCode, responseBody, log)
	} else {
		fmt.Println(string(userInfoJson))
	}
}
