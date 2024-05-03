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

	query := fmt.Sprintf("Create 8 total problems for %s graders for the subject %s. Make sure you give me 2 word problems maximum. Also make sure the problems start with what a number and type of problem it is so can parese it easy for example (1.word problem: \"the wordproblem\"). Can you not add special characters in the response.",
		grade,
		subject)
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
				"content": query,
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
