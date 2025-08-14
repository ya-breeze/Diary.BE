# \ItemsAPI

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**GetItems**](ItemsAPI.md#GetItems) | **Get** /v1/items | get diary items



## GetItems

> ItemsResponse GetItems(ctx).Date(date).Execute()

get diary items

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
    "time"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	date := time.Now() // string | filter items by date (optional) (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ItemsAPI.GetItems(context.Background()).Date(date).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ItemsAPI.GetItems``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetItems`: ItemsResponse
	fmt.Fprintf(os.Stdout, "Response from `ItemsAPI.GetItems`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiGetItemsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **date** | **string** | filter items by date (optional) | 

### Return type

[**ItemsResponse**](ItemsResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

