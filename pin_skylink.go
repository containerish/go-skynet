package skynet

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"gitlab.com/NebulousLabs/errors"
)

var (
	PinEndpoint      = "/skynet/pin/"
	SkylinkHeaderKey = "skynet-skylink"
)

func (sc *SkynetClient) PinSkylink(skylink string) (string, error) {
	skylink = strings.TrimPrefix(skylink, "sia://")

	resp, err := sc.executeRequest(
		requestOptions{
			reqBody: nil,
			query:   url.Values{},
			Options: Options{
				EndpointPath:      PinEndpoint,
				SkynetAPIKey:      sc.Options.SkynetAPIKey,
				CustomUserAgent:   sc.Options.CustomUserAgent,
				customContentType: sc.Options.customContentType,
			},
			method:    http.MethodPost,
			extraPath: skylink,
		},
	)
	if err != nil {
		return "", errors.AddContext(err, "could not execute request")
	}

	// previously, skynet returned no content(204), with verison bump it has started returning a sia path in    //body and headers along with statusOK
	if resp.StatusCode == http.StatusNoContent || resp.StatusCode == http.StatusOK {
		pinLink := resp.Header.Get(SkylinkHeaderKey)
		return pinLink, nil
	}
	return "", fmt.Errorf("expected response status code to be %d or %d but got %d", http.StatusNoContent, http.StatusOK, resp.StatusCode)
}
