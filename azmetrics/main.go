package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/monitor/query/azmetrics"
)

func main() {

	credentials, err := azidentity.NewClientSecretCredential(
		os.Getenv("AZURE_TENANT_ID"),
		os.Getenv("AZURE_CLIENT_ID"),
		os.Getenv("AZURE_CLIENT_SECRET"),
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	client, err := azmetrics.NewClient(
		"https://eastus2.metrics.monitor.azure.com",
		credentials,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Ready to go!")

	subscriptionID := os.Getenv("AZURE_SUBSCRIPTION_ID")
	resourceGroup := os.Getenv("AZURE_RESOURCE_GROUP")
	namespace := "Microsoft.KeyVault/vaults"
	metricNames := []string{
		"ServiceApiResult",
	}
	resourceIDs := []string{
		fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.KeyVault/vaults/mbrancar10s001", subscriptionID, resourceGroup),
		fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.KeyVault/vaults/mbrancar10s002", subscriptionID, resourceGroup),
		fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.KeyVault/vaults/mbrancar10s003", subscriptionID, resourceGroup),
		fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.KeyVault/vaults/mbrancar10s004", subscriptionID, resourceGroup),
		fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.KeyVault/vaults/mbrancar10s005", subscriptionID, resourceGroup),
		fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.KeyVault/vaults/mbrancar10s006", subscriptionID, resourceGroup),
		fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.KeyVault/vaults/mbrancar10s007", subscriptionID, resourceGroup),
		fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.KeyVault/vaults/mbrancar10s008", subscriptionID, resourceGroup),
		fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.KeyVault/vaults/mbrancar10s009", subscriptionID, resourceGroup),
		fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.KeyVault/vaults/mbrancar10s010", subscriptionID, resourceGroup),
	}
	options := azmetrics.QueryResourcesOptions{
		Aggregation: ptr("Count"),
		StartTime:   ptr("2024-11-18T07:18:13.001Z"),
		EndTime:     ptr("2024-11-18T07:19:13.001Z"),
		Filter:      ptr("ActivityType eq '*' AND ActivityName eq '*' AND StatusCode eq '*' AND StatusCodeClass eq '*'"),
		Interval:    ptr("PT1M"),
	}

	// fmt.Println("----------------------------------------------------")
	// fmt.Println("Querying SINGLE resources")
	// fmt.Println("----------------------------------------------------")

	// for _, resourceID := range resourceIDs {
	// 	res, err := client.QueryResources(
	// 		context.Background(),
	// 		subscriptionID,
	// 		namespace,
	// 		metricNames,
	// 		azmetrics.ResourceIDList{
	// 			ResourceIDs: []string{resourceID},
	// 		},
	// 		&options,
	// 	)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}

	// 	displayResponse(res)
	// }

	fmt.Println("----------------------------------------------------")
	fmt.Println("Querying resources as a GROUP")
	fmt.Println("----------------------------------------------------")

	res, err := client.QueryResources(
		context.Background(),
		subscriptionID,
		namespace,
		metricNames,
		azmetrics.ResourceIDList{
			ResourceIDs: resourceIDs,
		},
		&options,
	)
	if err != nil {
		log.Fatal(err)
	}

	displayResponse(res)
}

func displayResponse(res azmetrics.QueryResourcesResponse) {
	for _, metricData := range res.Values {
		fmt.Println(*metricData.ResourceID)
		for _, metricValue := range metricData.Values {
			// fmt.Println("timeseries", len(metricValue.TimeSeries))
			for _, timeseries := range metricValue.TimeSeries {
				for _, data := range timeseries.Data {
					displayMetricValue(data)
				}
			}
		}
	}
}

func ptr(s string) *string {
	return &s
}

func displayMetricValue(metricValue azmetrics.MetricValue) {
	if metricValue.Average != nil {
		fmt.Println("average", *metricValue.Average)
	}
	if metricValue.Count != nil {
		fmt.Println("count", *metricValue.Count)
	}
	if metricValue.Maximum != nil {
		fmt.Println("maximum", *metricValue.Maximum)
	}
	if metricValue.Minimum != nil {
		fmt.Println("minimum", *metricValue.Minimum)
	}
	if metricValue.Total != nil {
		fmt.Println("total", *metricValue.Total)
	}
}
