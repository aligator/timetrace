package clockify

import (
	"encoding/json"
	"github.com/dominikbraun/timetrace/clockify/clockify/model"
	"github.com/dominikbraun/timetrace/core"
	"path"
	"strings"
	"time"
)

type projectData struct {
	model.Project
}

type workspaceData struct {
	model.Workspace
	projects map[string]projectData
}

type fs struct {
	API *API

	workspaces map[string]workspaceData
}

func NewFs(config *Config) core.Filesystem {
	fs := &fs{
		API: &API{
			Config: config,
			url:    nil,
		},
	}

	fs.workspaces = make(map[string]workspaceData)
	return fs
}

func removeDisplayPart(key string) string {
	// Remove everything after " " as it is only the display name (if it even exists).
	return strings.SplitN(key, " ", 2)[0]
}

func (f *fs) ProjectFilepath(key string) string {
	return removeDisplayPart(key)
}

func (f *fs) ProjectFilepaths() ([]string, error) {
	projects, err := f.API.Projects()
	if err != nil {
		return nil, err
	}

	var result []string

	for _, project := range projects {
		result = append(result, path.Join(project.WorkspaceID, project.ID))

		w, ok := f.workspaces[project.WorkspaceID]
		if !ok {
			w = workspaceData{
				Workspace: *project.Workspace,
			}
			w.projects = make(map[string]projectData)
		}

		w.projects[project.ID] = projectData{
			project,
		}

		f.workspaces[project.WorkspaceID] = w
	}

	return result, nil
}

func (f *fs) RecordFilepath(start time.Time) string {
	panic("implement me")
}

func (f *fs) RecordFilepaths(dir string, less func(a string, b string) bool) ([]string, error) {
	panic("implement me")
}

func (f *fs) RecordDirs() ([]string, error) {
	panic("implement me")
}

func (f *fs) RecordDirFromDate(date time.Time) string {
	panic("implement me")
}

func (f *fs) EnsureDirectories() error {
	return nil
}

func (f *fs) EnsureRecordDir(date time.Time) error {
	return nil
}

func (f *fs) ReadFile(filename string) ([]byte, error) {
	pathParts := strings.Split(removeDisplayPart(filename), "/")

	// First check if the workspace is already loaded:
	workspace, ok := f.workspaces[pathParts[0]]
	if !ok {
		// If not load the workspacees.
		workspaces, err := f.API.Workspaces()
		if err != nil {
			return nil, err
		}

		for _, w := range workspaces {
			f.workspaces[w.ID] = workspaceData{
				Workspace: w,
				projects:  make(map[string]projectData),
			}

			if w.ID == pathParts[0] {
				workspace = f.workspaces[w.ID]
			}
		}
	}

	// Now do the same for the project
	project, ok := f.workspaces[pathParts[0]].projects[pathParts[1]]
	if !ok {
		// If not load the project.
		p, err := f.API.Project(workspace.ID, pathParts[1])
		if err != nil {
			return nil, err
		}

		f.workspaces[workspace.ID].projects[p.ID] = projectData{
			p,
		}

		project = f.workspaces[workspace.ID].projects[p.ID]
	}

	key := path.Join(workspace.ID, project.ID) + " " + path.Join(workspace.Name, project.Name)

	return json.Marshal(core.Project{
		Key: key,
	})
}

func (f *fs) Exists(path string) (bool, error) {
	pathParts := strings.Split(removeDisplayPart(path), "/")
	p, err := f.API.Project(pathParts[0], pathParts[1])
	return p != (model.Project{}), err
}

func (f *fs) Write(path string, bytes []byte) (int, error) {
	pathParts := strings.Split(removeDisplayPart(path), "/")

	var project core.Project
	err := json.Unmarshal(bytes, &project)
	if err != nil {
		return 0, err
	}

	err = f.API.CreateProject(pathParts[0], strings.TrimPrefix(project.Key, pathParts[0]+"/"))
	if err != nil {
		return 0, err
	}

	return len(bytes), nil
}
