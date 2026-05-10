package jira

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"face-recognition-fyp/domain"
)

var statusMap = map[string]int{
	"to do":       1,
	"in progress": 2,
	"in review":   3,
	"done":        4,
}

type Service struct {
	baseURL    string 
	email      string
	token      string
	projectKey string
	httpClient *http.Client
}

func NewService(baseURL, email, token, projectKey string) *Service {
	return &Service{
		baseURL:    strings.TrimRight(baseURL, "/"),
		email:      email,
		token:      token,
		projectKey: projectKey,
		httpClient: &http.Client{},
	}
}

type jiraSearchResponse struct {
	Issues []jiraIssue `json:"issues"`
}

type jiraIssue struct {
	Key    string     `json:"key"`
	Fields jiraFields `json:"fields"`
}

type jiraFields struct {
	Summary string      `json:"summary"`
	Status  jiraStatus  `json:"status"`
	Parent  *jiraParent `json:"parent"`
}

type jiraStatus struct {
	Name string `json:"name"`
}

type jiraParent struct {
	Key    string          `json:"key"`
	Fields jiraParentFields `json:"fields"`
}

type jiraParentFields struct {
	Summary string `json:"summary"`
}


func (s *Service) GetTasksForUser(ctx context.Context, email string) ([]domain.JiraTask, error) {
	jql := fmt.Sprintf(`project = %s AND assignee = "%s"`, s.projectKey, email)

	endpoint := fmt.Sprintf("%s/rest/api/3/search?jql=%s&maxResults=100&fields=summary,status,parent",
		s.baseURL, url.QueryEscape(jql),
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(s.email, s.token)
	req.Header.Set("Accept", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("jira request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("jira returned status %d", resp.StatusCode)
	}

	var result jiraSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode jira response: %w", err)
	}

	tasks := make([]domain.JiraTask, 0, len(result.Issues))
	for _, issue := range result.Issues {
		statusName := strings.ToLower(strings.TrimSpace(issue.Fields.Status.Name))
		statusID, ok := statusMap[statusName]
		if !ok {
			statusID = 1 // Default to "To Do" if unknown, same as Python
		}

		storyJira := ""
		if issue.Fields.Parent != nil {
			storyJira = issue.Fields.Parent.Key
		}

		tasks = append(tasks, domain.JiraTask{
			TicketKey: issue.Key,
			StoryJira: storyJira,
			Status:    statusID,
		})
	}

	return tasks, nil
}
