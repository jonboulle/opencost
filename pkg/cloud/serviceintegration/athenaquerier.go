package serviceintegration

import (
	"context"
	"fmt"
	"github.com/opencost/opencost/pkg/cloud/config"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/athena"
	"github.com/aws/aws-sdk-go-v2/service/athena/types"
	"github.com/opencost/opencost/pkg/kubecost"
	"github.com/opencost/opencost/pkg/log"
	"github.com/opencost/opencost/pkg/util/stringutil"
)

type AthenaQuerier struct {
	config.AthenaConfiguration
}

func (aq *AthenaQuerier) Equals(config config.Config) bool {
	thatConfig, ok := config.(*AthenaQuerier)
	if !ok {
		return false
	}

	return aq.AthenaConfiguration.Equals(&thatConfig.AthenaConfiguration)
}

// QueryAthenaPaginated executes athena query and processes results. An error from this method indicates a
// FAILED_CONNECTION CloudConnectionStatus and should immediately stop the caller to maintain the correct CloudConnectionStatus
func (aq *AthenaQuerier) QueryAthenaPaginated(ctx context.Context, query string, fn func(*athena.GetQueryResultsOutput) bool) error {

	queryExecutionCtx := &types.QueryExecutionContext{
		Database: aws.String(aq.Database),
	}

	resultConfiguration := &types.ResultConfiguration{
		OutputLocation: aws.String(aq.Bucket),
	}
	startQueryExecutionInput := &athena.StartQueryExecutionInput{
		QueryString:           aws.String(query),
		QueryExecutionContext: queryExecutionCtx,
		ResultConfiguration:   resultConfiguration,
	}

	// Only set if there is a value, the default input is nil
	if aq.Workgroup != "" {
		startQueryExecutionInput.WorkGroup = aws.String(aq.Workgroup)
	}

	// Create Athena Client
	cfg, err := aq.CreateAWSConfig()
	if err != nil {
		log.Errorf("QueryAthenaPaginated: Could not retrieve AWS Configuration: %s", err.Error())
	}
	cli := athena.NewFromConfig(cfg)

	// Query Athena
	startQueryExecutionOutput, err := cli.StartQueryExecution(ctx, startQueryExecutionInput)
	if err != nil {
		return fmt.Errorf("QueryAthenaPaginated: start query error: %s", err.Error())
	}
	err = waitForQueryToComplete(ctx, cli, startQueryExecutionOutput.QueryExecutionId)
	if err != nil {
		return fmt.Errorf("QueryAthenaPaginated: query execution error: %s", err.Error())
	}
	queryResultsInput := &athena.GetQueryResultsInput{
		QueryExecutionId: startQueryExecutionOutput.QueryExecutionId,
	}
	getQueryResultsPaginator := athena.NewGetQueryResultsPaginator(cli, queryResultsInput)
	for getQueryResultsPaginator.HasMorePages() {
		pg, err := getQueryResultsPaginator.NextPage(ctx)
		if err != nil {
			log.Errorf("queryAthenaPaginated: NextPage error: %s", err.Error())
			continue
		}
		fn(pg)
	}
	return nil
}

func waitForQueryToComplete(ctx context.Context, client *athena.Client, queryExecutionID *string) error {
	inp := &athena.GetQueryExecutionInput{
		QueryExecutionId: queryExecutionID,
	}
	isQueryStillRunning := true
	for isQueryStillRunning {
		qe, err := client.GetQueryExecution(ctx, inp)
		if err != nil {
			return err
		}
		if qe.QueryExecution.Status.State == "SUCCEEDED" {
			isQueryStillRunning = false
			continue
		}
		if qe.QueryExecution.Status.State != "RUNNING" && qe.QueryExecution.Status.State != "QUEUED" {
			return fmt.Errorf("no query results available for query %s", *queryExecutionID)
		}
		time.Sleep(2 * time.Second)
	}
	return nil
}

// GetAthenaRowValue retrieve value from athena row based on column names and used stringutil.Bank() to prevent duplicate
// allocation of strings
func GetAthenaRowValue(row types.Row, queryColumnIndexes map[string]int, columnName string) string {
	columnIndex, ok := queryColumnIndexes[columnName]
	if !ok {
		return ""
	}
	valuePointer := row.Data[columnIndex].VarCharValue
	if valuePointer == nil {
		return ""
	}
	return stringutil.Bank(*valuePointer)
}

// getAthenaRowValueFloat retrieve value from athena row based on column names and convert to float if possible
func GetAthenaRowValueFloat(row types.Row, queryColumnIndexes map[string]int, columnName string) (float64, error) {

	columnIndex, ok := queryColumnIndexes[columnName]
	if !ok {
		return 0.0, fmt.Errorf("getAthenaRowValueFloat: missing column index: %s", columnName)
	}

	valuePointer := row.Data[columnIndex].VarCharValue
	if valuePointer == nil {
		return 0.0, fmt.Errorf("getAthenaRowValueFloat: nil field")
	}

	cost, err := strconv.ParseFloat(*valuePointer, 64)
	if err != nil {
		return cost, fmt.Errorf("getAthenaRowValueFloat: failed to parse %s: '%s': %s", columnName, *valuePointer, err.Error())
	}
	return cost, nil
}

func SelectAWSCategory(isNode, isVol, isNetwork bool, service string) string {
	// Network has the highest priority and is based on the usage type ending in "Bytes"
	if isNetwork {
		return kubecost.NetworkCategory
	}
	// The node and volume conditions are mutually exclusive.
	// Provider ID has prefix "i-"
	if isNode {
		return kubecost.ComputeCategory
	}
	// Provider ID has prefix "vol-"
	if isVol {
		return kubecost.StorageCategory
	}

	// Default categories based on service
	switch strings.ToUpper(service) {
	case "AWSELB", "AWSGLUE", "AMAZONROUTE53":
		return kubecost.NetworkCategory
	case "AMAZONEC2", "AWSLAMBDA", "AMAZONELASTICACHE":
		return kubecost.ComputeCategory
	case "AMAZONEKS":
		return kubecost.ManagementCategory
	case "AMAZONS3", "AMAZONATHENA", "AMAZONRDS", "AMAZONDYNAMODB", "AWSSECRETSMANAGER", "AMAZONFSX":
		return kubecost.StorageCategory
	default:
		return kubecost.OtherCategory
	}
}

var parseARNRx = regexp.MustCompile("^.+\\/(.+)?") // Capture "a406f7761142e4ef58a8f2ba478d2db2" from "arn:aws:elasticloadbalancing:us-east-1:297945954695:loadbalancer/a406f7761142e4ef58a8f2ba478d2db2"

func ParseARN(id string) string {
	match := parseARNRx.FindStringSubmatch(id)
	if len(match) == 0 {
		if id != "" {
			log.DedupedInfof(10, "aws.parseARN: failed to parse %s", id)
		}
		return id
	}
	return match[len(match)-1]
}

func GetAthenaQueryFunc(fn func(types.Row)) func(*athena.GetQueryResultsOutput) bool {
	pageNum := 0
	processItemQueryResults := func(page *athena.GetQueryResultsOutput) bool {
		if page == nil {
			log.Errorf("AthenaQuerier: Athena page is nil")
			return false
		} else if page.ResultSet == nil {
			log.Errorf("AthenaQuerier: Athena page.ResultSet is nil")
			return false
		}
		rows := page.ResultSet.Rows
		if pageNum == 0 {
			rows = page.ResultSet.Rows[1:len(page.ResultSet.Rows)]
		}

		for _, row := range rows {
			fn(row)
		}
		pageNum++
		return true
	}
	return processItemQueryResults
}
