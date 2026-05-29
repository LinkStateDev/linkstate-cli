package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

var ErrUnauthorized = errors.New("unauthorized")

type Client struct {
	Server string
	Token  string
	HTTP   *http.Client
}

func New(server, token string) *Client {
	return &Client{Server: server, Token: token, HTTP: &http.Client{Timeout: 15 * time.Second}}
}

type LoginResponse struct {
	Token string `json:"token"`
}

func (c *Client) Login(email, password string) (string, error) {
	body := map[string]string{"email": email, "password": password}
	resp, err := c.do("POST", "/api/auth/login", body, "")
	if err != nil {
		return "", fmt.Errorf("login: %w", err)
	}
	var r LoginResponse
	if err := decode(resp, &r); err != nil {
		return "", err
	}
	return r.Token, nil
}

type Task struct {
	ID        int    `json:"id"`
	LessonID  int    `json:"lesson_id"`
	Slug      string `json:"slug"`
	Title     string `json:"title"`
	SortOrder int    `json:"sort_order"`
}

type Lesson struct {
	ID        int    `json:"id"`
	TrackSlug string `json:"track_slug"`
	Slug      string `json:"slug"`
	Title     string `json:"title"`
	Template  string `json:"template"`
}

type lessonResponse struct {
	Lesson Lesson `json:"lesson"`
	Tasks  []Task `json:"tasks"`
}

func (c *Client) GetLesson(id int) (*lessonResponse, error) {
	resp, err := c.do("GET", fmt.Sprintf("/api/lessons/%d", id), nil, c.Token)
	if err != nil {
		return nil, err
	}
	var lr lessonResponse
	if err := decode(resp, &lr); err != nil {
		return nil, err
	}
	return &lr, nil
}

func (c *Client) GetLessonBySlug(slug string) (*Lesson, error) {
	resp, err := c.do("GET", "/api/s/"+slug, nil, c.Token)
	if err != nil {
		return nil, err
	}
	var lr lessonResponse
	if err := decode(resp, &lr); err != nil {
		return nil, err
	}
	return &lr.Lesson, nil
}

type HintResponse struct {
	Level int    `json:"level"`
	Total int    `json:"total"`
	Hint  string `json:"hint"`
}

func (c *Client) GetHint(taskID int, level int) (*HintResponse, error) {
	resp, err := c.do("GET", fmt.Sprintf("/api/tasks/%d/hints/%d", taskID, level), nil, c.Token)
	if err != nil {
		return nil, err
	}
	var r HintResponse
	if err := decode(resp, &r); err != nil {
		return nil, err
	}
	return &r, nil
}

type SubmitResponse struct {
	TaskID          int     `json:"task_id"`
	LessonCompleted bool    `json:"lesson_completed,omitempty"`
	NextTaskID      *int    `json:"next_task_id,omitempty"`
	NextTaskSlug    *string `json:"next_task_slug,omitempty"`
	NextLessonID    *int    `json:"next_lesson_id,omitempty"`
	NextLessonSlug  *string `json:"next_lesson_slug,omitempty"`
}

func (c *Client) Submit(taskID int, status string) (*SubmitResponse, error) {
	body := map[string]string{"status": status}
	resp, err := c.do("POST", fmt.Sprintf("/api/tasks/%d/submit", taskID), body, c.Token)
	if err != nil {
		return nil, err
	}
	var r SubmitResponse
	if err := decode(resp, &r); err != nil {
		return nil, err
	}
	return &r, nil
}

type StartResponse struct {
	Task            Task   `json:"task"`
	Lesson          Lesson `json:"lesson"`
	ModuleCompleted bool   `json:"module_completed"`
}

func (c *Client) GetStartLesson(trackSlug, moduleSlug string) (*StartResponse, error) {
	resp, err := c.do("GET", fmt.Sprintf("/api/start/%s/%s", trackSlug, moduleSlug), nil, c.Token)
	if err != nil {
		return nil, err
	}
	var r StartResponse
	if err := decode(resp, &r); err != nil {
		return nil, err
	}
	return &r, nil
}

type ProgressItem struct {
	TaskID       int     `json:"task_id"`
	TaskSlug     string  `json:"task_slug"`
	TaskTitle    string  `json:"task_title"`
	LessonID     int     `json:"lesson_id"`
	LessonSlug   string  `json:"lesson_slug"`
	LessonTitle  string  `json:"lesson_title"`
	TrackSlug    string  `json:"track_slug"`
	TrackTitle   string  `json:"track_title"`
	Status       string  `json:"status"`
	CompletedAt  *string `json:"completed_at,omitempty"`
}

func (c *Client) GetProgress() ([]ProgressItem, error) {
	resp, err := c.do("GET", "/api/progress", nil, c.Token)
	if err != nil {
		return nil, err
	}
	var items []ProgressItem
	if err := decode(resp, &items); err != nil {
		return nil, err
	}
	return items, nil
}

func (c *Client) do(method, path string, body any, token string) (*http.Response, error) {
	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshal body: %w", err)
		}
		bodyReader = bytes.NewReader(data)
	}
	req, err := http.NewRequest(method, c.Server+path, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	return c.HTTP.Do(req)
}

type apiError struct {
	Error string `json:"error"`
}

func decode(resp *http.Response, v any) error {
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		var e apiError
		msg := fmt.Sprintf("HTTP %d", resp.StatusCode)
		if err := json.NewDecoder(resp.Body).Decode(&e); err == nil && e.Error != "" {
			msg = e.Error
		}
		if resp.StatusCode == http.StatusUnauthorized {
			return fmt.Errorf("%s: %w", msg, ErrUnauthorized)
		}
		return fmt.Errorf("%s", msg)
	}
	return json.NewDecoder(resp.Body).Decode(v)
}
