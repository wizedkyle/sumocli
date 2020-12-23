package create

import (
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/tidwall/gjson"
	"github.com/wizedkyle/sumocli/api"
	"github.com/wizedkyle/sumocli/pkg/cmd/factory"
	"github.com/wizedkyle/sumocli/pkg/logging"
	"io/ioutil"
	"strings"
)

func NewCmdCollectorCreate() *cobra.Command {
	var (
		name        string
		description string
		category    string
		fields      string
		output      string
	)
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates a Sumo Logic Hosted Collector",
		Run: func(cmd *cobra.Command, args []string) {
			log := logging.GetConsoleLogger()
			CreateCollector(name, description, category, fields, output, log)
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "Specify the name of the collector")
	cmd.Flags().StringVar(&description, "description", "", "Specify a description for the collector")
	cmd.Flags().StringVar(&category, "category", "", "Source category for the collector, this will overwrite the categories on configured sources")
	cmd.Flags().StringVar(&fields, "fields", "", "Specify fields for the collector, they need to be formatted as field1:value1,field2:value2")
	cmd.Flags().StringVar(&output, "output", "", "Specify the field to export the value from")
	cmd.MarkFlagRequired("name")
	return cmd
}

func CreateCollector(name string, description string, category string, fields string, output string, log zerolog.Logger) {
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
		log.Fatal().Err(err).Msg("failed to make http request")
	}

	defer response.Body.Close()
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal().Err(err).Msg("error reading response body from request")
	}

	jsonErr := json.Unmarshal(responseBody, &createCollectorResponse)
	if jsonErr != nil {
		log.Fatal().Err(jsonErr).Msg("error unmarshalling response body")
	}

	createCollectorResponseJson, err := json.MarshalIndent(createCollectorResponse, "", "    ")

	if response.StatusCode != 201 {
		factory.HttpError(response.StatusCode, responseBody, log)
	} else {
		if factory.ValidateCollectorOutput(output) == true {
			value := gjson.Get(string(createCollectorResponseJson), output)
			formattedValue := strings.Trim(value.String(), `"[]"`)
			fmt.Println(formattedValue)
		} else {
			fmt.Println(string(createCollectorResponseJson))
			log.Info().Msg(createCollectorResponse.Collector.Name + " collector successfully created")
		}
	}
}
