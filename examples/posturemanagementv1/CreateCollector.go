package examples

import (
	"fmt"
	"github.com/google/uuid"
	scc "github.com/ibm/scc-go-sdk/posturemanagementv1"
)

func CreateCollector(options scc.PostureManagementV1Options, accountId string) (int, *string) {

	service, _ := scc.NewPostureManagementV1(&options)

	source := service.NewCreateCollectorOptions(accountId)
	source.SetCollectorName("demo-" + uuid.NewString())
	source.SetCollectorDescription("ibm managed collector for sdk demo")
	source.SetManagedBy("ibm")
	source.SetIsPublic(true)
	source.SetPassPhrase("secret")

	reply, response, err := service.CreateCollector(source)

	if err != nil {
		fmt.Println(response.Result)
		fmt.Println("Failed to create collector: ", err)
		panic(err)
	}

	return response.StatusCode, reply.CollectorID
}
