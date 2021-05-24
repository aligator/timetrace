package clockify

import (
	"bytes"
	"encoding/json"
	"github.com/dominikbraun/timetrace/clockify/clockify/model"
	"io"
	"net/http"
	"net/url"
	"path"
)

type API struct {
	Config *Config

	url    *url.URL
	client http.Client
}

func (a *API) createRequest(method string, requestPath string, body io.Reader) (*http.Request, error) {
	if a.url == nil {
		var err error
		if a.Config.Endpoint == "" {
			a.Config.Endpoint = "https://api.clockify.me/api/v1"
		}
		a.url, err = url.Parse(a.Config.Endpoint)
		if err != nil {
			return nil, err
		}
	}

	// Copy the cached URL to modify it.
	reqURL := *a.url
	reqURL.Path = path.Join(reqURL.Path, requestPath)
	req, err := http.NewRequest(method, reqURL.String(), body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Api-Key", a.Config.APIKey)
	return req, nil
}

func (a *API) Projects() ([]model.Project, error) {
	workspaces, err := a.Workspaces()
	if err != nil {
		return nil, err
	}

	var result []model.Project

	for _, workspace := range workspaces {
		req, err := a.createRequest("GET", path.Join("workspaces", workspace.ID, "projects"), nil)
		if err != nil {
			return nil, err
		}

		q := req.URL.Query()
		q.Add("archived", "false")
		req.URL.RawQuery = q.Encode()

		res, err := a.client.Do(req)
		if err != nil {
			return nil, err
		}

		body, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		var projects []model.Project
		err = json.Unmarshal(body, &projects)
		if err != nil {
			return nil, err
		}

		for i := range projects {
			projects[i].Workspace = &workspace
		}

		result = append(result, projects...)
	}

	return result, nil
}

func (a *API) Workspaces() ([]model.Workspace, error) {
	req, err := a.createRequest("GET", "workspaces", nil)
	if err != nil {
		return nil, err
	}

	res, err := a.client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var workspaces []model.Workspace
	err = json.Unmarshal(body, &workspaces)

	return workspaces, err
}

func (a *API) Project(workspaceId, projectId string) (model.Project, error) {
	req, err := a.createRequest("GET", path.Join("workspaces", workspaceId, "projects", projectId), nil)
	if err != nil {
		return model.Project{}, err
	}

	res, err := a.client.Do(req)
	if err != nil {
		return model.Project{}, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return model.Project{}, err
	}

	var project model.Project
	err = json.Unmarshal(body, &project)

	return project, err
}

func (a *API) CreateProject(workspaceId, name string) error {
	data, err := json.Marshal(model.CreateProject{
		Name: name,
	})
	if err != nil {
		return err
	}

	req, err := a.createRequest("POST", path.Join("workspaces", workspaceId, "projects"), bytes.NewReader(data))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	_, err = a.client.Do(req)
	if err != nil {
		return err
	}
	return nil
}
