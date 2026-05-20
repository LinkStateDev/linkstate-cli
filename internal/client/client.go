package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	Server string
	Token  string
	HTTP   *http.Client
}

func New(server, token string) *Client {
	return &Client{
		Server: server,
		Token:  token,
		HTTP:   &http.Client{Timeout: 15 * time.Second},
	}
}

type LoginResponse struct {
	Token string `json:"token"`
}

func (c *Client) Login(email, password string) (string, error) {
	body := map[string]string{"email": email, "password": password}
	resp, err := c.do("POST", "/api/auth/login", body, "")
	if err != nil {
		return "", fmt.Errorf("login request: %w", err)
	}
	var r LoginResponse
	if err := decode(resp, &r); err != nil {
		return "", fmt.Errorf("login decode: %w", err)
	}
	return r.Token, nil
}

type Course struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (c *Client) ListCourses() ([]Course, error) {
	resp, err := c.do("GET", "/api/courses", nil, "")
	if err != nil {
		return nil, fmt.Errorf("list courses: %w", err)
	}
	var courses []Course
	if err := decode(resp, &courses); err != nil {
		return nil, fmt.Errorf("list courses decode: %w", err)
	}
	return courses, nil
}

type Lesson struct {
	ID        int    `json:"id"`
	CourseID  int    `json:"course_id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	IsFree    bool   `json:"is_free"`
	SortOrder int    `json:"sort_order"`
	Status    string `json:"status,omitempty"`
}

type CourseDetail struct {
	ID          int      `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Lessons     []Lesson `json:"lessons"`
}

func (c *Client) GetCourse(id int) (*CourseDetail, error) {
	resp, err := c.do("GET", fmt.Sprintf("/api/courses/%d", id), nil, c.Token)
	if err != nil {
		return nil, fmt.Errorf("get course: %w", err)
	}
	var d CourseDetail
	if err := decode(resp, &d); err != nil {
		return nil, fmt.Errorf("get course decode: %w", err)
	}
	return &d, nil
}

func (c *Client) GetLesson(id int) (*Lesson, error) {
	resp, err := c.do("GET", fmt.Sprintf("/api/lessons/%d", id), nil, c.Token)
	if err != nil {
		return nil, fmt.Errorf("get lesson: %w", err)
	}
	var l Lesson
	if err := decode(resp, &l); err != nil {
		return nil, fmt.Errorf("get lesson decode: %w", err)
	}
	return &l, nil
}

type Task struct {
	ID         int    `json:"id"`
	LessonID   int    `json:"lesson_id"`
	Title      string `json:"title"`
	TaskType   string `json:"task_type"`
	Template   string `json:"template"`
	TestConfig string `json:"test_config"`
	SortOrder  int    `json:"sort_order"`
}

func (c *Client) GetTask(id int) (*Task, error) {
	resp, err := c.do("GET", fmt.Sprintf("/api/tasks/%d", id), nil, c.Token)
	if err != nil {
		return nil, fmt.Errorf("get task: %w", err)
	}
	var t Task
	if err := decode(resp, &t); err != nil {
		return nil, fmt.Errorf("get task decode: %w", err)
	}
	return &t, nil
}

type SubmitResponse struct {
	TaskID          int    `json:"task_id"`
	Status          string `json:"status"`
	LessonCompleted bool   `json:"lesson_completed,omitempty"`
	NextLessonID    *int   `json:"next_lesson_id,omitempty"`
}

func (c *Client) Submit(taskID int, status string) (*SubmitResponse, error) {
	body := map[string]string{"status": status}
	resp, err := c.do("POST", fmt.Sprintf("/api/tasks/%d/submit", taskID), body, c.Token)
	if err != nil {
		return nil, fmt.Errorf("submit: %w", err)
	}
	var r SubmitResponse
	if err := decode(resp, &r); err != nil {
		return nil, fmt.Errorf("submit decode: %w", err)
	}
	return &r, nil
}

type ProgressItem struct {
	LessonID    int     `json:"lesson_id"`
	Status      string  `json:"status"`
	CompletedAt *string `json:"completed_at,omitempty"`
}

func (c *Client) GetProgress() ([]ProgressItem, error) {
	resp, err := c.do("GET", "/api/progress", nil, c.Token)
	if err != nil {
		return nil, fmt.Errorf("progress: %w", err)
	}
	var items []ProgressItem
	if err := decode(resp, &items); err != nil {
		return nil, fmt.Errorf("progress decode: %w", err)
	}
	return items, nil
}

func (c *Client) do(method, path string, body any, token string) (*http.Response, error) {
	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewReader(data)
	}
	req, err := http.NewRequest(method, c.Server+path, bodyReader)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

type apiError struct {
	Error string `json:"error"`
}

func decode(resp *http.Response, v any) error {
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		var e apiError
		json.NewDecoder(resp.Body).Decode(&e)
		if e.Error != "" {
			return fmt.Errorf("%s (HTTP %d)", e.Error, resp.StatusCode)
		}
		return fmt.Errorf("HTTP %d", resp.StatusCode)
	}
	return json.NewDecoder(resp.Body).Decode(v)
}
