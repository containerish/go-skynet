package skynet

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type UploadReponse struct {
	Skylink string `json:"skylink"`
}

type UploadOptions struct {
	portalUrl           string
	portalUploadPath    string
	portalFileFieldname string
	customFilename      string
	tryParseResponse    bool
}

var DefaultUploadOptions = UploadOptions{
	portalUrl:           "https://siasky.net/",
	portalUploadPath:    "/api/skyfile",
	portalFileFieldname: "file",
	customFilename:      "",
	tryParseResponse:    true,
}

func UploadFile(path string, opts UploadOptions) (string, error) {
	// open the file
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// set filename
	var filename string
	if opts.customFilename != "" {
		filename = opts.customFilename
	} else {
		filename = filepath.Base(path)
	}

	// prepare formdata
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(opts.portalFileFieldname, filename)
	if err != nil {
		return "", err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return "", err
	}
	err = writer.Close()
	if err != nil {
		return "", err
	}

	// prepare the request
	url := fmt.Sprintf("%s/%s", strings.TrimRight(opts.portalUrl, "/"), strings.TrimLeft(opts.portalUploadPath, "/"))
	req, err := http.NewRequest("POST", url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	if err != nil {
		return "", err
	}

	// upload the file to skynet
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	// parse the response
	body = &bytes.Buffer{}
	_, err = body.ReadFrom(resp.Body)
	if err != nil {
		return "", err
	}
	resp.Body.Close()

	if !opts.tryParseResponse {
		return body.String(), nil
	}

	var apiResponse UploadReponse
	err = json.Unmarshal(body.Bytes(), &apiResponse)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s\n", strings.TrimRight(opts.portalUrl, "/"), apiResponse.Skylink), nil
}