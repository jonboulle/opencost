package serviceintegration

import (
	"bytes"
	"context"
	"fmt"
	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/opencost/opencost/pkg/cloud/config"
	"github.com/opencost/opencost/pkg/log"
	"net/url"
	"strings"
)

// AzureStorageConnection provides access to Azure Storage
type AzureStorageConnection struct {
	config.AzureStorageConfiguration
}

func (asc *AzureStorageConnection) Equals(config config.Config) bool {
	thatConfig, ok := config.(*AzureStorageConnection)
	if !ok {
		return false
	}

	return asc.AzureStorageConfiguration.Equals(&thatConfig.AzureStorageConfiguration)
}

func (asc *AzureStorageConnection) getContainer() (*azblob.ContainerURL, error) {

	credential, err := asc.Configurer.GetBlobCredentials()
	if err != nil {
		return nil, err
	}

	p := azblob.NewPipeline(credential, azblob.PipelineOptions{})

	// From the Azure portal, get your storage account blob service URL endpoint.
	URL, _ := url.Parse(
		fmt.Sprintf(asc.getBlobURLTemplate(), asc.Account, asc.Container))

	// Create a ContainerURL object that wraps the container URL and a request
	// pipeline to make requests.
	containerURL := azblob.NewContainerURL(*URL, p)
	return &containerURL, nil
}

// getBlobURLTemplate returns the correct BlobUrl for whichever Cloud storage account is specified by the AzureCloud configuration
// defaults to the Public Cloud template
func (asc *AzureStorageConnection) getBlobURLTemplate() string {
	// Use gov cloud blob url if gov is detected in AzureCloud
	if strings.Contains(strings.ToLower(asc.Cloud), "gov") {
		return "https://%s.blob.core.usgovcloudapi.net/%s"
	}
	// default to Public Cloud template
	return "https://%s.blob.core.windows.net/%s"
}

func (asc *AzureStorageConnection) DownloadBlob(blobName string, containerURL *azblob.ContainerURL, ctx context.Context) ([]byte, error) {
	log.Infof("Azure Storage: retrieving blob: %v", blobName)

	blobURL := containerURL.NewBlobURL(blobName)
	downloadResponse, err := blobURL.Download(ctx, 0, azblob.CountToEnd, azblob.BlobAccessConditions{}, false, azblob.ClientProvidedKeyOptions{})
	if err != nil {
		return nil, err
	}
	// NOTE: automatically retries are performed if the connection fails
	bodyStream := downloadResponse.Body(azblob.RetryReaderOptions{MaxRetryRequests: 20})

	// read the body into a buffer
	downloadedData := bytes.Buffer{}
	_, err = downloadedData.ReadFrom(bodyStream)
	if err != nil {
		return nil, err
	}
	return downloadedData.Bytes(), nil
}
