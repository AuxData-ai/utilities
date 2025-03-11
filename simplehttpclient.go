package utilities

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type SimpleHttpClient struct {
	Headers     map[string]string
	Body        string
	Url         string
	Method      string
	ContentType string
}

func (client *SimpleHttpClient) AddHeader(key string, value string) {

	if client.Headers == nil {
		client.Headers = make(map[string]string)
	}

	client.Headers[key] = value
}

func (client *SimpleHttpClient) AddBearerAuthentificationToken(token string) {

	if client.Headers == nil {
		client.Headers = make(map[string]string)
	}

	client.Headers["Authorization"] = fmt.Sprintf("Bearer %s", token)
}

func (client *SimpleHttpClient) AddObjectAsBody(object interface{}) {

	data, err := json.Marshal(object)
	if err != nil {
		log.Println(err)
	}

	client.Body = string(data)
}

func (client SimpleHttpClient) Execute() (string, error) {

	// Create a new HTTP Client
	httpClient := &http.Client{}

	// Create a new HTTP Request
	req, err := http.NewRequest(client.Method, client.Url, nil)

	if err != nil {
		log.Println(err)
		return "", err
	}

	// Set the headers for the request
	for key, value := range client.Headers {
		req.Header.Set(key, value)
	}

	if client.ContentType == "" {
		req.Header.Set("Content-Type", "application/json")
	} else {
		req.Header.Set("Content-Type", client.ContentType)
	}

	req.Body = io.NopCloser(strings.NewReader(client.Body))
	httpClient.Timeout = time.Duration(time.Duration.Minutes(2))

	resp, err := httpClient.Do(req)
	if err != nil {
		log.Println(err)
		return "", err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return "", err
	}

	return string(body), nil
}

func (client SimpleHttpClient) ExecuteStream() (string, error) {

	// Create a new HTTP Client
	httpClient := &http.Client{}

	// Create a new HTTP Request
	req, err := http.NewRequest(client.Method, client.Url, nil)

	if err != nil {
		log.Println(err)
		return "", err
	}

	// Set the headers for the request
	for key, value := range client.Headers {
		req.Header.Set(key, value)
	}

	if client.ContentType == "" {
		req.Header.Set("Content-Type", "application/json")
	} else {
		req.Header.Set("Content-Type", client.ContentType)
	}

	req.Body = io.NopCloser(strings.NewReader(client.Body))
	httpClient.Timeout = time.Duration(time.Duration.Minutes(2))

	resp, err := httpClient.Do(req)
	if err != nil {
		log.Println(err)
		return "", err
	}

	defer resp.Body.Close()

	var body []byte
	reader := bufio.NewReader(resp.Body)

	for {
		line, err := reader.ReadBytes('\n')

		if err == io.EOF || string(line) == "data: [DONE]\n" {
			break
		}

		if err != nil {
			log.Println(err)
			return "", err
		}

		var obj map[string]interface{}
		index := strings.Index(string(line), "{")

		if index == -1 {
			continue
		}

		// We create a new []byte, which contains only the part beginning at "{"
		content := line[index:]
		err = json.Unmarshal(content, &obj)

		if err != nil {
			log.Println(err)
			break
		}

		choices := obj["choices"].([]interface{})
		data := choices[0].(map[string]interface{})
		resultContent := data["delta"].(map[string]interface{})

		if resultContent["content"] != nil {
			body = append(body, []byte(resultContent["content"].(string))...)
		}
	}

	return string(body), nil
}

func (client SimpleHttpClient) ExecuteObjectResult() (map[string]interface{}, error) {

	response, err := client.Execute()

	if err != nil {
		log.Println(err)
		return nil, err
	}

	response, err = GetJsonPartOfResponse(response)

	if err != nil {
		log.Println(err)
		return nil, errors.New(response)
	}

	var responseObject map[string]interface{}
	err = json.Unmarshal([]byte(response), &responseObject)

	if err != nil {
		log.Println(err)
		log.Println(response)
		return nil, errors.New(response)
	}

	return responseObject, nil
}

func (client SimpleHttpClient) ExecuteObjectResultAsStream() (map[string]interface{}, error) {

	result, err := client.ExecuteStream()

	if err != nil {
		log.Println(err)
		return nil, err
	}

	var obj map[string]interface{}
	err = json.Unmarshal([]byte(result), &obj)

	if err != nil {
		log.Println(err)
		obj = make(map[string]interface{})
		obj["result"] = result
		return obj, nil
	}

	return obj, nil
}
