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

	fmt.Println("Creating Collector")
	_, collector := scc.CreateCollector(options, accountId, "customer")

	fmt.Println("Collector created.")
	fmt.Println("Collector ID: " + *collector.CollectorID)
	fmt.Println("Press enter to continue...")
	fmt.Scanln()

	fmt.Println("Creating Credentials")
	credentials, _ := scc.CreateCredentials(options, accountId, config["CREDENTIAL_PATH"], config["PEM_PATH"])
	fmt.Println("Credentials created.")
	fmt.Println("Credential ID: " + *credentials.CredentialID)
	fmt.Println("Credential Name: " + *credentials.CredentialName)
	fmt.Println("Time created: " + credentials.CreatedTime.String())

	//collectorIds = append(collectorIds, *collector.CollectorID)
	collectorIds = append(collectorIds, "2082")

	fmt.Println("Press enter to continue...")
	fmt.Scanln()

	fmt.Println("Creating Scope")
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

	fmt.Println("Checking Discovery Status:")
	for scopeCondition == false {
		scopeCondition, _ = scc.ListScopes(options, accountId, *scope.ScopeName, *scope.ScopeID, "discovery_completed")
	}

	fmt.Println("Press enter to continue...")
	fmt.Scanln()

	fmt.Println("Listing Profiles")
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

	fmt.Println("Checking Scan Status:")
	scopeCondition = false
	for scopeCondition == false {
		scopeCondition, _ = scc.ListScopes(options, accountId, *scope.ScopeName, *scope.ScopeID, "validation_completed")
	}

	fmt.Println("Press enter to continue...")
	fmt.Scanln()

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

	fmt.Println("Press enter to continue...")
	fmt.Scanln()

	fmt.Println("Retrieving summary of scan " + "32821")
	_, scanSummary := scc.RetrieveScanSummary(options, accountId, "32821", "48")
	for _, control := range scanSummary.Controls {
		fmt.Println("Control ID: " + *control.ControlID)
		fmt.Println("Control Description: " + *control.ControlDesciption)
		fmt.Println("Control Status: " + *control.Status)
		for _, goal := range control.Goals {
			fmt.Println("Goal Description: " + *goal.GoalDescription)
			fmt.Println("Control Status: " + *goal.Status)
			fmt.Println("Control Status: " + *goal.Severity)
			for _, resource := range goal.ResourceResult {
				if resource.ResourceName != nil {
					fmt.Println("Resource name: " + *resource.ResourceName)
					fmt.Println("Resource type: " + *resource.ResourceType)
					fmt.Println("Resource status: " + *resource.ResourceStatus)
					fmt.Println("Display Expected Value: " + *resource.DisplayExpectedValue)
					fmt.Println("Actual Value: " + *resource.ActualValue)
					fmt.Println("Results Info: " + *resource.ResultsInfo)
					fmt.Println("NA Reason: " + *resource.NaReason)
				}
			}
		}
	}

	fmt.Println("Press enter to continue...")
	fmt.Scanln()

	fmt.Println("Listing Validation Runs:")
	_, summariesList := scc.ListValidationRuns(options, accountId, "28632", "48")

	for _, summary := range summariesList.Summaries {
		fmt.Println("Scan ID: " + *summary.ScanID)
		fmt.Println("Scan Name: " + *summary.ScanName)
		fmt.Println("Scope Name: " + *summary.ScopeName)
		fmt.Println("\nValidation Result for " + *summary.Profile.ProfileName)
		fmt.Println("Goals Pass Count: " + strconv.FormatInt(*summary.Profile.ValidationResult.GoalsPassCount, 10))
		fmt.Println("Goals U2P Count: " + strconv.FormatInt(*summary.Profile.ValidationResult.GoalsU2pCount, 10))
		fmt.Println("Goals NA Count: " + strconv.FormatInt(*summary.Profile.ValidationResult.GoalsNaCount, 10))
		fmt.Println("Goals Fail Count: " + strconv.FormatInt(*summary.Profile.ValidationResult.GoalsFailCount, 10))
		fmt.Println("Goals Total Count: " + strconv.FormatInt(*summary.Profile.ValidationResult.GoalsTotalCount, 10))
		fmt.Println("Controls Pass Count: " + strconv.FormatInt(*summary.Profile.ValidationResult.ControlsPassCount, 10))
		fmt.Println("Controls U2P Count: " + strconv.FormatInt(*summary.Profile.ValidationResult.ControlsU2pCount, 10))
		fmt.Println("Controls NA Count: " + strconv.FormatInt(*summary.Profile.ValidationResult.ControlsNaCount, 10))
		fmt.Println("Controls Fail Count: " + strconv.FormatInt(*summary.Profile.ValidationResult.ControlsFailCount, 10))
		fmt.Println("Controls Total Count: " + strconv.FormatInt(*summary.Profile.ValidationResult.ControlsTotalCount, 10))
	}

}
