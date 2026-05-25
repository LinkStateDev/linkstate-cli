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
	return &Client{Server: server, Token: token, HTTP: &http.Client{Timeout: 15 * time.Second}}
}

type LoginResponse struct {
	Token string `json:"token"`
}

func (c *Client) Login(email, password string) (string, error) {
	body := map[string]string{"email": email, "password": password}
	resp, err := c.do("POST", "/api/auth/login", body, "")
	if err != nil { return "", fmt.Errorf("login: %w", err) }
	var r LoginResponse
	if err := decode(resp, &r); err != nil { return "", err }
	return r.Token, nil
}

type Lesson struct {
	ID         int    `json:"id"`
	CourseSlug string `json:"course_slug"`
	Slug       string `json:"slug"`
	Title      string `json:"title"`
	Template   string `json:"template"`
	TestConfig string `json:"test_config"`
}

func (c *Client) GetLessonBySlug(slug string) (*Lesson, error) {
	resp, err := c.do("GET", "/api/lessons/slug/"+slug, nil, c.Token)
	if err != nil { return nil, err }
	var l Lesson
	if err := decode(resp, &l); err != nil { return nil, err }
	return &l, nil
}

func (c *Client) GetHint(slug string, level int) (map[string]any, error) {
	resp, err := c.do("GET", fmt.Sprintf("/api/lessons/slug/%s/hints/%d", slug, level), nil, c.Token)
	if err != nil { return nil, err }
	var r map[string]any
	if err := decode(resp, &r); err != nil { return nil, err }
	return r, nil
}

type SubmitResponse struct {
	LessonCompleted bool   `json:"lesson_completed,omitempty"`
	NextLessonID    *int   `json:"next_lesson_id,omitempty"`
}

func (c *Client) Submit(lessonID int, status string) (*SubmitResponse, error) {
	body := map[string]string{"status": status}
	resp, err := c.do("POST", fmt.Sprintf("/api/lessons/%d/submit", lessonID), body, c.Token)
	if err != nil { return nil, err }
	var r SubmitResponse
	if err := decode(resp, &r); err != nil { return nil, err }
	return &r, nil
}

type ProgressItem struct {
	LessonID    int     `json:"lesson_id"`
	Status      string  `json:"status"`
	CompletedAt *string `json:"completed_at,omitempty"`
}

func (c *Client) GetProgress() ([]ProgressItem, error) {
	resp, err := c.do("GET", "/api/progress", nil, c.Token)
	if err != nil { return nil, err }
	var items []ProgressItem
	if err := decode(resp, &items); err != nil { return nil, err }
	return items, nil
}

func (c *Client) do(method, path string, body any, token string) (*http.Response, error) {
	var bodyReader io.Reader
	if body != nil {
		data, _ := json.Marshal(body)
		bodyReader = bytes.NewReader(data)
	}
	req, _ := http.NewRequest(method, c.Server+path, bodyReader)
	if body != nil { req.Header.Set("Content-Type", "application/json") }
	if token != "" { req.Header.Set("Authorization", "Bearer "+token) }
	return c.HTTP.Do(req)
}

type apiError struct{ Error string `json:"error"` }

func decode(resp *http.Response, v any) error {
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		var e apiError
		json.NewDecoder(resp.Body).Decode(&e)
		if e.Error != "" { return fmt.Errorf("%s", e.Error) }
		return fmt.Errorf("HTTP %d", resp.StatusCode)
	}
	return json.NewDecoder(resp.Body).Decode(v)
}
