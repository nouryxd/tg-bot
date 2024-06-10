package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// WolframAlphaShort queries the WolframAlpha Short Answers API with
// the given query and returns the result.
func WolframAlphaShort(query, appid string) (string, error) {
	escaped := url.QueryEscape(query)
	url := fmt.Sprintf("http://api.wolframalpha.com/v1/result?appid=%s&i=%s", appid, escaped)

	resp, err := http.Get(url)
	if err != nil {
		return "", ErrInternalServerError
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", ErrInternalServerError
	}

	reply := string(body)
	return reply, nil
}

type waConversationResponse struct {
	Result         string `json:"result"`
	ConversationID string `json:"conversationID"`
	Host           string `json:"host"`
	S              string `json:"s"`
}

// WolframAlphaConv queries the WolframAlpha Conversational API with
// the given query and returns the result.
func WolframAlphaConv(query, appid string) (string, error) {
	escaped := url.QueryEscape(query)
	url := fmt.Sprintf("http://api.wolframalpha.com/v1/conversation.jsp?appid=%s&i=%s", appid, escaped)

	response, err := http.Get(url)
	if err != nil {
		return "", ErrInternalServerError
	}
	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		return "", ErrInternalServerError
	}
	var responseObject waConversationResponse
	if err = json.Unmarshal(responseData, &responseObject); err != nil {
		return "", ErrInternalServerError
	}

	return responseObject.Result, nil
}
