package examples

import (
	"fmt"
	scc "github.com/ibm/scc-go-sdk/posturemanagementv1"
)

func ListLatestScans(options scc.PostureManagementV1Options, accountId string) *scc.ScansList {
	service, _ := scc.NewPostureManagementV1(&options)

	source := service.NewListLatestScansOptions(accountId)

	scansList, response, err := service.ListLatestScans(source)

	if err != nil {
		fmt.Println(response.Result)
		fmt.Println("Failed to create scope: ", err)
		panic(err)
	}
	return scansList

}
