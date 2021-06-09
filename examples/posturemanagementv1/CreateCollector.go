package examples

import (
	"fmt"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/google/uuid"
	scc "github.com/ibm/scc-go-sdk/posturemanagementv1"
)

func CreateCollector(options scc.PostureManagementV1Options, accountId string, managedBy string) (int, *scc.Collector) {

	service, _ := scc.NewPostureManagementV1(&options)

	source := service.NewCreateCollectorOptions(accountId)
	source.CollectorName = core.StringPtr("test-" + uuid.NewString())
	source.CollectorDescription = core.StringPtr("test collector")
	source.ManagedBy = core.StringPtr(managedBy)
	source.IsPublic = core.BoolPtr(true)
	source.PassPhrase = core.StringPtr("secret")

	reply, response, err := service.CreateCollector(source)

	if err != nil {
		fmt.Println(response.Result)
		fmt.Println("Failed to create collector: ", err)
		panic(err)
	}

	return response.StatusCode, reply
}
