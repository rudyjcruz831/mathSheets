package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func getProblemsFromOpenAI(grade string, subject string) {
	url := "https://api.openai.com/v1/chat/completions"
	openaiAPIKey := os.Getenv("OPEN_AI_API_KEY")

	// query := fmt.Sprintf("Create 8 math problems for %s graders for the subject %s and make sure you give me 3 word problems")
	//
	data := map[string]interface{}{
		"model": "gpt-3.5-turbo",
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": "You are a teacher assistant, skilled in explaining math problems for first graders.",
			},
			{
				"role":    "user",
				"content": "Create 6 problems for math for first graders but make sure you only do 3 word problems.",
			},
		},
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating HTTP request:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+openaiAPIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making HTTP request:", err)
		return
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		fmt.Println("Error decoding JSON response:", err)
		return
	}

	fmt.Println(result)
}
