package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

// 데이터 모델
type User struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type Issue struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	User        *User     `json:"user,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// 요청/응답 구조체
type CreateIssueRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	UserID      *uint  `json:"userId,omitempty"`
}

type UpdateIssueRequest struct {
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	Status      *string `json:"status,omitempty"`
	UserID      *uint   `json:"userId,omitempty"`
}

type IssueListResponse struct {
	Issues []Issue `json:"issues"`
}

type ErrorResponse struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}

// 전역 변수 (실제 환경에서는 데이터베이스 사용)
var (
	users      []User
	issues     []Issue
	nextUserID uint = 4
	nextIssueID uint = 1
)

// 유효한 상태값
var validStatuses = map[string]bool{
	"PENDING":     true,
	"IN_PROGRESS": true,
	"COMPLETED":   true,
	"CANCELLED":   true,
}

// 초기 데이터 설정
func init() {
	users = []User{
		{ID: 1, Name: "김개발"},
		{ID: 2, Name: "이디자인"},
		{ID: 3, Name: "박기획"},
	}
}

// 유틸리티 함수들
func findUserByID(id uint) *User {
	for _, user := range users {
		if user.ID == id {
			return &user
		}
	}
	return nil
}

func findIssueByID(id uint) *Issue {
	for i, issue := range issues {
		if issue.ID == id {
			return &issues[i]
		}
	}
	return nil
}

func validateStatus(status string) bool {
	return validStatuses[status]
}

func sendErrorResponse(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message, Code: code})
}

func sendJSONResponse(w http.ResponseWriter, data interface{}, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(data)
}

// API 핸들러들
func createIssue(w http.ResponseWriter, r *http.Request) {
	var req CreateIssueRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendErrorResponse(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// 필수 필드 검증
	if strings.TrimSpace(req.Title) == "" {
		sendErrorResponse(w, "Title is required", http.StatusBadRequest)
		return
	}

	// 사용자 검증 (존재하는 경우)
	var assignedUser *User
	status := "PENDING"
	
	if req.UserID != nil {
		user := findUserByID(*req.UserID)
		if user == nil {
			sendErrorResponse(w, "User not found", http.StatusBadRequest)
			return
		}
		assignedUser = user
		status = "IN_PROGRESS"
	}

	// 새 이슈 생성
	now := time.Now()
	issue := Issue{
		ID:          nextIssueID,
		Title:       strings.TrimSpace(req.Title),
		Description: strings.TrimSpace(req.Description),
		Status:      status,
		User:        assignedUser,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	issues = append(issues, issue)
	nextIssueID++

	sendJSONResponse(w, issue, http.StatusCreated)
}

func getIssues(w http.ResponseWriter, r *http.Request) {
	statusFilter := r.URL.Query().Get("status")
	
	// 상태 필터 유효성 검증
	if statusFilter != "" && !validateStatus(statusFilter) {
		sendErrorResponse(w, "Invalid status value", http.StatusBadRequest)
		return
	}

	var filteredIssues []Issue
	for _, issue := range issues {
		if statusFilter == "" || issue.Status == statusFilter {
			filteredIssues = append(filteredIssues, issue)
		}
	}

	response := IssueListResponse{Issues: filteredIssues}
	sendJSONResponse(w, response, http.StatusOK)
}

func getIssue(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		sendErrorResponse(w, "Invalid issue ID", http.StatusBadRequest)
		return
	}

	issue := findIssueByID(uint(id))
	if issue == nil {
		sendErrorResponse(w, "Issue not found", http.StatusNotFound)
		return
	}

	sendJSONResponse(w, *issue, http.StatusOK)
}

func updateIssue(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		sendErrorResponse(w, "Invalid issue ID", http.StatusBadRequest)
		return
	}

	issue := findIssueByID(uint(id))
	if issue == nil {
		sendErrorResponse(w, "Issue not found", http.StatusNotFound)
		return
	}

	// COMPLETED 또는 CANCELLED 상태에서는 업데이트 불가
	if issue.Status == "COMPLETED" || issue.Status == "CANCELLED" {
		sendErrorResponse(w, "Cannot update completed or cancelled issue", http.StatusBadRequest)
		return
	}

	var req UpdateIssueRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendErrorResponse(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// 업데이트 전 상태 백업
	originalStatus := issue.Status
	originalUser := issue.User

	// 필드별 업데이트
	if req.Title != nil {
		trimmed := strings.TrimSpace(*req.Title)
		if trimmed == "" {
			sendErrorResponse(w, "Title cannot be empty", http.StatusBadRequest)
			return
		}
		issue.Title = trimmed
	}

	if req.Description != nil {
		issue.Description = strings.TrimSpace(*req.Description)
	}

	// 사용자 업데이트 처리
	if req.UserID != nil {
		user := findUserByID(*req.UserID)
		if user == nil {
			sendErrorResponse(w, "User not found", http.StatusBadRequest)
			return
		}
		issue.User = user
	}

	// 상태 업데이트 처리
	if req.Status != nil {
		if !validateStatus(*req.Status) {
			sendErrorResponse(w, "Invalid status value", http.StatusBadRequest)
			return
		}
		issue.Status = *req.Status
	}

	// 비즈니스 규칙 검증
	// 1. 담당자 없이 PENDING, CANCELLED 이외 상태로 변경 불가
	if issue.User == nil && issue.Status != "PENDING" && issue.Status != "CANCELLED" {
		sendErrorResponse(w, "Cannot set status to IN_PROGRESS or COMPLETED without assignee", http.StatusBadRequest)
		return
	}

	// 2. PENDING 상태에서 담당자 할당 시 상태 변경
	if originalStatus == "PENDING" && req.UserID != nil && req.Status == nil {
		issue.Status = "IN_PROGRESS"
	}

	// 3. 담당자 제거 시 상태를 PENDING으로 변경
	if originalUser != nil && req.UserID != nil && *req.UserID == 0 {
		issue.User = nil
		issue.Status = "PENDING"
	}

	// 특별한 경우: userId가 명시적으로 null로 설정된 경우 (JSON에서 null 값 처리)
	// 이는 실제로는 요청에서 "userId": null로 전달되는 경우를 처리하기 위한 것입니다.
	// 현재 구현에서는 포인터를 사용하여 이를 구분하기 어려우므로,
	// 프론트엔드와의 협의가 필요한 부분입니다.

	issue.UpdatedAt = time.Now()

	sendJSONResponse(w, *issue, http.StatusOK)
}

func main() {
	router := mux.NewRouter()
	
	// API 라우트
	router.HandleFunc("/issue", createIssue).Methods("POST")
	router.HandleFunc("/issues", getIssues).Methods("GET")
	router.HandleFunc("/issue/{id}", getIssue).Methods("GET")
	router.HandleFunc("/issue/{id}", updateIssue).Methods("PATCH")
	
	// 헬스 체크
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	fmt.Println("Server starting on port 8080...")
	if err := http.ListenAndServe(":8080", router); err != nil {
		fmt.Printf("Server failed to start: %v\n", err)
	}
}