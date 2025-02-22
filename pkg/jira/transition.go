package jira

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// TransitionRequest struct holds request data for transition request.
type TransitionRequest struct {
	Transition *TransitionRequestData `json:"transition"`
	Update     *TransitionUpdateData  `json:"update"`
}
type TransitionRequestNoUpdate struct {
	Transition *TransitionRequestData `json:"transition"`
}

// TransitionRequestData is a transition request data.
type TransitionRequestData struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type TransitionUpdateData struct {
	Worklog []WorklogObj `json:"worklog"`
}
type WorklogObj struct {
	Add map[string]string `json:"add"`
}

type transitionResponse struct {
	Expand      string        `json:"expand"`
	Transitions []*Transition `json:"transitions"`
}

// Transitions fetches valid transitions for an issue using v3 version of the GET /issue/{key}/transitions endpoint.
func (c *Client) Transitions(key string) ([]*Transition, error) {
	return c.transitions(key, apiVersion3)
}

// TransitionsV2 fetches valid transitions for an issue using v2 version of the GET /issue/{key}/transitions endpoint.
func (c *Client) TransitionsV2(key string) ([]*Transition, error) {
	return c.transitions(key, apiVersion2)
}

func (c *Client) transitions(key, ver string) ([]*Transition, error) {
	path := fmt.Sprintf("/issue/%s/transitions", key)

	var (
		res *http.Response
		err error
	)

	switch ver {
	case apiVersion2:
		res, err = c.GetV2(context.Background(), path, nil)
	default:
		res, err = c.Get(context.Background(), path, nil)
	}

	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, ErrEmptyResponse
	}
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusOK {
		return nil, formatUnexpectedResponse(res)
	}

	var out transitionResponse

	err = json.NewDecoder(res.Body).Decode(&out)

	return out.Transitions, err
}

func (c *Client) TransitionNoUpdate(key string, data *TransitionRequestNoUpdate) (int, error) {
	body, err := json.Marshal(&data)
	if err != nil {
		return 0, err
	}

	path := fmt.Sprintf("/issue/%s/transitions", key)

	res, err := c.PostV2(context.Background(), path, body, Header{
		"Accept":       "application/json",
		"Content-Type": "application/json",
	})
	if err != nil {
		return 0, err
	}
	if res == nil {
		return 0, ErrEmptyResponse
	}
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusNoContent {
		return res.StatusCode, formatUnexpectedResponse(res)
	}
	return res.StatusCode, nil
}

// Transition moves issue from one state to another using POST /issue/{key}/transitions endpoint.
func (c *Client) Transition(key string, data *TransitionRequest) (int, error) {
	body, err := json.Marshal(&data)
	if err != nil {
		return 0, err
	}

	path := fmt.Sprintf("/issue/%s/transitions", key)

	res, err := c.PostV2(context.Background(), path, body, Header{
		"Accept":       "application/json",
		"Content-Type": "application/json",
	})
	if res == nil {
		return 0, ErrEmptyResponse
	}
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusNoContent {
		return res.StatusCode, formatUnexpectedResponse(res)
	}
	return res.StatusCode, nil
}
