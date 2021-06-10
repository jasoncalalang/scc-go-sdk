package examples

import (
	"encoding/hex"
	"fmt"
	"github.com/google/uuid"
	scc "github.com/ibm/scc-go-sdk/posturemanagementv1"
)

func CreateScope(options scc.PostureManagementV1Options, accountId string, credentialId string, collectorIds []string) *scc.Scope {
	service, _ := scc.NewPostureManagementV1(&options)
	uid := hex.EncodeToString(uuid.New().NodeID()[:4])

	source := service.NewCreateScopeOptions(accountId)
	source.SetScopeName("SDK DEMO " + uid)
	source.SetScopeDescription("This is a scope created from the SCC SDK")
	source.SetEnvironmentType("ibm")
	source.SetCollectorIds(collectorIds)
	source.SetCredentialID(credentialId)

	scope, response, err := service.CreateScope(source)

	if err != nil {
		fmt.Println(response.Result)
		fmt.Println("Failed to create scope: ", err)
		panic(err)
	}

	return scope

}
