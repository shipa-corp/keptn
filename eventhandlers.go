package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/brunoa19/shipa-keptn/shipa"
	cloudevents "github.com/cloudevents/sdk-go/v2" // make sure to use v2 cloudevents here
	keptn "github.com/keptn/go-utils/pkg/lib"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"log"
)

/**
* Here are all the handler functions for the individual event
* See https://github.com/keptn/spec/blob/0.8.0-alpha/cloudevents.md for details on the payload
**/

// GenericLogKeptnCloudEventHandler is a generic handler for Keptn Cloud Events that logs the CloudEvent
func GenericLogKeptnCloudEventHandler(myKeptn *keptnv2.Keptn, incomingEvent cloudevents.Event, data interface{}) error {
	log.Printf("Handling %s Event: %s", incomingEvent.Type(), incomingEvent.Context.GetID())
	log.Printf("CloudEvent %T: %v", data, data)

	return nil
}

// OldHandleConfigureMonitoringEvent handles old configure-monitoring events
// TODO: add in your handler code
func OldHandleConfigureMonitoringEvent(myKeptn *keptnv2.Keptn, incomingEvent cloudevents.Event, data *keptn.ConfigureMonitoringEventData) error {
	log.Printf("Handling old configure-monitoring Event: %s", incomingEvent.Context.GetID())

	return nil
}

// HandleConfigureMonitoringTriggeredEvent handles configure-monitoring.triggered events
// TODO: add in your handler code
func HandleConfigureMonitoringTriggeredEvent(myKeptn *keptnv2.Keptn, incomingEvent cloudevents.Event, data *keptnv2.ConfigureMonitoringTriggeredEventData) error {
	log.Printf("Handling configure-monitoring.triggered Event: %s", incomingEvent.Context.GetID())

	return nil
}

// HandleDeploymentTriggeredEvent handles deployment.triggered events
// TODO: add in your handler code
func HandleDeploymentTriggeredEvent(myKeptn *keptnv2.Keptn, incomingEvent cloudevents.Event, data *keptnv2.DeploymentTriggeredEventData) error {
	log.Printf("Handling deployment.triggered Event: %s", incomingEvent.Context.GetID())

	return nil
}

// HandleTestTriggeredEvent handles test.triggered events
// TODO: add in your handler code
func HandleTestTriggeredEvent(myKeptn *keptnv2.Keptn, incomingEvent cloudevents.Event, data *keptnv2.TestTriggeredEventData) error {
	log.Printf("Handling test.triggered Event: %s", incomingEvent.Context.GetID())

	return nil
}

// HandleApprovalTriggeredEvent handles approval.triggered events
// TODO: add in your handler code
func HandleApprovalTriggeredEvent(myKeptn *keptnv2.Keptn, incomingEvent cloudevents.Event, data *keptnv2.ApprovalTriggeredEventData) error {
	log.Printf("Handling approval.triggered Event: %s", incomingEvent.Context.GetID())

	return nil
}

// HandleEvaluationTriggeredEvent handles evaluation.triggered events
// TODO: add in your handler code
func HandleEvaluationTriggeredEvent(myKeptn *keptnv2.Keptn, incomingEvent cloudevents.Event, data *keptnv2.EvaluationTriggeredEventData) error {
	log.Printf("Handling evaluation.triggered Event: %s", incomingEvent.Context.GetID())

	return nil
}

// HandleReleaseTriggeredEvent handles release.triggered events
// TODO: add in your handler code
func HandleReleaseTriggeredEvent(myKeptn *keptnv2.Keptn, incomingEvent cloudevents.Event, data *keptnv2.ReleaseTriggeredEventData) error {
	log.Printf("Handling release.triggered Event: %s", incomingEvent.Context.GetID())

	return nil
}

// HandleGetSliTriggeredEvent handles get-sli.triggered events if SLIProvider == shipa-keptn
// This function acts as an example showing how to handle get-sli events by sending .started and .finished events
// TODO: adapt handler code to your needs
func HandleGetSliTriggeredEvent(myKeptn *keptnv2.Keptn, incomingEvent cloudevents.Event, data *keptnv2.GetSLITriggeredEventData) error {
	log.Printf("Handling get-sli.triggered Event: %s", incomingEvent.Context.GetID())

	// Step 1 - Do we need to do something?
	// Lets make sure we are only processing an event that really belongs to our SLI Provider
	if data.GetSLI.SLIProvider != "shipa-keptn" {
		log.Printf("Not handling get-sli event as it is meant for %s", data.GetSLI.SLIProvider)
		return nil
	}

	// Step 2 - Send out a get-sli.started CloudEvent
	// The get-sli.started cloud-event is new since Keptn 0.8.0 and is required to be send when the task is started
	_, err := myKeptn.SendTaskStartedEvent(data, ServiceName)

	if err != nil {
		errMsg := fmt.Sprintf("Failed to send task started CloudEvent (%s), aborting...", err.Error())
		log.Println(errMsg)
		return err
	}

	// Step 4 - prep-work
	// Get any additional input / configuration data
	// - Labels: get the incoming labels for potential config data and use it to pass more labels on result, e.g: links
	// - SLI.yaml: if your service uses SLI.yaml to store query definitions for SLIs get that file from Keptn
	labels := data.Labels
	if labels == nil {
		labels = make(map[string]string)
	}
	testRunID := labels["testRunId"]

	// Step 5 - get SLI Config File
	// Get SLI File from shipa-keptn subdirectory of the config repo - to add the file use:
	//   keptn add-resource --project=PROJECT --stage=STAGE --service=SERVICE --resource=my-sli-config.yaml  --resourceUri=shipa-keptn/sli.yaml
	sliFile := "shipa-keptn/sli.yaml"
	sliConfigFileContent, err := myKeptn.GetKeptnResource(sliFile)

	// FYI you do not need to "fail" if sli.yaml is missing, you can also assume smart defaults like we do
	// in keptn-contrib/dynatrace-service and keptn-contrib/prometheus-service
	if err != nil {
		// failed to fetch sli config file
		errMsg := fmt.Sprintf("Failed to fetch SLI file %s from config repo: %s", sliFile, err.Error())
		log.Println(errMsg)
		// send a get-sli.finished event with status=error and result=failed back to Keptn

		_, err = myKeptn.SendTaskFinishedEvent(&keptnv2.EventData{
			Status: keptnv2.StatusErrored,
			Result: keptnv2.ResultFailed,
		}, ServiceName)

		return err
	}

	fmt.Println(sliConfigFileContent)

	// Step 6 - do your work - iterate through the list of requested indicators and return their values
	// Indicators: this is the list of indicators as requested in the SLO.yaml
	// SLIResult: this is the array that will receive the results
	indicators := data.GetSLI.Indicators
	sliResults := []*keptnv2.SLIResult{}

	for _, indicatorName := range indicators {
		sliResult := &keptnv2.SLIResult{
			Metric: indicatorName,
			Value:  123.4, // ToDo: Fetch the values from your monitoring tool here
		}
		sliResults = append(sliResults, sliResult)
	}

	// Step 7 - add additional context via labels (e.g., a backlink to the monitoring or CI tool)
	labels["Link to Data Source"] = "https://mydatasource/myquery?testRun=" + testRunID

	// Step 8 - Build get-sli.finished event data
	getSliFinishedEventData := &keptnv2.GetSLIFinishedEventData{
		EventData: keptnv2.EventData{
			Status: keptnv2.StatusSucceeded,
			Result: keptnv2.ResultPass,
		},
		GetSLI: keptnv2.GetSLIFinished{
			IndicatorValues: sliResults,
			Start:           data.GetSLI.Start,
			End:             data.GetSLI.End,
		},
	}

	_, err = myKeptn.SendTaskFinishedEvent(getSliFinishedEventData, ServiceName)

	if err != nil {
		errMsg := fmt.Sprintf("Failed to send task finished CloudEvent (%s), aborting...", err.Error())
		log.Println(errMsg)
		return err
	}

	return nil
}

// HandleProblemEvent handles two problem events:
// - ProblemOpenEventType = "sh.keptn.event.problem.open"
// - ProblemEventType = "sh.keptn.events.problem"
// TODO: add in your handler code
func HandleProblemEvent(myKeptn *keptnv2.Keptn, incomingEvent cloudevents.Event, data *keptn.ProblemEventData) error {
	log.Printf("Handling Problem Event: %s", incomingEvent.Context.GetID())

	// Deprecated since Keptn 0.7.0 - use the HandleActionTriggeredEvent instead

	return nil
}

// HandleActionTriggeredEvent handles action.triggered events
// TODO: add in your handler code
func HandleActionTriggeredEvent(myKeptn *keptnv2.Keptn, incomingEvent cloudevents.Event, data *keptnv2.ActionTriggeredEventData) error {
	log.Printf("Handling Action Triggered Event: %s", incomingEvent.Context.GetID())
	log.Printf("Action=%s\n", data.Action.Action)
	log.Println("Value", data.Action.Value)

	handler, err := NewShipaHandler()
	if err != nil {
		return err
	}

	// check if action is supported
	switch data.Action.Action {
	case "create.framework":
		return handler.action(myKeptn, data, handler.createFramework)
	case "update.framework":
		return handler.action(myKeptn, data, handler.updateFramework)
	case "create.cluster":
		return handler.action(myKeptn, data, handler.createCluster)
	case "update.cluster":
		return handler.action(myKeptn, data, handler.updateCluster)
	case "remove.cluster":
		return handler.action(myKeptn, data, handler.removeCluster)
	case "create.application":
		return handler.action(myKeptn, data, handler.createApp)
	case "deploy.application":
		return handler.action(myKeptn, data, handler.deployApp)

	default:
		log.Printf("Retrieved unknown action %s, skipping...", data.Action.Action)
	}

	return nil
}

type ShipaHandler struct {
	client *shipa.Client
}

const (
	ShipaHost  = "https://target.shipa.cloud"
	ShipaToken = "..."
)

func NewShipaHandler() (*ShipaHandler, error) {
	client, err := shipa.NewClient(ShipaHost, ShipaToken)
	if err != nil {
		log.Println("ERR: failed to create shipa client:", err)
		return nil, err
	}

	return &ShipaHandler{
		client: client,
	}, nil
}

func (s *ShipaHandler) action(myKeptn *keptnv2.Keptn, data *keptnv2.ActionTriggeredEventData, actionFn func(ctx context.Context, data []byte) error) error {
	log.Println("1. Send Action.Started Cloud-Event")
	// -----------------------------------------------------
	// 1. Send Action.Started Cloud-Event
	// -----------------------------------------------------
	myKeptn.SendTaskStartedEvent(data, ServiceName)

	log.Println("2. Implement your remediation action here")
	// -----------------------------------------------------
	// 2. Implement your remediation action here
	// -----------------------------------------------------

	rawData, err := json.Marshal(data.Action.Value)
	if err != nil {
		log.Println("ERR: failed to marshal framework:", err)
		return err
	}

	err = actionFn(context.Background(), rawData)
	if err != nil {
		myKeptn.SendTaskFinishedEvent(&keptnv2.EventData{
			Status:  keptnv2.StatusErrored, // alternative: keptnv2.StatusErrored
			Result:  keptnv2.ResultFailed,  // alternative: keptnv2.ResultFailed
			Message: err.Error(),
		}, ServiceName)
		return err
	}

	log.Println("3. Send Action.Finished Cloud-Event")
	// -----------------------------------------------------
	// 3. Send Action.Finished Cloud-Event
	// -----------------------------------------------------
	myKeptn.SendTaskFinishedEvent(&keptnv2.EventData{
		Status:  keptnv2.StatusSucceeded, // alternative: keptnv2.StatusErrored
		Result:  keptnv2.ResultPass,      // alternative: keptnv2.ResultFailed
		Message: "Successfully created!",
	}, ServiceName)

	return nil
}

func (s *ShipaHandler) createFramework(ctx context.Context, data []byte) error {
	framework := &shipa.PoolConfig{}
	err := json.Unmarshal(data, framework)
	if err != nil {
		log.Println("ERR: failed to unmarshal framework:", err)
		return err
	}

	err = s.client.CreatePoolConfig(ctx, framework)
	if err != nil {
		log.Println("ERR: failed to create framework:", err)
		return err
	}

	return nil
}

func (s *ShipaHandler) updateFramework(ctx context.Context, data []byte) error {
	framework := &shipa.PoolConfig{}
	err := json.Unmarshal(data, framework)
	if err != nil {
		log.Println("ERR: failed to unmarshal framework:", err)
		return err
	}

	err = s.client.UpdatePoolConfig(ctx, framework)
	if err != nil {
		log.Println("ERR: failed to update framework:", err)
		return err
	}

	return nil
}

func (s *ShipaHandler) createCluster(ctx context.Context, data []byte) error {
	cluster := &shipa.Cluster{}
	err := json.Unmarshal(data, cluster)
	if err != nil {
		log.Println("ERR: failed to unmarshal cluster:", err)
		return err
	}

	err = s.client.CreateCluster(ctx, cluster)
	if err != nil {
		log.Println("ERR: failed to create cluster:", err)
		return err
	}

	return nil
}

func (s *ShipaHandler) updateCluster(ctx context.Context, data []byte) error {
	cluster := &shipa.Cluster{}
	err := json.Unmarshal(data, cluster)
	if err != nil {
		log.Println("ERR: failed to unmarshal cluster:", err)
		return err
	}

	err = s.client.UpdateCluster(ctx, cluster)
	if err != nil {
		log.Println("ERR: failed to update cluster:", err)
		return err
	}

	return nil
}

func (s *ShipaHandler) removeCluster(ctx context.Context, data []byte) error {
	cluster := &shipa.Cluster{}
	err := json.Unmarshal(data, cluster)
	if err != nil {
		log.Println("ERR: failed to unmarshal cluster:", err)
		return err
	}

	err = s.client.DeleteCluster(ctx, cluster.Name)
	if err != nil {
		log.Println("ERR: failed to delete cluster:", err)
		return err
	}

	return nil
}

func (s *ShipaHandler) createApp(ctx context.Context, data []byte) error {
	app := &shipa.App{}
	err := json.Unmarshal(data, app)
	if err != nil {
		log.Println("ERR: failed to unmarshal app:", err)
		return err
	}

	err = s.client.CreateApp(ctx, app)
	if err != nil {
		log.Println("ERR: failed to create app:", err)
		return err
	}

	return nil
}

type AppDeployConfig struct {
	Name   string           `json:"name"`
	Deploy *shipa.AppDeploy `json:"deploy"`
}

func (s *ShipaHandler) deployApp(ctx context.Context, data []byte) error {
	app := &AppDeployConfig{}
	err := json.Unmarshal(data, app)
	if err != nil {
		log.Println("ERR: failed to unmarshal app deploy config:", err)
		return err
	}

	err = s.client.DeployApp(ctx, app.Name, app.Deploy)
	if err != nil {
		log.Println("ERR: failed to deploy app:", err)
		return err
	}

	return nil
}
