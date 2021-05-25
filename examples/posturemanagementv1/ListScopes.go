package examples

import (
	"fmt"
	"github.com/IBM/go-sdk-core/v5/core"
	scc "github.com/ibm-cloud-security/scc-go-sdk/posturemanagementv1"
	"os"
)

func ListScopes() {
	apiKey := os.Getenv("IAM_API_KEY")
	url := os.Getenv("IAM_APIKEY_URL")
	accountId := os.Getenv("ACCOUNT_ID")
	authenticator := &core.IamAuthenticator{
		ApiKey: apiKey,
		URL:    url, //use for dev/preprod env

	}
	service, _ := scc.NewPostureManagementV1(&scc.PostureManagementV1Options{
		Authenticator: authenticator,
		URL:           "https://asap-dev.compliance.test.cloud.ibm.com", //Specify url or use default
	})

	source := service.NewListScopesOptions(accountId)

	result, response, err := service.ListScopes(source)

	if err != nil {
		fmt.Println(response.Result)
		fmt.Println("Failed to create scope: ", err)
		return
	}
	fmt.Println("Status Code: ", response.GetStatusCode())
	fmt.Println("Scope ID: ", *result.Scopes[0].ScopeID)
	fmt.Println("Created Time: ", *result.Scopes[0].CreatedTime)
	fmt.Println("Modified Time: ", *result.Scopes[0].ModifiedTime)

}
