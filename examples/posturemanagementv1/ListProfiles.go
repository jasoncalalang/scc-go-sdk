package examples

import (
	"fmt"
	"github.com/IBM/go-sdk-core/v5/core"
	scc "github.com/ibm/scc-go-sdk/posturemanagementv1"
)

func ListProfiles(options scc.PostureManagementV1Options, accountId string) (int, []scc.ProfileItem) {
	service, _ := scc.NewPostureManagementV1(&options)
	source := service.NewListProfilesOptions(accountId)
	source.Name = core.StringPtr("IBM Cloud Best Practices Controls 1.0")
	reply, response, err := service.ListProfiles(source)

	if err != nil {
		fmt.Println(response.Result)
		fmt.Println("Failed to list profiles: ", err)
		panic(err)
	}

	return response.StatusCode, reply.Profiles
}
