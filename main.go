package main

import (
	"fmt"
	"github.com/IBM/go-sdk-core/v5/core"
	scc "github.com/ibm/scc-go-sdk/examples/posturemanagementv1"
	"github.com/ibm/scc-go-sdk/posturemanagementv1"
	"os"
)

var collectorIds []string

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

	_, collectorId := scc.CreateCollector(options, accountId)
	credentialId, _ := scc.CreateCredentials(options, accountId, config["CREDENTIAL_PATH"], config["PEM_PATH"])

	collectorIds = append(collectorIds, *collectorId)
	_, scopeId, scopeName := scc.CreateScope(options, accountId, credentialId, collectorIds)
	fmt.Println("Scope created.")
	fmt.Println("Scope ID ", scopeId)
	fmt.Println("Scope Name ", scopeName)
	// loop here
	//scc.ListScopes(options, accountId)
	// end loop
	_, profiles := scc.ListProfiles(options, accountId)

	for _, profile := range profiles {
		fmt.Println("Profile ID: ", *profile.ProfileID)
		fmt.Println("Profile Name: ", *profile.Name)
		fmt.Println("Profile Description: ", *profile.Description)
	}

	_, scanMessage := scc.InitiateValidationScan(options, accountId, scopeId, profileId)
	scc.ListLatestScans()
	scc.ListValidationRuns()
	scc.RetrieveScanSummary()
}
