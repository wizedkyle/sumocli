package create

import (
	"encoding/json"
	"github.com/spf13/cobra"
	"github.com/wizedkyle/sumocli/api"
	"github.com/wizedkyle/sumocli/pkg/cmd/factory"
	"github.com/wizedkyle/sumocli/pkg/logging"
	"io"
	"strings"
)

func NewCmdCollectorCreate() *cobra.Command {
	var (
		name        string
		description string
		category    string
		fields      string
	)
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates a Sumo Logic Hosted Collector",
		Run: func(cmd *cobra.Command, args []string) {
			Collector(name, description, category, fields)
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "Specify the name of the collector")
	cmd.Flags().StringVar(&description, "description", "", "Specify a description for the collector")
	cmd.Flags().StringVar(&category, "category", "", "sourceCategory for the collector, this will overwrite the categories on configured sources")
	cmd.MarkFlagRequired("name")
	return cmd
}

func Collector(name string, description string, category string, fields string) api.CollectorResponse {
	log := logging.GetConsoleLogger()
	var createCollectorResponse api.CollectorResponse
	requestBodySchema := &api.CreateCollectorRequest{
		Collector: api.CreateCollector{
			CollectorType: "Hosted",
			Name:          name,
			Description:   description,
			Category:      category,
			Fields:        nil,
		},
	}

	if fields != "" {
		fieldsMap := make(map[string]string)
		splitStrings := strings.Split(fields, ",")
		for i, splitString := range splitStrings {
			components := strings.Split(splitString, ":")
			fieldsMap[components[0]] = components[1]
			i++
		}
		requestBodySchema.Collector.Fields = fieldsMap
	}

	requestBody, _ := json.Marshal(requestBodySchema)
	client, request := factory.NewHttpRequestWithBody("POST", "v1/collectors", requestBody)
	response, err := client.Do(request)
	if err != nil {
		log.Error().Err(err).Msg("failed to make http request")
	}

	defer response.Body.Close()
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		log.Error().Err(err).Msg("error reading response body from request")
	}

	err = json.Unmarshal(responseBody, &createCollectorResponse)
	if err != nil {
		log.Error().Err(err).Msg("error unmarshalling response body")
	}
	return createCollectorResponse
}
