package az

import (
	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2020-06-01/web"
	"github.com/Azure/go-autorest/autorest/to"
)

func ReaderAppSettings(storageAccountName string, storageAccountConnection string, instrumentationKey *string, eventHubKey *string,
	serviceBusKey *string) []web.NameValuePair {
	readerAppSettings := []web.NameValuePair{
		{Name: to.StringPtr("FUNCTIONS_EXTENSION_VERSION"), Value: to.StringPtr("~1")},
		{Name: to.StringPtr("Project"), Value: to.StringPtr("BlockBlobReader/target/producer_build/")},
		{Name: to.StringPtr("AzureWebJobsDashboard"), Value: to.StringPtr(storageAccountConnection)},
		{Name: to.StringPtr("AzureWebJobsStorage"), Value: to.StringPtr(storageAccountConnection)},
		{Name: to.StringPtr("APPINSIGHTS_INSTRUMENTATIONKEY"), Value: instrumentationKey},
		{Name: to.StringPtr("TABLE_NAME"), Value: to.StringPtr("FileOffsetMap")},
		{Name: to.StringPtr("AzureEventHubConnectionString"), Value: eventHubKey},
		{Name: to.StringPtr("TaskQueueConnectionString"), Value: serviceBusKey},
		{Name: to.StringPtr("WEBSITE_NODE_DEFAULT_VERSION"), Value: to.StringPtr("6.5.0")},
		{Name: to.StringPtr("FUNCTION_APP_EDIT_MODE"), Value: to.StringPtr("readwrite")},
		{Name: to.StringPtr("WEBSITE_CONTENTAZUREFILECONNECTIONSTRING"), Value: to.StringPtr(storageAccountConnection)},
		{Name: to.StringPtr("WEBSITE_CONTENTSHARE"), Value: to.StringPtr(storageAccountName)},
	}
	return readerAppSettings
}

func ConsumerAppSettings(storageAccountName string, storageAccountConnection string, instrumentationKey *string,
	serviceBusKey *string, sumoLogicSource string) []web.NameValuePair {
	consumerAppSettings := []web.NameValuePair{
		{Name: to.StringPtr("FUNCTIONS_EXTENSION_VERSION"), Value: to.StringPtr("~1")},
		{Name: to.StringPtr("Project"), Value: to.StringPtr("BlockBlobReader/target/consumer_build/")},
		{Name: to.StringPtr("AzureWebJobsDashboard"), Value: to.StringPtr(storageAccountConnection)},
		{Name: to.StringPtr("AzureWebJobsStorage"), Value: to.StringPtr(storageAccountConnection)},
		{Name: to.StringPtr("APPINSIGHTS_INSTRUMENTATIONKEY"), Value: instrumentationKey},
		{Name: to.StringPtr("SumoLogEndpoint"), Value: to.StringPtr(sumoLogicSource)},
		{Name: to.StringPtr("TaskQueueConnectionString"), Value: serviceBusKey},
		{Name: to.StringPtr("WEBSITE_NODE_DEFAULT_VERSION"), Value: to.StringPtr("6.5.0")},
		{Name: to.StringPtr("FUNCTION_APP_EDIT_MODE"), Value: to.StringPtr("readwrite")},
		{Name: to.StringPtr("WEBSITE_CONTENTAZUREFILECONNECTIONSTRING"), Value: to.StringPtr(storageAccountConnection)},
		{Name: to.StringPtr("WEBSITE_CONTENTSHARE"), Value: to.StringPtr(storageAccountName)},
	}
	return consumerAppSettings
}

func DlqAppSettings(storageAccountName string, storageAccountConnection string, instrumentationKey *string,
	serviceBusKey *string, queueName string, sumoLogicSource string) []web.NameValuePair {
	dlqAppSettings := []web.NameValuePair{
		{Name: to.StringPtr("FUNCTIONS_EXTENSION_VERSION"), Value: to.StringPtr("~1")},
		{Name: to.StringPtr("Project"), Value: to.StringPtr("BlockBlobReader/target/dlqprocessor_build/")},
		{Name: to.StringPtr("AzureWebJobsDashboard"), Value: to.StringPtr(storageAccountConnection)},
		{Name: to.StringPtr("AzureWebJobsStorage"), Value: to.StringPtr(storageAccountConnection)},
		{Name: to.StringPtr("APPINSIGHTS_INSTRUMENTATIONKEY"), Value: instrumentationKey},
		{Name: to.StringPtr("SumoLogEndpoint"), Value: to.StringPtr(sumoLogicSource)},
		{Name: to.StringPtr("TaskQueueConnectionString"), Value: serviceBusKey},
		{Name: to.StringPtr("TASKQUEUE_NAME"), Value: to.StringPtr(queueName)},
		{Name: to.StringPtr("WEBSITE_NODE_DEFAULT_VERSION"), Value: to.StringPtr("6.5.0")},
		{Name: to.StringPtr("FUNCTION_APP_EDIT_MODE"), Value: to.StringPtr("readwrite")},
		{Name: to.StringPtr("WEBSITE_CONTENTAZUREFILECONNECTIONSTRING"), Value: to.StringPtr(storageAccountConnection)},
		{Name: to.StringPtr("WEBSITE_CONTENTSHARE"), Value: to.StringPtr(storageAccountName)},
	}
	return dlqAppSettings
}

func DiagnosticLogsAppSettings(logsStorageAccountName string, logsStorageAccountConnection string, failedStorageAccountConnection string,
	eventHubKey *string, sumoLogicSource string) []web.NameValuePair {
	diagnosticAppSettings := []web.NameValuePair{
		{Name: to.StringPtr("FUNCTIONS_EXTENSION_VERSION"), Value: to.StringPtr("~1")},
		{Name: to.StringPtr("Project"), Value: to.StringPtr("EventHubs/target/logs_build/")},
		{Name: to.StringPtr("AzureWebJobsDashboard"), Value: to.StringPtr(logsStorageAccountConnection)},
		{Name: to.StringPtr("AzureWebJobsStorage"), Value: to.StringPtr(logsStorageAccountConnection)},
		{Name: to.StringPtr("SumoLogsEndpoint"), Value: to.StringPtr(sumoLogicSource)},
		{Name: to.StringPtr("AzureEventHubConnectionString"), Value: eventHubKey},
		{Name: to.StringPtr("StorageConnectionString"), Value: to.StringPtr(failedStorageAccountConnection)},
		{Name: to.StringPtr("WEBSITE_CONTENTAZUREFILECONNECTIONSTRING"), Value: to.StringPtr(logsStorageAccountConnection)},
		{Name: to.StringPtr("WEBSITE_CONTENTSHARE"), Value: to.StringPtr(logsStorageAccountName)},
	}
	return diagnosticAppSettings
}

func MetricsAppSettings(logsStorageAccountName string, logsStorageAccountConnection string, failedStorageAccountConnection string,
	eventHubKey *string, sumoLogicSource string) []web.NameValuePair {
	metricsAppSettings := []web.NameValuePair{
		{Name: to.StringPtr("FUNCTIONS_EXTENSION_VERSION"), Value: to.StringPtr("~1")},
		{Name: to.StringPtr("Project"), Value: to.StringPtr("EventHubs/target/metrics_build/")},
		{Name: to.StringPtr("AzureWebJobsDashboard"), Value: to.StringPtr(logsStorageAccountConnection)},
		{Name: to.StringPtr("AzureWebJobsStorage"), Value: to.StringPtr(logsStorageAccountConnection)},
		{Name: to.StringPtr("SumoLabsMetricEndpoint"), Value: to.StringPtr(sumoLogicSource)},
		{Name: to.StringPtr("AzureEventHubConnectionString"), Value: eventHubKey},
		{Name: to.StringPtr("StorageConnectionString"), Value: to.StringPtr(failedStorageAccountConnection)},
		{Name: to.StringPtr("WEBSITE_CONTENTAZUREFILECONNECTIONSTRING"), Value: to.StringPtr(logsStorageAccountConnection)},
		{Name: to.StringPtr("WEBSITE_CONTENTSHARE"), Value: to.StringPtr(logsStorageAccountName)},
	}
	return metricsAppSettings
}
