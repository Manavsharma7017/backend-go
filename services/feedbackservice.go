package services

import (
	"backend/config"
	"backend/database"
	"backend/models"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

type QA struct {
	Number   int    `json:"number"`
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

type UserSubmission struct {
	Question   string `json:"question"`
	Answer     string `json:"answer"`
	ResponceId string `json:"responceId"`
	Domain     string `json:"domain"`
	History    []QA   `json:"history"`
}

type UserSubmissionResponse struct {
	Question               string `json:"Question"`
	Answer                 string `json:"Answer"`
	Clarity                string `json:"Clarity"`
	Tone                   string `json:"Tone"`
	Relevance              string `json:"Relevance"`
	OverallScore           string `json:"OverallScore"`
	Suggestio              string `json:"Suggestio"`
	Nextquestion           string `json:"Nextquestion"`
	NextQuestionDifficulty string `json:"NextQuestionDifficulty"`
	Explanation            string `json:"Explanation"`
}

func GetFeedbackByResponseID(responseID string, feddback *models.Feedback) (*models.Feedback, error) {
	if err := database.DB.Where("ResponseID = ?", responseID).First(feddback).Error; err != nil {
		return nil, err
	}
	return feddback, nil
}
func historymanager(sessionId string) ([]QA, error) {
	var questions []models.UserQuestion

	err := database.DB.Preload("Responses").
		Where("sessionid = ?", sessionId).
		Order("questionnumber DESC").
		Limit(5).
		Find(&questions).Error

	if err != nil {
		return nil, err
	}

	historyList := make([]QA, 0, len(questions))
	for i := len(questions) - 1; i >= 0; i-- {
		q := questions[i]
		answer := ""
		if q.Responses != nil {
			answer = q.Responses.Answer
		}

		historyList = append(historyList, QA{
			Number:   q.Questionnumber,
			Question: q.Text,
			Answer:   answer,
		})
	}

	log.Printf("📚 Final historyList size: %d", len(historyList))
	return historyList, nil
}

func CreateFeedback(feedback *models.Feedback, userrespocedata models.UserQuestionfedback) (*models.Feedback, error) {
	// Get history from DB
	history, err := historymanager(userrespocedata.SessionID)
	if err != nil {
		log.Printf("❌ Failed to fetch history: %v", err)
		return nil, err
	}
	log.Printf("📚 Fetched %d history items for session %s", len(history), userrespocedata.SessionID)

	// Build HTTP request payload
	payload := UserSubmission{
		Question:   userrespocedata.Question,
		Answer:     userrespocedata.Answer,
		ResponceId: userrespocedata.ResponceId,
		Domain:     userrespocedata.Domain,
		History:    history,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		log.Printf("❌ JSON marshaling error: %v", err)
		return nil, err
	}

	// Send POST request to FastAPI
	url := config.GetBAckedUrl()
	fullURL := fmt.Sprintf("%s/submit", url)
	resp, err := http.Post(fullURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("❌ HTTP request error: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("❌ Bad response from server: %d", resp.StatusCode)
		return nil, errors.New("server error: non-200 response")
	}

	// Parse JSON response
	var feedbackResp UserSubmissionResponse
	if err := json.NewDecoder(resp.Body).Decode(&feedbackResp); err != nil {
		log.Printf("❌ Failed to decode response: %v", err)
		return nil, err
	}

	// Populate feedback model
	feedback.Clarity = feedbackResp.Clarity
	feedback.Tone = feedbackResp.Tone
	feedback.Relevance = feedbackResp.Relevance
	feedback.OverallScore = feedbackResp.OverallScore
	feedback.Suggestion = feedbackResp.Suggestio
	feedback.NextQuestion = feedbackResp.Nextquestion
	feedback.NextQuestionDifficuilty = feedbackResp.NextQuestionDifficulty
	feedback.Explanation = feedbackResp.Explanation
	feedback.ResponseID = userrespocedata.ResponceId

	// Save feedback to DB
	if err := database.DB.Create(feedback).Error; err != nil {
		return nil, errors.New("failed to create feedback in the database")
	}

	// Link feedback to response
	var response models.Response
	if err := database.DB.Preload("Feedback").First(&response, "ResponseID = ?", userrespocedata.ResponceId).Error; err != nil {
		return nil, errors.New("response not found")
	}
	response.Feedback = *feedback

	if err := database.DB.Save(&response).Error; err != nil {
		return nil, errors.New("failed to save response with feedback")
	}

	return feedback, nil
}
