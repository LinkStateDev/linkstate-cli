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

// ErrUnauthorized is returned (wrapped) by Client methods when the server
// responds with HTTP 401. Callers detect it via errors.Is so they can clear
// the local token and prompt the user to re-authenticate.
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
	if err != nil {
		return nil, err
	}
	var l Lesson
	if err := decode(resp, &l); err != nil {
		return nil, err
	}
	return &l, nil
}

type HintResponse struct {
	Level int    `json:"level"`
	Total int    `json:"total"`
	Hint  string `json:"hint"`
}

func (c *Client) GetHint(slug string, level int) (*HintResponse, error) {
	resp, err := c.do("GET", fmt.Sprintf("/api/lessons/slug/%s/hints/%d", slug, level), nil, c.Token)
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
	LessonCompleted bool    `json:"lesson_completed,omitempty"`
	NextLessonID    *int    `json:"next_lesson_id,omitempty"`
	NextLessonSlug  *string `json:"next_lesson_slug,omitempty"`
}

func (c *Client) Submit(lessonID int, status string) (*SubmitResponse, error) {
	body := map[string]string{"status": status}
	resp, err := c.do("POST", fmt.Sprintf("/api/lessons/%d/submit", lessonID), body, c.Token)
	if err != nil {
		return nil, err
	}
	var r SubmitResponse
	if err := decode(resp, &r); err != nil {
		return nil, err
	}
	return &r, nil
}

type ProgressItem struct {
	LessonID    int     `json:"lesson_id"`
	LessonSlug  string  `json:"lesson_slug"`
	LessonTitle string  `json:"lesson_title"`
	CourseSlug  string  `json:"course_slug"`
	CourseTitle string  `json:"course_title"`
	Status      string  `json:"status"`
	CompletedAt *string `json:"completed_at,omitempty"`
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
