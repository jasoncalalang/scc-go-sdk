package examples

import (
	"encoding/json"
	"fmt"
	scc "github.com/ibm/scc-go-sdk/posturemanagementv1"
	"io/ioutil"
)

func RetrieveScanSummary(options scc.PostureManagementV1Options, accountId string, scanId string, profileId string) {

	service, _ := scc.NewPostureManagementV1(&options)

	source := service.NewScansSummaryOptions(accountId, scanId, profileId)

	reply, response, err := service.ScansSummary(source)
	if err != nil {
		fmt.Println(response.Result)
		fmt.Println("Failed to retrieve scan summary: ", err)
		return
	}

	out, err := json.Marshal(reply)
	if err != nil {
		fmt.Println("Failed to marshal json file: ", err)
		return
	}

	err = ioutil.WriteFile("scans.json", out, 0644)
	if err != nil {
		fmt.Println("Failed to write output file: ", err)
		return
	}

}
