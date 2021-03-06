package api

type DashboardSyncDefinition struct {
	Type        string                      `json:"type"`
	Name        string                      `json:"name"`
	Description string                      `json:"description"`
	DetailLevel int                         `json:"detailLevel"`
	Properties  string                      `json:"properties"`
	Panels      []reportPanelSyncDefinition `json:"panels"`
	Filters     []filtersSyncDefinition     `json:"filters"`
}

type FolderSyncDefinition struct {
	Type        string                  `json:"type"`
	Name        string                  `json:"name"`
	Description string                  `json:"description"`
	Children    []contentSyncDefinition `json:"children"`
}

type GetContentResponse struct {
	CreatedAt   string   `json:"createdAt"`
	CreatedBy   string   `json:"createdBy"`
	ModifiedAt  string   `json:"modifiedAt"`
	ModifiedBy  string   `json:"modifiedBy"`
	Id          string   `json:"id"`
	Name        string   `json:"name"`
	ItemType    string   `json:"itemType"`
	ParentId    string   `json:"parentId"`
	Permissions []string `json:"permissions"`
}

type GetPathResponse struct {
	Path string `json:"path"`
}

type LookupTableSyncDefinition struct {
	Type            string   `json:"type"`
	Name            string   `json:"name"`
	Description     string   `json:"description"`
	Fields          []fields `json:"fields"`
	PrimaryKeys     []string `json:"primaryKeys"`
	TTL             int      `json:"ttl"`
	SizeLimitAction string   `json:"sizeLimitAction"`
}

type MetricsSavedSearchSyncDefinition struct {
	Type                      string                     `json:"type"`
	Name                      string                     `json:"name"`
	Description               string                     `json:"description"`
	TimeRange                 timeRangeDefinition        `json:"timeRange"`
	LogQuery                  string                     `json:"logQuery"`
	MetricsQueries            []metricsQueriesDefinition `json:"metricsQueries"`
	DesiredQuantizationInSecs int                        `json:"desiredQuantizationInSecs"`
	Properties                string                     `json:"properties"`
}

type MetricsSearchSyncDefinition struct {
	Type           string              `json:"type"`
	Name           string              `json:"name"`
	TimeRange      timeRangeDefinition `json:"timeRange"`
	Description    string              `json:"description"`
	Queries        []queries           `json:"queries"`
	VisualSettings string              `json:"visualSettings"`
}

type MewboardSyncDefinition struct {
	Type             string                  `json:"type"`
	Name             string                  `json:"name"`
	Description      string                  `json:"description"`
	Title            string                  `json:"title"`
	RootPanel        rootPanelDefinition     `json:"rootPanel"`
	Theme            string                  `json:"theme"`
	TopologyLabelMap topologyLabelMap        `json:"topologyLabelMap"`
	RefreshInterval  int                     `json:"refreshInterval"`
	TimeRange        timeRangeDefinition     `json:"timeRange"`
	Layout           layout                  `json:"layout"`
	Panels           panelsDefinition        `json:"panels"`
	Variables        variablesDefinition     `json:"variables"`
	ColoringRules    coloringRulesDefinition `json:"coloringRules"`
}

type MoveResponse struct {
	Id     string       `json:"id"`
	Errors []moveErrors `json:"errors"`
}

type ResponseType struct {
	Type string `json:"type"`
}

type SavedSearchWithScheduleSyncDefinition struct {
	Type           string         `json:"type"`
	Name           string         `json:"name"`
	Search         search         `json:"search"`
	SearchSchedule searchSchedule `json:"searchSchedule"`
	Description    string         `json:"description"`
}

type StartExportResponse struct {
	Id string `json:"id"`
}

type ExportStatusResponse struct {
	Status        string      `json:"status"`
	StatusMessage string      `json:"statusMessage,omitempty"`
	Error         exportError `json:"error,omitempty"`
}

type autoComplete struct {
	AutoCompleteType   string               `json:"autoCompleteType"`
	AutoCompleteKey    string               `json:"autoCompleteKey"`
	AutoCompleteValues []autoCompleteValues `json:"autoCompleteValues"`
	LookupFileName     string               `json:"lookupFileName"`
	LookupLabelColumn  string               `json:"lookupLabelColumn"`
	LookupValueColumn  string               `json:"lookupValueColumn"`
}

type autoCompleteValues struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

type autoParsingInfo struct {
	Mode string `json:"mode"`
}

type coloringRulesDefinition struct {
	Scope                           string          `json:"scope"`
	SingleSeriesAggregateFunction   string          `json:"singleSeriesAggregateFunction"`
	MultipleSeriesAggregateFunction string          `json:"multipleSeriesAggregateFunction"`
	ColorThresholds                 colorThresholds `json:"colorThresholds"`
}

type colorThresholds struct {
	Color string `json:"color"`
	Min   int    `json:"min"`
	Max   int    `json:"max"`
}

type contentSyncDefinition struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

type exportError struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	Detail  string `json:"detail,omitempty"`
}

type fields struct {
	FieldName string `json:"fieldName"`
	FieldType string `json:"fieldType"`
}

type filtersSyncDefinition struct {
	FieldName    string   `json:"fieldName"`
	Label        string   `json:"label"`
	DefaultValue string   `json:"defaultValue"`
	FilterType   string   `json:"filterType"`
	Properties   string   `json:"properties"`
	PanelIds     []string `json:"panelIds"`
}

type layout struct {
	LayoutType       string            `json:"layoutType"`
	LayoutStructures []layoutStructure `json:"layoutStructures"`
}

type layoutStructure struct {
	Key       string `json:"key"`
	Structure string `json:"structure"`
}

type metricsQueriesDefinition struct {
	Query string `json:"query"`
	RowId string `json:"rowId"`
}

type moveErrors struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type panelsDefinition struct {
	Id                                     string `json:"id"`
	Key                                    string `json:"key"`
	Title                                  string `json:"title"`
	visualSettings                         string `json:"visualSettings"`
	KeepVisualSettingsConsistentWithParent bool   `json:"keepVisualSettingsConsistentWithParent"`
	PanelType                              string `json:"panelType"`
}

type queries struct {
	QueryString string `json:"queryString"`
	QueryType   string `json:"queryType"`
	QueryKey    string `json:"queryKey"`
}

type queryParameters struct {
	Name         string       `json:"name"`
	Label        string       `json:"label"`
	Description  string       `json:"description"`
	DataType     string       `json:"dataType"`
	Value        string       `json:"value"`
	AutoComplete autoComplete `json:"autoComplete"`
}

type rootPanelDefinition struct {
	Id                                     string                    `json:"id"`
	Key                                    string                    `json:"key"`
	Title                                  string                    `json:"title"`
	VisualSettings                         string                    `json:"visualSettings"`
	KeepVisualSettingsConsistentWithParent bool                      `json:"keepVisualSettingsConsistentWithParent"`
	PanelType                              string                    `json:"panelType"`
	Layout                                 layout                    `json:"layout"`
	Panels                                 []panelsDefinition        `json:"panels"`
	Variables                              []variablesDefinition     `json:"variables"`
	ColoringRules                          []coloringRulesDefinition `json:"coloringRules"`
}

type reportPanelSyncDefinition struct {
	Name                      string                     `json:"name"`
	ViewerType                string                     `json:"viewerType"`
	DetailLevel               int                        `json:"detailLevel"`
	QueryString               string                     `json:"queryString"`
	MetricsQueries            []metricsQueriesDefinition `json:"metricsQueries"`
	TimeRange                 timeRangeDefinition        `json:"timeRange"`
	X                         int                        `json:"x"`
	Y                         int                        `json:"y"`
	Width                     int                        `json:"width"`
	Height                    int                        `json:"height"`
	Properties                string                     `json:"properties"`
	Id                        string                     `json:"id"`
	DesiredQuantizationInSecs int                        `json:"desiredQuantizationInSecs"`
	QueryParameters           []queryParameters          `json:"queryParameters"`
	AutoParsingInfo           autoParsingInfo            `json:"autoParsingInfo"`
}

type search struct {
	QueryText        string            `json:"queryText"`
	DefaultTimeRange string            `json:"defaultTimeRange"`
	ByReceiptTime    bool              `json:"byReceiptTime"`
	ViewName         string            `json:"viewName"`
	ViewStartTime    string            `json:"viewStartTime"`
	QueryParameters  []queryParameters `json:"queryParameters"`
	ParsingMode      string            `json:"parsingMode"`
}

type searchSchedule struct {
	CronExpression       string                     `json:"cronExpression"`
	DisplayableTimeRange string                     `json:"displayableTimeRange"`
	ParseableTimeRange   timeRangeDefinition        `json:"parseableTimeRange"`
	TimeZone             string                     `json:"timeZone"`
	Threshold            searchScheduleThreshold    `json:"threshold"`
	Notification         searchScheduleNotification `json:"notification"`
	ScheduleType         string                     `json:"scheduleType"`
	MuteErrorEmails      bool                       `json:"muteErrorEmails"`
	Parameters           []searchScheduleParamters  `json:"parameters"`
}

type searchScheduleNotification struct {
	TaskType string `json:"taskType"`
}

type searchScheduleParamters struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type searchScheduleThreshold struct {
	ThresholdType string `json:"thresholdType"`
	Operator      string `json:"operator"`
	Count         int    `json:"count"`
}

type timeRangeDefinition struct {
	Type string                  `json:"type"`
	From timeRangeFromDefinition `json:"from"`
}

type timeRangeFromDefinition struct {
	Type         string `json:"type"`
	RelativeTime string `json:"relativeTime"`
}

type topologyLabelMap struct {
	Service []string `json:"service"`
}

type variablesDefinition struct {
	Id               string                    `json:"id"`
	Name             string                    `json:"name"`
	DisplayName      string                    `json:"displayName"`
	DefaultValue     string                    `json:"defaultValue"`
	SourceDefinition variablesSourceDefinition `json:"sourceDefinition"`
	AllowMultiSelect bool                      `json:"allowMultiSelect"`
	IncludeAllOption bool                      `json:"includeAllOption"`
	HideFromUI       bool                      `json:"hideFromUI"`
}

type variablesSourceDefinition struct {
	VariableSourceType string `json:"variableSourceType"`
}
