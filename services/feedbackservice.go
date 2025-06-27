package services

import (
	"backend/config"
	"backend/database"
	"backend/models"
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net"
	"net/http"
	"time"
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

	url := config.GetBAckedUrl()

	// Create custom HTTP client with longer timeout
	client := &http.Client{
		Timeout: 60 * time.Second, // Full request timeout
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second, // TCP connection timeout
				KeepAlive: 30 * time.Second,
			}).DialContext,
			TLSHandshakeTimeout: 10 * time.Second,
		},
	}

	var resp *http.Response
	const maxRetries = 5

	for attempt := 1; attempt <= maxRetries; attempt++ {
		resp, err = client.Post(url, "application/json", bytes.NewBuffer(jsonData))
		if err == nil {
			break // success
		}
		log.Printf("⚠️ HTTP request failed on attempt %d: %v", attempt, err)
		time.Sleep(5 * time.Second)
	}

	if err != nil {
		log.Printf("❌ HTTP request failed after retries: %v", err)
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

// package services

// import (
// 	pb "backend/common"
// 	"backend/database"
// 	"backend/models"
// 	"context"
// 	"errors"
// 	"log"

// 	"google.golang.org/grpc"

// 	"google.golang.org/grpc/credentials/insecure"
// )

// type history struct {
// 	number   int
// 	question string
// 	answer   string
// }

// func GetFeedbackByResponseID(responseID string, feddback *models.Feedback) (*models.Feedback, error) {
// 	if err := database.DB.Where("ResponseID = ?", responseID).First(feddback).Error; err != nil {
// 		return nil, err
// 	}
// 	return feddback, nil
// }
// func historymanager(sessionId string) ([]history, error) {
// 	var questions []models.UserQuestion

// 	err := database.DB.Preload("Responses").
// 		Where("sessionid = ?", sessionId).
// 		Order("questionnumber DESC").
// 		Limit(5).
// 		Find(&questions).Error

// 	log.Printf("📦 Loaded %d UserQuestions for session %s", len(questions), sessionId)

// 	if err != nil {
// 		return nil, err
// 	}

// 	historyList := make([]history, 0, len(questions))
// 	for i := len(questions) - 1; i >= 0; i-- {
// 		q := questions[i]

// 		ans := ""
// 		if q.Responses != nil {
// 			ans = q.Responses.Answer
// 		} else {
// 			log.Printf("⚠️ No response for question #%d: %s", q.Questionnumber, q.Text)
// 		}

// 		historyList = append(historyList, history{
// 			number:   q.Questionnumber,
// 			question: q.Text,
// 			answer:   ans,
// 		})
// 	}
// 	log.Printf("📚 Final historyList size: %d", len(historyList))
// 	return historyList, nil
// }
// func CreateFeedback(feedback *models.Feedback, userrespocedata models.UserQuestionfedback) (*models.Feedback, error) {
// 	dialer := grpc.WithTransportCredentials(insecure.NewCredentials())
// 	conn, err := grpc.Dial("llmserver-2gpa.onrender.com:50051", dialer)

// 	if err != nil {
// 		log.Printf("❌ Failed to connect to gRPC server: %v", err)
// 		return nil, err
// 	}
// 	defer conn.Close()

// 	client := pb.NewUserSubmittionServiceClient(conn)

// 	// Fetch history from DB
// 	history, err := historymanager(userrespocedata.SessionID)
// 	if err != nil {
// 		log.Printf("❌ Failed to fetch history: %v", err)
// 		return nil, err
// 	}
// 	log.Printf("📚 Fetched %d history items for session %s", len(history), userrespocedata.SessionID)

// 	// Build gRPC request
// 	req := &pb.UserSubmittion{
// 		Question:   userrespocedata.Question,
// 		Answer:     userrespocedata.Answer,
// 		ResponceId: userrespocedata.ResponceId,
// 		Domain:     userrespocedata.Domain,
// 		History:    make([]*pb.QA, 0, len(history)),
// 	}

// 	// Append history into the gRPC request
// 	for _, h := range history {
// 		log.Printf("📤 Adding QA: #%d | Q: %s | A: %s", h.number, h.question, h.answer)
// 		req.History = append(req.History, &pb.QA{
// 			Number:   int32(h.number),
// 			Question: h.question,
// 			Answer:   h.answer,
// 		})
// 	}
// 	log.Printf("✅ Total QA entries added to gRPC request: %d", len(req.History))

// 	// Submit request to gRPC server
// 	resp, err := client.SubmitUserSubmittion(context.Background(), req)
// 	if err != nil {
// 		log.Printf("❌ SubmitUserSubmittion gRPC error: %v", err)
// 		return nil, err
// 	}

// 	// Populate feedback fields
// 	feedback.Clarity = resp.Clarity
// 	feedback.Tone = resp.Tone
// 	feedback.Relevance = resp.Relevance
// 	feedback.OverallScore = resp.OverallScore
// 	feedback.Suggestion = resp.Suggestio
// 	feedback.NextQuestion = resp.Nextquestion
// 	feedback.NextQuestionDifficuilty = resp.NextQuestionDifficulty
// 	feedback.Explanation = resp.Explanation
// 	feedback.ResponseID = userrespocedata.ResponceId

// 	// Save feedback
// 	if err := database.DB.Create(feedback).Error; err != nil {
// 		return nil, errors.New("failed to create feedback in the database")
// 	}

// 	// Fetch the associated response and update it with the feedback
// 	var response models.Response
// 	if err := database.DB.Preload("Feedback").First(&response, "ResponseID = ?", userrespocedata.ResponceId).Error; err != nil {
// 		return nil, errors.New("response not found")
// 	}
// 	response.Feedback = *feedback

// 	if err := database.DB.Save(&response).Error; err != nil {
// 		return nil, errors.New("failed to save response with feedback")
// 	}

// 	return feedback, nil
// }
