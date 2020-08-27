package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type (
	Resource struct {
		Name       string `json:"name"`
		Path       string `json:"path"`
		Created    string `json:"created"`
		ResourceId string `json:"resource_id"`
		Type       string `json:"type"`
		MimeType   string `json:"mime_type"`
		Embedded   struct {
			Items []Resource `json:"items"`
			Path  string     `json:"path"`
		} `json:"_embedded"`
	}
)

func apiRequest(path, method string) (*http.Response, error) {
	client := http.Client{}
	url := fmt.Sprintf("%s/%s", ynxUrl, path)
	req, _ := http.NewRequest(method, url, nil)
	req.Header.Add("Authorization", fmt.Sprintf("OAuth %s", authToken))
	return client.Do(req)
}

func createFolder(path string) error {
	_, err := apiRequest(fmt.Sprintf("resources?path=%s", path), "PUT")
	return err
}

func uploadFile(localPath, remotePath string) error {
	getUploadUrl := func(path string) (string, error) {
		res, err := apiRequest(fmt.Sprintf("resources/upload?path=%s&overwrite=true", path), "GET")
		if err != nil {
			return "", err
		}
		var resultJson struct {
			Href string `json:"href"`
		}
		err = json.NewDecoder(res.Body).Decode(&resultJson)
		if err != nil {
			return "", err
		}
		return resultJson.Href, err
	}

	data, err := os.Open(localPath)
	if err != nil {
		return err
	}
	href, err := getUploadUrl(remotePath)
	if err != nil {
		return err
	}
	//fmt.Printf("localPath: %s\n", localPath)
	//fmt.Printf("remotePath: %s href: %s\n", remotePath, href)
	defer data.Close()
	req, err := http.NewRequest("PUT", href, data)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", fmt.Sprintf("OAuth %s", authToken))

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return nil
}

func deleteFile(path string) error {
	_, err := apiRequest(fmt.Sprintf("resources?path=%s&permanently=true", path), "DELETE")
	return err
}

func getResource(path string) (*Resource, error) {
	res, err := apiRequest(fmt.Sprintf("resources?path=%s&limit=50&sort=-created", path), "GET")
	if err != nil {
		return nil, err
	}

	var result *Resource
	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
