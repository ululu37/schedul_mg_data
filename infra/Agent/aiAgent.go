// NewAiAgent creates a new Agent with default API key and URL

package aiAgent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

//const apiKey = "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" // ใส่ API Key ของคุณที่นี่
//const apiURL = "https://openrouter.ai/api/v1/chat/completions"                             // ตัวอย่าง URL, เปลี่ยนตามเอกสารจริง

type Agent struct {
	ApiKey string
	ApiURL string
	Model  string
}

func NewAiAgent(apiKey string, apiURL string, model string) *Agent {
	return &Agent{
		ApiKey: apiKey,
		ApiURL: apiURL,
		Model:  model,
	}
}

func (g *Agent) Chat(messages []Message) (*ChatCompletionResponse, error) {
	reqBody := GPTMini4Request{
		Model:    g.Model,
		Messages: messages,
	}
	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	// Log payload size for debugging
	sizeMB := float64(len(bodyBytes)) / 1024 / 1024
	fmt.Printf("AI Request Payload Size: %.2f MB\n", sizeMB)

	req, err := http.NewRequest("POST", g.ApiURL, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+g.ApiKey)

	client := &http.Client{
		Timeout: 5 * time.Minute,
	}
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("AI API error (status %d): %s", resp.StatusCode, string(respBody))
	}
	//fmt.Println(string(respBody))

	var result ChatCompletionResponse

	errJ := json.Unmarshal(respBody, &result)
	if errJ != nil {
		return nil, errJ
	}

	//fmt.Println(result.Choices[0].Message.Content)
	return &result, nil
}
