package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/jung-kurt/gofpdf"
	"github.com/rudyjcruz831/mathSheets/util/errors"
)

type ChatResponse struct {
	Id      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Usage   Usage    `json:"usage"`
	Choices []Choice `json:"choices"`
}

type Choice struct {
	Index         int64   `json:"index"`
	Message       Message `json:"message"`
	Logprobs      bool    `json:"logprobs"`
	Finish_reason string  `json:"finish_reason"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"` // this is the important stuff
}

type Usage struct {
	Prompt_tokens     int64 `json:"prompt_tokens"`
	Completion_tokens int64 `json:"completion_tokens"`
	Total_tokens      int64 `json:"total_tokens"`
}

func createPfd(grade, subject string) (bytes.Buffer, *errors.MathSheetsError) {
	// here we are going to try to qurey OpenAI API
	// panic("CreatedPDF...")
	// c.JSON(200, gin.H{"message": "Created PDF"})
	url := "https://api.openai.com/v1/chat/completions"
	openaiAPIKey := os.Getenv("OPEN_AI_API_KEY")

	// This the query being asked to the OpenAI API

	// "Create 8 total problems for %s graders for the subject %s.
	// Make sure you give me 2 word problems maximum.
	// Also make sure the problems start with what a number
	// and type of problem it is so can parese it easy for
	// example (1.word problem: \"the wordproblem\").
	// Can you not add special characters in the response."

	query := fmt.Sprintf("Create 8 total problems for %s graders for the subject %s. Make sure you give me 2 word problems maximum. Also make sure the problems start with what a number and type of problem it is so can parese it easy for example (1.word problem: \"the wordproblem\"). Can you not add special characters in the response.", grade, subject)
	data := map[string]interface{}{
		"model": "gpt-3.5-turbo",
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": "You are a teacher assistant, skilled in explaining math problems for 1-6 graders.",
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
		mathErr := errors.NewInternalServerError("Error encoding JSON")
		return bytes.Buffer{}, mathErr
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating HTTP request:", err)
		mathErr := errors.NewInternalServerError("Error creating HTTP request")
		return bytes.Buffer{}, mathErr
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+openaiAPIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making HTTP request:", err)
		mathErr := errors.NewInternalServerError("Error making HTTP request")
		return bytes.Buffer{}, mathErr
	}
	defer resp.Body.Close()

	var result ChatResponse

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		fmt.Println("Error decoding JSON response:", err)
		mathErr := errors.NewInternalServerError("Error decoding JSON response")
		return bytes.Buffer{}, mathErr
	}

	fmt.Println("Response Result: ", result)

	if len(result.Choices) < 1 {
		mathErr := errors.NewInternalServerError("No choices found")
		return bytes.Buffer{}, mathErr
	}

	s := strings.Split(result.Choices[0].Message.Content, "\n")

	buf, err := pdf(s, subject, grade)
	if err != nil {
		return bytes.Buffer{}, errors.NewInternalServerError("Error creating PDF: " + err.Error())
	}
	return buf, nil

}

func pdf(problems []string, subject string, grade string) (bytes.Buffer, error) {

	pdf := gofpdf.New("P", "mm", "A4", "")
	// pdf.AddPage()

	title := fmt.Sprintf("%s Work Sheet for %s graders", subject, grade)
	var opt gofpdf.ImageOptions
	pdf.SetTopMargin(30)
	pdf.SetHeaderFuncMode(func() {
		pdf.ImageOptions("./main-logo-black.png", 4, 4, 30, 0, false, opt, 0, "")
		pdf.SetY(5)
		// pdf.SetX(3)
		pdf.SetFont("Times", "B", 17)
		pdf.Cell(80, 0, "")
		pdf.SetFont("Times", "B", 13)
		// pdf.SetX(4)
		pdf.CellFormat(30, 10, title, "", 0, "C", false, 0, "")
		// pdf.Ln(50)
	}, true)

	pdf.SetFooterFunc(func() {
		pdf.SetY(-15)
		pdf.SetFont("Times", "I", 8)
		pdf.CellFormat(0, 10, fmt.Sprintf("Page %d/{nb}", pdf.PageNo()),
			"", 0, "C", false, 0, "")
	})
	pdf.AliasNbPages("")
	pdf.AddPage()
	// pdf.SetTopMargin(30)
	pdf.SetFont("Times", "I", 7)
	for i := 0; i < len(problems); i++ {
		pdf.CellFormat(0, 50, problems[i], "1", 1, "LT", false, 0, "")
		pdf.CellFormat(0, 10, "", "0", 1, "LT", false, 0, "")
	}

	// create bytes to pass to frontend
	var buf bytes.Buffer
	err := pdf.Output(&buf)

	// err := pdf.OutputFileAndClose("mathworksheet.pdf")
	if err != nil {
		return buf, err
	}

	return buf, nil

	// fmt.Println(&buf)
}
