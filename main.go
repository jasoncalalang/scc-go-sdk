package main

import (
	"fmt"
	"github.com/IBM/go-sdk-core/v5/core"
	scc "github.com/ibm/scc-go-sdk/examples/posturemanagementv1"
	"github.com/ibm/scc-go-sdk/posturemanagementv1"
	"os"
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

	fmt.Println("Creating IBM-Managed Collector...")
	_, collector := scc.CreateCollector(options, accountId, "ibm")

	fmt.Println("Collector created.")
	fmt.Println("Collector ID: " + *collector.CollectorID)
	fmt.Println("Press enter to continue...")
	fmt.Scanln()

	fmt.Println("Creating Credentials...")
	credentials, _ := scc.CreateCredentials(options, accountId, config["CREDENTIAL_PATH"], config["PEM_PATH"])
	fmt.Println("Credentials created.")
	fmt.Println("Credential ID: " + *credentials.CredentialID)
	fmt.Println("Credential Name: " + *credentials.CredentialName)
	fmt.Println("Time created: " + credentials.CreatedTime.String())

	collectorIds = append(collectorIds, *collector.CollectorID) //make sure that managedBy=ibm for IBM MC
	//collectorIds = append(collectorIds, "2082") // for customer-managed collector

	fmt.Println("Press enter to continue...")
	fmt.Scanln()

	fmt.Println("Creating Scope...")
	scope := scc.CreateScope(options, accountId, *credentials.CredentialID, collectorIds)

	fmt.Println("Scope created.")
	fmt.Println("Scope ID: " + *scope.ScopeID)
	fmt.Println("Scope Name: " + *scope.ScopeName)
	fmt.Println("Scope Description: " + *scope.ScopeDescription)
	fmt.Println("Scope Environment Type: " + *scope.EnvironmentType)
	fmt.Println("Created time: " + scope.CreatedTime.String())
	fmt.Println("Modified time: " + scope.ModifiedTime.String())
	fmt.Println("Discovery Triggered.")

	fmt.Println("Press enter to continue...")
	fmt.Scanln()

	fmt.Println("Checking Discovery Status...")
	for scopeCondition == false {
		scopeCondition, _ = scc.ListScopes(options, accountId, *scope.ScopeName, *scope.ScopeID, "discovery_completed")
	}

	fmt.Println("Press enter to continue...")
	fmt.Scanln()

	fmt.Println("Listing Profiles...")
	_, profiles := scc.ListProfiles(options, accountId)

	for _, profile := range profiles {
		fmt.Println("Profile ID: ", *profile.ProfileID)
		fmt.Println("Profile Name: ", *profile.Name)
		fmt.Println("Profile Description: ", *profile.Description)
	}

	fmt.Println("Press enter to continue...")
	fmt.Scanln()

	fmt.Println("Triggered Validation Scan with IBM Cloud Best Practices Controls 1.0")
	_, scanMessage := scc.InitiateValidationScan(options, accountId, *scope.ScopeID, "48")

	fmt.Println("Scan initiated: " + *scanMessage)

	fmt.Println("Checking Scan Status...")
	scopeCondition = false
	for scopeCondition == false {
		scopeCondition, _ = scc.ListScopes(options, accountId, *scope.ScopeName, *scope.ScopeID, "validation_completed")
	}

	fmt.Println("Press enter to continue...")
	fmt.Scanln()

	scansList := scc.ListLatestScans(options, accountId)
	fmt.Println("Showing scan of " + *scansList.LatestScans[0].ScanID) //showing scan id "scanId"

	if scansList.LatestScans != nil {
		scanId = *scansList.LatestScans[0].ScanID
	}
	fmt.Println("Scan ID: ", *scansList.LatestScans[0].ScanID)
	fmt.Println("Profile: ", *scansList.LatestScans[0].ProfileName)
	fmt.Println("Report run by: ", *scansList.LatestScans[0].ReportRunBy)
	fmt.Println("Scan Controls Pass Count: ", *scansList.LatestScans[0].Result.ControlsPassCount)
	fmt.Println("Scan Controls Fail Count: ", *scansList.LatestScans[0].Result.ControlsFailCount)
	fmt.Println("Scan Controls NA Count: ", *scansList.LatestScans[0].Result.ControlsNaCount)
	fmt.Println("Scan Controls U2P Count: ", *scansList.LatestScans[0].Result.ControlsU2pCount)
	fmt.Println("Scan Controls Total Count: ", *scansList.LatestScans[0].Result.ControlsTotalCount)

	fmt.Println("Press enter to continue...")
	fmt.Scanln()

	fmt.Println("Retrieving summary of scan " + scanId + "...")
	scc.RetrieveScanSummary(options, accountId, scanId, "48")

	fmt.Println("Press enter to continue...")
	fmt.Scanln()

}
