// NewAiAgent creates a new Agent with default API key and URL

package aiAgent

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

//const apiKey = "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" // ใส่ API Key ของคุณที่นี่
//const apiURL = "https://openrouter.ai/api/v1/chat/completions"                             // ตัวอย่าง URL, เปลี่ยนตามเอกสารจริง

type Agent struct {
	ApiKey string
	ApiURL string
}

func NewAiAgent(apiKey string, apiURL string) *Agent {
	return &Agent{
		ApiKey: apiKey,
		ApiURL: apiURL,
	}
}

func (g *Agent) Chat(messages []Message) (*ChatCompletionResponse, error) {
	reqBody := GPTMini4Request{
		Model:    "google/gemini-3-flash-preview",
		Messages: messages,
	}
	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", g.ApiURL, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+g.ApiKey)

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(respBody))

	var result ChatCompletionResponse

	errJ := json.Unmarshal(respBody, &result)
	if errJ != nil {
		return nil, errJ
	}

	//fmt.Println(result.Choices[0].Message.Content)
	return &result, nil
}
