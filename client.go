package skynet

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"gitlab.com/NebulousLabs/errors"
)

type (
	// SkynetClient is the Skynet Client which can be used to access Skynet.
	SkynetClient struct {
		PortalURL string
		Options   Options
	}

	// requestOptions contains the options for a request.
	requestOptions struct {
		Options

		method    string
		reqBody   io.Reader
		extraPath string
		query     url.Values
	}
)

// NewSkynetClient creates a new Skynet Client which can be used to access Skynet.
// Pass in "" for the portal to let the function select one for you.
func NewSkynetClient(portalURL string) SkynetClient {
	if portalURL == "" {
		portalURL = DefaultPortalURL
	}
	return SkynetClient{
		PortalURL: portalURL,
	}
}

// SetCustomOptions sets the custom options for this client.
func (sc *SkynetClient) SetCustomOptions(customOptions Options) {
	sc.Options = customOptions
}

// executeRequest makes and executes a request.
func (sc *SkynetClient) executeRequest(config requestOptions) (*http.Response, error) {
	url := sc.PortalURL
	method := config.method
	reqBody := config.reqBody

	// Set options, prioritizing options passed to the API calls.
	opts := sc.Options
	if config.EndpointPath != "" {
		opts.EndpointPath = config.EndpointPath
	}
	if config.APIKey != "" {
		opts.APIKey = config.APIKey
	}
	if config.CustomUserAgent != "" {
		opts.CustomUserAgent = config.CustomUserAgent
	}
	if config.customContentType != "" {
		opts.customContentType = config.customContentType
	}

	// Make the URL.
	url = makeURL(url, opts.EndpointPath, nil)
	url = makeURL(url, config.extraPath, config.query)

	// Create the request.
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, errors.AddContext(err, fmt.Sprintf("could not create %v request", method))
	}
	if opts.APIKey != "" {
		req.SetBasicAuth("", opts.APIKey)
	}
	if opts.CustomUserAgent != "" {
		req.Header.Set("User-Agent", opts.CustomUserAgent)
	}
	if opts.customContentType != "" {
		req.Header.Set("Content-Type", opts.customContentType)
	}

	// Execute the request.
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.AddContext(err, "could not execute request")
	}
	if resp.StatusCode >= 400 {
		return nil, errors.AddContext(makeResponseError(resp), "error code received")
	}

	return resp, nil
}