package main

import (
	"fmt"
	"github.com/IBM/go-sdk-core/v5/core"
	scc "github.com/ibm/scc-go-sdk/examples/posturemanagementv1"
	"github.com/ibm/scc-go-sdk/posturemanagementv1"
	"os"
	"strconv"
)

var (
	collectorIds   []string
	scopeCondition bool
	scanId         string
)

func main() {
	const externalConfigFile = "./posture_management_v1.env"

	os.Stat(externalConfigFile)

	os.Setenv("IBM_CREDENTIALS_FILE", externalConfigFile)
	config, _ := core.GetServiceProperties(posturemanagementv1.DefaultServiceName)
	accountId := config["ACCOUNT_ID"]
	apiKey := config["IAM_API_KEY"]

	authUrl := config["IAM_APIKEY_URL"]
	apiUrl := config["API_URL"]

	authenticator := core.IamAuthenticator{
		ApiKey: apiKey,
		URL:    authUrl,
	}

	options := posturemanagementv1.PostureManagementV1Options{
		Authenticator: &authenticator,
		URL:           apiUrl,
	}

	//_, collectorId := scc.CreateCollector(options, accountId)

	fmt.Println("Creating Collector")
	_, collector := scc.CreateCollector(options, accountId, "customer")

	fmt.Println("Collector created.")
	fmt.Println("Collector ID: " + *collector.CollectorID)

	fmt.Println("Creating Credentials")
	credentials, _ := scc.CreateCredentials(options, accountId, config["CREDENTIAL_PATH"], config["PEM_PATH"])
	fmt.Println("Credentials created.")
	fmt.Println("Credential ID: " + *credentials.CredentialID)
	fmt.Println("Credential Name: " + *credentials.CredentialName)
	fmt.Println("Time created: " + credentials.CreatedTime.String())

	//collectorIds = append(collectorIds, *collector.CollectorID)
	collectorIds = append(collectorIds, "2082")

	fmt.Println("Creating Scope")
	scope := scc.CreateScope(options, accountId, *credentials.CredentialID, collectorIds)

	fmt.Println("Scope created.")
	fmt.Println("Scope ID: " + *scope.ScopeID)
	fmt.Println("Scope Name: " + *scope.ScopeName)
	fmt.Println("Scope Description: " + *scope.ScopeDescription)
	fmt.Println("Scope Environment Type: " + *scope.EnvironmentType)
	fmt.Println("Created time: " + scope.CreatedTime.String())
	fmt.Println("Modified time: " + scope.ModifiedTime.String())

	fmt.Println("Checking Discovery Status:")
	for scopeCondition == false {
		scopeCondition, _ = scc.ListScopes(options, accountId, *scope.ScopeName, *scope.ScopeID, "discovery_completed")
	}
	fmt.Println("Listing Profiles")
	_, profiles := scc.ListProfiles(options, accountId)

	for _, profile := range profiles {
		fmt.Println("Profile ID: ", *profile.ProfileID)
		fmt.Println("Profile Name: ", *profile.Name)
		fmt.Println("Profile Description: ", *profile.Description)
	}

	fmt.Println("Creating Scan")
	_, scanMessage := scc.InitiateValidationScan(options, accountId, *scope.ScopeID, "48")

	fmt.Println("Scan initiated: " + *scanMessage)

	fmt.Println("Checking Scan Status:")
	scopeCondition = false
	for scopeCondition == false {
		scopeCondition, _ = scc.ListScopes(options, accountId, *scope.ScopeName, *scope.ScopeID, "validation_completed")
	}
	fmt.Println("Listing Latest Scans:")
	scansList := scc.ListLatestScans(options, accountId)

	fmt.Println("Total scans count: " + strconv.FormatInt(*scansList.TotalCount, 10))
	if scansList.LatestScans != nil {
		scanId = *scansList.LatestScans[0].ScanID
	}
	for _, scan := range scansList.LatestScans {
		fmt.Println("Scan ID: ", *scan.ScanID)
		fmt.Println("Profile: ", *scan.ProfileName)
		fmt.Println("Report run by: ", *scan.ReportRunBy)
		fmt.Println("Scan Controls Pass Count: ", *scan.Result.ControlsPassCount)
		fmt.Println("Scan Controls Fail Count: ", *scan.Result.ControlsFailCount)
		fmt.Println("Scan Controls NA Count: ", *scan.Result.ControlsNaCount)
		fmt.Println("Scan Controls U2P Count: ", *scan.Result.ControlsU2pCount)
		fmt.Println("Scan Controls Total Count: ", *scan.Result.ControlsTotalCount)
	}

	fmt.Println("Retrieving summary of scan " + scanId)
	_, scanSummary := scc.RetrieveScanSummary(options, accountId, scanId, "48")
	for _, control := range scanSummary.Controls {
		fmt.Println("Control ID: " + *control.ControlID)
		fmt.Println("Control Description: " + *control.ControlDesciption)
	}

	fmt.Println("Listing Validation Runs:")
	_, summariesList := scc.ListValidationRuns(options, accountId, *scope.ScopeID, "48")

	for _, summary := range summariesList.Summaries {
		fmt.Println("Scan ID: " + *summary.ScanID)
		fmt.Println("Scan Status: " + *summary.Status)
	}

}
