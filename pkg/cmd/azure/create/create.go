package create

import (
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/profiles/latest/resources/mgmt/features"
	"github.com/Azure/azure-sdk-for-go/profiles/latest/storage/mgmt/storage"
	"github.com/Azure/azure-sdk-for-go/services/appinsights/mgmt/2015-05-01/insights"
	"github.com/Azure/azure-sdk-for-go/services/eventgrid/mgmt/2020-06-01/eventgrid"
	"github.com/Azure/azure-sdk-for-go/services/eventhub/mgmt/2017-04-01/eventhub"
	"github.com/Azure/azure-sdk-for-go/services/servicebus/mgmt/2017-04-01/servicebus"
	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2020-06-01/web"
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/wizedkyle/sumocli/internal/az"
	"github.com/wizedkyle/sumocli/internal/clients"
	"github.com/wizedkyle/sumocli/internal/config"
	"github.com/wizedkyle/sumocli/pkg/cmd/collectors/create"
	sources "github.com/wizedkyle/sumocli/pkg/cmd/sources/create"
	"github.com/wizedkyle/sumocli/pkg/logging"
)

func NewCmdAzureCreate() *cobra.Command {
	var (
		category   string
		prefix     string
		diagnostic bool
		metrics    bool
		blob       bool
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create Azure infrastructure to collect logs or metrics",
		Run: func(cmd *cobra.Command, args []string) {
			logger := logging.GetConsoleLogger()
			log := logging.GetConsoleLogger()
			logger.Debug().Msg("Create Azure infrastructure request started")
			if blob == true {
				azureCreateBlobCollection(category, prefix, log)
			} else if metrics == true {
				azureCreateMetricCollection(category, prefix, log)
			} else if diagnostic == true {
				azureCreateDiagLogCollection(category, prefix, log)
			} else {
				fmt.Println("Please select either --diagnostic, --logs or --metrics")
			}
			logger.Debug().Msg("Create Azure infrastructure request finished.")
		},
	}

	cmd.Flags().BoolVar(&blob, "blob", false, "Deploys infrastructure for Azure Blob collection.")
	cmd.Flags().BoolVar(&diagnostic, "diagnostic", false, "Deploys infrastructure for Azure Diagnostic Log collection")
	cmd.Flags().BoolVar(&metrics, "metrics", false, "Deploys infrastructure for Azure Metrics collection")
	cmd.Flags().StringVar(&category, "category", "", "Specify the source category for the Sumo Logic source")
	cmd.Flags().StringVar(&prefix, "prefix", "", "User defined string that is used to name Azure resources can only contain numbers, lowercase letters and a max length of 10")
	_ = cmd.MarkFlagRequired("category")
	_ = cmd.MarkFlagRequired("prefix")
	return cmd
}

func azureCreateBlobCollection(category string, prefix string, log zerolog.Logger) {
	ctx := context.Background()
	logsName := "scliblob"
	rgName := logsName + prefix
	sgName := logsName + prefix
	sourceSgName := "sclisrc" + prefix
	nsName := logsName + prefix
	nsAuthName := logsName + prefix
	queueName := logsName + prefix
	ehNsName := logsName + prefix + "ehns"
	ehName := logsName + prefix + "eh"
	ehAuthName := logsName + prefix + "ehrule"
	cgName := logsName + prefix + "cg"
	eventSubName := logsName + prefix + "sub"
	insightsName := logsName + prefix
	appPlanName := logsName + prefix
	functionName := logsName + prefix
	collectorName := logsName + prefix
	sourceName := logsName + prefix
	appRepoUrl := "https://github.com/SumoLogic/sumologic-azure-function"
	branch := "master"

	createResourceGroup(ctx, rgName, log)
	_, _ = createStorageAccount(ctx, rgName, sgName, log)
	sourceSgAcc, _ := createStorageAccount(ctx, rgName, sourceSgName, log)
	createStorageAccountTable(ctx, rgName, sgName, log)
	_, _ = createServiceBusNamespace(ctx, rgName, nsName, log)
	createServiceBusAuthRule(ctx, rgName, sgName, nsAuthName, log)
	sbKey := getServiceBusConnectionString(ctx, rgName, nsName, nsAuthName, log)
	createServiceBusQueue(ctx, rgName, nsName, queueName, log)
	_, _ = createEventHubNamespace(ctx, rgName, ehNsName, log)
	eh := createEventHub(ctx, rgName, ehNsName, ehName, log)
	createEventHubAuthRule(ctx, rgName, ehNsName, ehName, ehAuthName, log)
	ehKey := getEventHubConnectionString(ctx, rgName, ehNsName, ehName, ehAuthName, log)
	createEventHubConsumerGroup(ctx, rgName, ehNsName, ehName, cgName, log)
	createEventGridSubscription(ctx, sourceSgAcc, eventSubName, eh, log)
	appInsights := createApplicationInsight(ctx, rgName, insightsName, log)
	appServicePlan, _ := createAppServicePlan(ctx, rgName, appPlanName, log)

	// Creates a Sumo Logic collector and HTTP source
	collector := create.Collector(collectorName, "", "", "")
	source := sources.HTTPSource(category, nil, false, false, sourceName, collector.Collector.Id, log)

	// Creates each function app, adds source control integration and provides custom App Settings
	// Blob collection requires three apps:  blob reader, consumer, dlq (dead letter queue)
	readerAppSettings := az.ReaderAppSettings(
		sgName,
		getStorageAccountConnectionString(ctx, rgName, sgName, log),
		appInsights.InstrumentationKey,
		ehKey.PrimaryConnectionString,
		sbKey.PrimaryConnectionString)
	readerFunctionName := functionName + "reader"
	createFunctionApp(ctx, rgName, readerFunctionName, appServicePlan, readerAppSettings, log)
	createFunctionAppSourceControl(ctx, rgName, readerFunctionName, appRepoUrl, branch, log)

	consumerAppSettings := az.ConsumerAppSettings(
		sgName,
		getStorageAccountConnectionString(ctx, rgName, sgName, log),
		appInsights.InstrumentationKey,
		sbKey.PrimaryConnectionString,
		source.Source.Url)
	consumerFunctionName := functionName + "consumer"
	createFunctionApp(ctx, rgName, consumerFunctionName, appServicePlan, consumerAppSettings, log)
	createFunctionAppSourceControl(ctx, rgName, consumerFunctionName, appRepoUrl, branch, log)

	dlqAppSettings := az.DlqAppSettings(
		sgName,
		getStorageAccountConnectionString(ctx, rgName, sgName, log),
		appInsights.InstrumentationKey,
		sbKey.PrimaryConnectionString,
		queueName,
		source.Source.Url)
	dlqFunctionName := functionName + "dlq"
	createFunctionApp(ctx, rgName, dlqFunctionName, appServicePlan, dlqAppSettings, log)
	createFunctionAppSourceControl(ctx, rgName, dlqFunctionName, appRepoUrl, branch, log)
}

func azureCreateDiagLogCollection(category string, prefix string, log zerolog.Logger) {
	ctx := context.Background()
	logsName := "sclidiag"
	rgName := logsName + prefix
	sgLogsName := logsName + prefix + "logs"
	sgFailedName := logsName + prefix + "failed"
	ehNsName := logsName + prefix + "ehns"
	ehName := logsName + prefix + "eh"
	ehAuthName := logsName + prefix + "ehrule"
	cgName := logsName + prefix + "cg"
	appPlanName := logsName + prefix
	functionName := logsName + prefix
	collectorName := logsName + prefix
	sourceName := logsName + prefix
	appRepoUrl := "https://github.com/SumoLogic/sumologic-azure-function"
	branch := "master"

	createResourceGroup(ctx, rgName, log)
	_, _ = createStorageAccount(ctx, rgName, sgLogsName, log)
	_, _ = createStorageAccount(ctx, rgName, sgFailedName, log)
	_, _ = createEventHubNamespace(ctx, rgName, ehNsName, log)
	createEventHub(ctx, rgName, ehNsName, ehName, log)
	createEventHubAuthRule(ctx, rgName, ehNsName, ehName, ehAuthName, log)
	ehKey := getEventHubConnectionString(ctx, rgName, ehNsName, ehName, ehAuthName, log)
	createEventHubConsumerGroup(ctx, rgName, ehNsName, ehName, cgName, log)
	appServicePlan, _ := createAppServicePlan(ctx, rgName, appPlanName, log)

	// Creates a Sumo Logic collector and HTTP source
	collector := create.Collector(collectorName, "", "", "")
	source := sources.HTTPSource(category, nil, false, false, sourceName, collector.Collector.Id, log)
	diagnosticAppSettings := az.DiagnosticLogsAppSettings(
		sgLogsName,
		getStorageAccountConnectionString(ctx, rgName, sgLogsName, log),
		getStorageAccountConnectionString(ctx, rgName, sgFailedName, log),
		ehKey.PrimaryConnectionString,
		source.Source.Url)
	diagnosticFunctionName := functionName + "diagnostic"
	createFunctionApp(ctx, rgName, diagnosticFunctionName, appServicePlan, diagnosticAppSettings, log)
	createFunctionAppSourceControl(ctx, rgName, diagnosticFunctionName, appRepoUrl, branch, log)
}

func azureCreateMetricCollection(category string, prefix string, log zerolog.Logger) {
	ctx := context.Background()
	metricsName := "sclimetrics"
	rgName := metricsName + prefix
	sgLogsName := metricsName + prefix + "logs"
	sgFailedName := metricsName + prefix + "failed"
	ehNsName := metricsName + prefix + "ehns"
	ehName := metricsName + prefix + "eh"
	ehAuthName := metricsName + prefix + "ehrule"
	cgName := metricsName + prefix + "cg"
	appPlanName := metricsName + prefix
	functionName := metricsName + prefix
	collectorName := metricsName + prefix
	sourceName := metricsName + prefix
	appRepoUrl := "https://github.com/SumoLogic/sumologic-azure-function"
	branch := "master"

	createResourceGroup(ctx, rgName, log)
	_, _ = createStorageAccount(ctx, rgName, sgLogsName, log)
	_, _ = createStorageAccount(ctx, rgName, sgFailedName, log)
	_, _ = createEventHubNamespace(ctx, rgName, ehNsName, log)
	createEventHub(ctx, rgName, ehNsName, ehName, log)
	createEventHubAuthRule(ctx, rgName, ehNsName, ehName, ehAuthName, log)
	ehKey := getEventHubConnectionString(ctx, rgName, ehNsName, ehName, ehAuthName, log)
	createEventHubConsumerGroup(ctx, rgName, ehNsName, ehName, cgName, log)
	appServicePlan, _ := createAppServicePlan(ctx, rgName, appPlanName, log)

	// Creates a Sumo Logic collector and HTTP source
	collector := create.Collector(collectorName, "", "", "")
	source := sources.HTTPSource(category, nil, false, false, sourceName, collector.Collector.Id, log)
	metricsAppSettings := az.MetricsAppSettings(
		sgLogsName,
		getStorageAccountConnectionString(ctx, rgName, sgLogsName, log),
		getStorageAccountConnectionString(ctx, rgName, sgFailedName, log),
		ehKey.PrimaryConnectionString,
		source.Source.Url)
	metricsFunctionName := functionName + "metrics"
	createFunctionApp(ctx, rgName, metricsFunctionName, appServicePlan, metricsAppSettings, log)
	createFunctionAppSourceControl(ctx, rgName, metricsFunctionName, appRepoUrl, branch, log)
}

func createApplicationInsight(ctx context.Context, rgName string, insightsName string, log zerolog.Logger) insights.ApplicationInsightsComponent {
	log.Info().Msg("creating or updating application appInsights: " + insightsName)
	insightsClient := clients.GetInsightsClient()
	appInsights, err := insightsClient.CreateOrUpdate(
		ctx,
		rgName,
		insightsName,
		insights.ApplicationInsightsComponent{
			Kind: to.StringPtr("web"),
			ApplicationInsightsComponentProperties: &insights.ApplicationInsightsComponentProperties{
				ApplicationID:              nil,
				AppID:                      nil,
				ApplicationType:            "",
				FlowType:                   "",
				RequestSource:              "",
				InstrumentationKey:         nil,
				CreationDate:               nil,
				TenantID:                   nil,
				HockeyAppID:                nil,
				HockeyAppToken:             nil,
				ProvisioningState:          nil,
				SamplingPercentage:         nil,
				ConnectionString:           nil,
				RetentionInDays:            nil,
				DisableIPMasking:           nil,
				ImmediatePurgeDataOn30Days: nil,
				PrivateLinkScopedResources: nil,
				IngestionMode:              "",
			},
			Location: to.StringPtr(config.GetDefaultLocation()),
			Tags:     config.GetAzureLogTags(),
		})

	if err != nil {
		log.Fatal().Err(err).Msg("cannot create or update application appInsights: " + insightsName)
	}

	log.Info().Msg("created or updated application appInsights: " + insightsName)
	return appInsights
}

func createAppServicePlan(ctx context.Context, rgName string, appPlanName string, log zerolog.Logger) (web.AppServicePlan, error) {
	log.Info().Msg("creating or updating app service plan " + appPlanName)
	appClient := clients.GetAppServicePlanClient()
	appPlan, err := appClient.CreateOrUpdate(
		ctx,
		rgName,
		appPlanName,
		web.AppServicePlan{
			AppServicePlanProperties: nil,
			Sku: &web.SkuDescription{
				Name: to.StringPtr("Y1"),
				Tier: to.StringPtr("Dynamic"),
				Size: to.StringPtr("Y1"),
			},
			Kind:     to.StringPtr("FunctionApp"),
			Location: to.StringPtr(config.GetDefaultLocation()),
			Tags:     config.GetAzureLogTags(),
		})

	if err != nil {
		log.Fatal().Err(err).Msg("cannot create or update app service plan " + appPlanName)
	}

	err = appPlan.WaitForCompletionRef(ctx, appClient.Client)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create or update app service plan " + appPlanName)
	}

	log.Info().Msg("created or updated app service plan " + appPlanName)
	return appPlan.Result(appClient)
}

func createEventGridSubscription(ctx context.Context, scope storage.Account, eventSubName string, eventhub eventhub.Model, log zerolog.Logger) eventgrid.EventSubscriptionsCreateOrUpdateFuture {
	log.Info().Msg("creating or updating event grid subscription " + eventSubName)
	egSubClient := clients.GetEventGridSubscriptionClient()
	subscription, err := egSubClient.CreateOrUpdate(
		ctx,
		to.String(scope.ID),
		eventSubName,
		eventgrid.EventSubscription{
			EventSubscriptionProperties: &eventgrid.EventSubscriptionProperties{
				Destination: eventgrid.EventHubEventSubscriptionDestination{
					EventHubEventSubscriptionDestinationProperties: &eventgrid.EventHubEventSubscriptionDestinationProperties{
						ResourceID: eventhub.ID,
					},
					EndpointType: eventgrid.EndpointTypeEventHub,
				},
				Filter: &eventgrid.EventSubscriptionFilter{
					IncludedEventTypes: &[]string{
						"Microsoft.Storage.BlobCreated",
					},
				},
			},
		})

	if err != nil {
		log.Fatal().Err(err).Msg("cannot create or update event grid subscription " + eventSubName)
	}
	err = subscription.WaitForCompletionRef(ctx, egSubClient.Client)

	if err != nil {
		log.Fatal().Err(err).Msg("cannot create or update event grid subscription " + eventSubName)
	}

	log.Info().Msg("created or updated event grid subscription " + eventSubName)
	return subscription
}

func createEventHubNamespace(ctx context.Context, rgName string, nsName string, log zerolog.Logger) (eventhub.EHNamespace, error) {
	log.Info().Msg("creating or updating event hub namespace " + nsName)
	ehClient := clients.GetEventHubNamespaceClient()
	ehNamespace, err := ehClient.CreateOrUpdate(
		ctx,
		rgName,
		nsName,
		eventhub.EHNamespace{
			Sku: &eventhub.Sku{
				Name:     eventhub.Standard,
				Capacity: to.Int32Ptr(1),
			},
			Location: to.StringPtr(config.GetDefaultLocation()),
			Tags:     config.GetAzureLogTags(),
		})

	if err != nil {
		log.Fatal().Err(err).Msg("cannot create or update event hub namespace " + nsName)
	}

	err = ehNamespace.WaitForCompletionRef(ctx, ehClient.Client)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create or update event hub namespace " + nsName)
	}

	log.Info().Msg("created or updated event hub namespace " + nsName)
	return ehNamespace.Result(ehClient)
}

func createEventHub(ctx context.Context, rgName string, ehNsName string, ehName string, log zerolog.Logger) eventhub.Model {
	log.Info().Msg("creating or updating event hub " + ehName)
	ehClient := clients.GetEventHubClient()
	eh, err := ehClient.CreateOrUpdate(
		ctx,
		rgName,
		ehNsName,
		ehName,
		eventhub.Model{
			Properties: &eventhub.Properties{
				MessageRetentionInDays: to.Int64Ptr(7),
				PartitionCount:         to.Int64Ptr(2),
			},
		})

	if err != nil {
		log.Fatal().Err(err).Msg("cannot create or update event hub " + ehName)
	}
	return eh
}

func createEventHubAuthRule(ctx context.Context, rgName string, ehNsName string, ehName string, ehAuthName string, log zerolog.Logger) eventhub.AuthorizationRule {
	log.Info().Msg("creating or updating event hub authorization rule " + ehAuthName)
	ehClient := clients.GetEventHubClient()
	ehAuthRule, err := ehClient.CreateOrUpdateAuthorizationRule(
		ctx,
		rgName,
		ehNsName,
		ehName,
		ehAuthName,
		eventhub.AuthorizationRule{
			AuthorizationRuleProperties: &eventhub.AuthorizationRuleProperties{
				Rights: &[]eventhub.AccessRights{
					"Listen",
					"Manage",
					"Send",
				}},
		})

	if err != nil {
		log.Fatal().Err(err).Msg("cannot create or update event hub authorization rule " + ehAuthName)
	}

	log.Info().Msg("created or updated event hub authorization rule " + ehAuthName)
	return ehAuthRule
}

func createEventHubConsumerGroup(ctx context.Context, rgName string, ehNsName string, ehName string, cgName string, log zerolog.Logger) {
	log.Info().Msg("creating or updating event hub consumer group " + cgName)
	csClient := clients.GetConsumerGroupsClient()
	_, err := csClient.CreateOrUpdate(
		ctx,
		rgName,
		ehNsName,
		ehName,
		cgName,
		eventhub.ConsumerGroup{
			ConsumerGroupProperties: nil,
		})

	if err != nil {
		log.Fatal().Err(err).Msg("cannot create or update event hub consumer group " + cgName)
	}
	log.Info().Msg("created or updated event hub consumer group " + cgName)
}

func getEventHubConnectionString(ctx context.Context, rgName string, ehNsName string, ehName string, ehAuthName string, log zerolog.Logger) eventhub.AccessKeys {
	log.Info().Msg("getting event hub keys for " + ehAuthName)
	ehClient := clients.GetEventHubClient()
	ehKey, err := ehClient.ListKeys(
		ctx,
		rgName,
		ehNsName,
		ehName,
		ehAuthName)

	if err != nil {
		log.Fatal().Err(err).Msg("cannot get event hub keys for " + ehAuthName)
	}

	log.Info().Msg("obtained event hub keys for " + ehAuthName)
	return ehKey
}

func createFunctionApp(ctx context.Context, rgName string, functionName string, appSerivceId web.AppServicePlan, appSettings []web.NameValuePair, log zerolog.Logger) web.AppsCreateOrUpdateFuture {
	log.Info().Msg("creating or updating azure function " + functionName)
	appClient := clients.GetAppServiceClient()
	functionApp, err := appClient.CreateOrUpdate(
		ctx,
		rgName,
		functionName,
		web.Site{
			SiteProperties: &web.SiteProperties{
				Enabled:      to.BoolPtr(true),
				ServerFarmID: appSerivceId.ID,
				SiteConfig: &web.SiteConfig{
					AppSettings: &appSettings,
					ScmType:     web.ScmTypeNone,
				},
				ClientAffinityEnabled: to.BoolPtr(true),
				DailyMemoryTimeQuota:  to.Int32Ptr(1000),
				HTTPSOnly:             to.BoolPtr(true),
			},
			Identity: &web.ManagedServiceIdentity{
				Type: web.ManagedServiceIdentityTypeSystemAssigned,
			},
			Kind:     to.StringPtr("FunctionApp"),
			Location: to.StringPtr(config.GetDefaultLocation()),
			Tags:     config.GetAzureLogTags(),
		})

	if err != nil {
		log.Fatal().Err(err).Msg("cannot create azure function app " + functionName)
	}

	err = functionApp.WaitForCompletionRef(ctx, appClient.Client)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create azure function app " + functionName)
	}

	log.Info().Msg("created or updated azure function app " + functionName)
	return functionApp
}

func createFunctionAppSourceControl(ctx context.Context, rgName string, functionName string, appRepoUrl string, branch string, log zerolog.Logger) web.AppsCreateOrUpdateSourceControlFuture {
	log.Info().Msg("creating or updating source control for function app " + functionName)
	appClient := clients.GetAppServiceClient()
	functionAppSc, err := appClient.CreateOrUpdateSourceControl(
		ctx,
		rgName,
		functionName,
		web.SiteSourceControl{
			SiteSourceControlProperties: &web.SiteSourceControlProperties{
				RepoURL:             to.StringPtr(appRepoUrl),
				Branch:              to.StringPtr(branch),
				IsManualIntegration: to.BoolPtr(true),
			},
		})

	if err != nil {
		log.Fatal().Err(err).Msg("cannot create source control settings on function app " + functionName)
	}

	err = functionAppSc.WaitForCompletionRef(ctx, appClient.Client)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create source control settings on function app " + functionName)
	}

	log.Info().Msg("created or updated source control settings on function app " + functionName)
	return functionAppSc
}

func createResourceGroup(ctx context.Context, rgName string, log zerolog.Logger) features.ResourceGroup {
	rgClient := clients.GetResourceGroupClient()
	log.Info().Msg("creating or updating resource group " + rgName)
	rg, err := rgClient.CreateOrUpdate(
		ctx,
		rgName,
		features.ResourceGroup{
			Name:     to.StringPtr(rgName),
			Location: to.StringPtr(config.GetDefaultLocation()),
			Tags:     config.GetAzureLogTags(),
		})

	if err != nil {
		log.Fatal().Err(err).Msg("cannot create or update resource group " + rgName)
	}
	log.Info().Msg("created or updated resource group " + rgName)
	return rg
}

func createStorageAccount(ctx context.Context, rgName string, sgName string, log zerolog.Logger) (storage.Account, error) {
	log.Info().Msg("creating or updating storage account " + sgName)
	sgClient := clients.GetStorageClient()

	// TODO: add storage account name check
	sgAccount, err := sgClient.Create(
		ctx,
		rgName,
		sgName,
		storage.AccountCreateParameters{
			Sku: &storage.Sku{
				Name: storage.StandardLRS,
				Tier: storage.Standard,
			},
			Kind:     storage.StorageV2,
			Location: to.StringPtr(config.GetDefaultLocation()),
			Tags:     config.GetAzureLogTags(),
		})

	if err != nil {
		log.Fatal().Err(err).Msg("cannot create or update storage account " + sgName)
	}

	err = sgAccount.WaitForCompletionRef(ctx, sgClient.Client)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create or update storage account " + sgName)
	}

	log.Info().Msg("created or updated storage account " + rgName)
	return sgAccount.Result(sgClient)
}

func createStorageAccountTable(ctx context.Context, rgName string, sgName string, log zerolog.Logger) {
	log.Info().Msg("creating FileOffsetMap table")
	tableClient := clients.GetStorageTableClient()
	_, err := tableClient.Create(
		ctx,
		rgName,
		sgName,
		"FileOffsetMap")

	if err != nil {
		log.Fatal().Err(err).Msg("cannot create FileOffsetMap table")
	}

	log.Info().Msg("created FileOffsetMap table")
}

func getStorageAccountConnectionString(ctx context.Context, rgName string, sgName string, log zerolog.Logger) string {
	log.Info().Msg("getting storage account connection string for " + sgName)
	sgClient := clients.GetStorageClient()
	sgKey, err := sgClient.ListKeys(
		ctx,
		rgName,
		sgName,
		storage.Kerb)

	if err != nil {
		log.Fatal().Err(err).Msg("cannot get storage account keys")
	}

	log.Info().Msg("connection string obtained for storage account " + sgName)
	return fmt.Sprintf("DefaultEndpointsProtocol=https;AccountName=%s;AccountKey=%s;EndpointSuffix=core.windows.net", sgName, to.String((*sgKey.Keys)[0].Value))
}

func createServiceBusNamespace(ctx context.Context, rgName string, nsName string, log zerolog.Logger) (servicebus.SBNamespace, error) {
	log.Info().Msg("creating or updating service bus namespace " + nsName)
	nsClient := clients.GetNamespaceClient()
	ns, err := nsClient.CreateOrUpdate(
		ctx,
		rgName,
		nsName,
		servicebus.SBNamespace{
			Sku: &servicebus.SBSku{
				Name: servicebus.Standard,
			},
			Location: to.StringPtr(config.GetDefaultLocation()),
			Tags:     config.GetAzureLogTags(),
		})

	if err != nil {
		log.Fatal().Err(err).Msg("cannot create service bus namespace " + nsName)
	}

	err = ns.WaitForCompletionRef(ctx, nsClient.Client)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create service bus namespace " + nsName)
	}

	log.Info().Msg("created or updated service bus namespace " + nsName)
	return ns.Result(nsClient)
}

func createServiceBusAuthRule(ctx context.Context, rgName string, nsName string, nsAuthName string, log zerolog.Logger) servicebus.SBAuthorizationRule {
	log.Info().Msg("creating or updating service bus namespace authorization rule " + nsAuthName)
	nsClient := clients.GetNamespaceClient()
	sbAuthRule, err := nsClient.CreateOrUpdateAuthorizationRule(
		ctx,
		rgName,
		nsName,
		nsAuthName,
		servicebus.SBAuthorizationRule{
			SBAuthorizationRuleProperties: &servicebus.SBAuthorizationRuleProperties{
				Rights: &[]servicebus.AccessRights{
					"Listen",
					"Manage",
					"Send",
				},
			},
		})

	if err != nil {
		log.Fatal().Err(err).Msg("cannot create service bus namespace authorization rule " + nsAuthName)
	}

	log.Info().Msg("created or updated service bus namespace authorization rule " + nsAuthName)
	return sbAuthRule
}

func getServiceBusConnectionString(ctx context.Context, rgName string, nsName string, nsAuthName string, log zerolog.Logger) servicebus.AccessKeys {
	log.Info().Msg("getting service bus connection string for " + nsAuthName)
	nsClient := clients.GetNamespaceClient()
	sbKeys, err := nsClient.ListKeys(
		ctx,
		rgName,
		nsName,
		nsAuthName)

	if err != nil {
		log.Fatal().Err(err).Msg("cannot get keys for service bus " + nsAuthName)
	}

	log.Info().Msg("obtained service bus connection string for " + nsAuthName)
	return sbKeys
}

func createServiceBusQueue(ctx context.Context, rgName string, nsName string, queueName string, log zerolog.Logger) {
	log.Info().Msg("creating or updating service bus queue " + queueName)
	queueClient := clients.GetQueueClient()
	_, err := queueClient.CreateOrUpdate(
		ctx,
		rgName,
		nsName,
		queueName,
		servicebus.SBQueue{
			SBQueueProperties: &servicebus.SBQueueProperties{
				LockDuration:                        to.StringPtr("PT5M"),
				MaxSizeInMegabytes:                  to.Int32Ptr(2048),
				RequiresDuplicateDetection:          to.BoolPtr(false),
				RequiresSession:                     to.BoolPtr(false),
				DefaultMessageTimeToLive:            to.StringPtr("P14D"),
				DeadLetteringOnMessageExpiration:    to.BoolPtr(true),
				DuplicateDetectionHistoryTimeWindow: to.StringPtr("PT10M"),
				MaxDeliveryCount:                    to.Int32Ptr(10),
				EnableBatchedOperations:             to.BoolPtr(true),
				AutoDeleteOnIdle:                    to.StringPtr("P10675199DT2H48M5.4775807S"),
				EnablePartitioning:                  to.BoolPtr(true),
				EnableExpress:                       to.BoolPtr(true),
			},
		})

	if err != nil {
		log.Fatal().Err(err).Msg("cannot create service bus queue " + queueName)
	}
	log.Info().Msg("created or updated service bus queue " + queueName)
}
