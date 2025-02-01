package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"go.woodpecker-ci.org/woodpecker/v3/server/model"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	FileTypeBlob            = "blob"
	AppJsonType             = "application/json"
	SessionStatusAuthorized = "authorized"
)

const (
	apiPath                 = "/api"
	apiV1Path               = "/v1"
	pathNode                = "%s/node"
	pathSession             = "%s/sessions/%s"
	pathProject             = "%s/projects/%s"
	pathProjects            = "%s/projects?show=all&page=%s&perPage=%s"
	pathProjectCommits      = "%s/projects/%s/commits?%s"
	pathProjectCommitFile   = "%s/projects/%s/blob/%s/%s"
	pathProjectCommitDir    = "%s/projects/%s/tree/%s/%s"
	pathProjectPatches      = "%s/projects/%s/patches?%s"
	pathProjectPatchComment = "%s/projects/%s/patches/%s"
	pathProjectWebhooks     = "%s/projects/%s/webhooks"
	pathLogin               = "%s/oauth?callback_url=%s/authorize"
)

type Client struct {
	*http.Client
	base  string
	ctx   context.Context
	token string
}

func NewClient(ctx context.Context, url string, secretToken string) *Client {
	return &Client{
		Client: http.DefaultClient,
		base:   url,
		ctx:    ctx,
		token:  secretToken,
	}
}

func GetLoginURL(radURL, hostURL string) string {
	return fmt.Sprintf(pathLogin, radURL+apiPath+apiV1Path, hostURL)
}

func (c *Client) GetNodeInfo() (*NodeInfo, error) {
	out := new(NodeInfo)
	uri := fmt.Sprintf(pathNode, c.base+apiPath+apiV1Path)
	_, err := c.do(uri, http.MethodGet, nil, out)
	return out, err
}

func (c *Client) GetSessionInfo() (*SessionInfo, error) {
	out := new(SessionInfo)
	uri := fmt.Sprintf(pathSession, c.base+apiPath+apiV1Path, c.token)
	_, err := c.do(uri, http.MethodGet, nil, out)
	return out, err
}

func (c *Client) GetProject(projectID string) (*Repository, error) {
	out := new(Repository)
	uri := fmt.Sprintf(pathProject, c.base+apiPath+apiV1Path, projectID)
	_, err := c.do(uri, http.MethodGet, nil, out)
	return out, err
}

func (c *Client) GetProjects() ([]*Repository, error) {
	var projects []*Repository
	var err error
	page := 0
	perPage := 100
	for {
		out := new([]*Repository)
		uri := fmt.Sprintf(pathProjects, c.base+apiPath+apiV1Path, strconv.Itoa(page), strconv.Itoa(perPage))
		_, err = c.do(uri, http.MethodGet, nil, out)
		if err != nil {
			return nil, err
		}
		if len(*out) == 0 {
			break
		}
		page++
		projects = append(projects, *out...)
	}
	return projects, nil
}

func (c *Client) GetProjectCommits(projectID string, listOpts ListOpts) ([]*RepositoryCommit, error) {
	out := new(RepositoryCommits)
	uri := fmt.Sprintf(pathProjectCommits, c.base+apiPath+apiV1Path, projectID, listOpts.Encode())
	_, err := c.do(uri, http.MethodGet, nil, out)
	return *out, err
}

func (c *Client) GetProjectCommitFile(projectID, commit, file string) (*ProjectFile, error) {
	out := new(ProjectFile)
	uri := fmt.Sprintf(pathProjectCommitFile, c.base+apiPath+apiV1Path, projectID, commit, file)
	_, err := c.do(uri, http.MethodGet, nil, out)
	return out, err
}

func (c *Client) GetProjectCommitDir(projectID, commit, path string) (FileTree, error) {
	out := new(FileTree)
	uri := fmt.Sprintf(pathProjectCommitDir, c.base+apiPath+apiV1Path, projectID, commit, path)
	_, err := c.do(uri, http.MethodGet, nil, out)
	return *out, err
}

func (c *Client) GetProjectPatches(projectID string, listOpts ListOpts) ([]*Patch, error) {
	out := new([]*Patch)
	uri := fmt.Sprintf(pathProjectPatches, c.base+apiPath+apiV1Path, projectID, listOpts.Encode())
	_, err := c.do(uri, http.MethodGet, nil, out)
	return *out, err
}

func (c *Client) AddProjectPatchComment(projectID model.ForgeRemoteID, patchID string,
	commentPayload CreatePatchComment) error {
	uri := fmt.Sprintf(pathProjectPatchComment, c.base+apiPath+apiV1Path, projectID, patchID)
	_, err := c.do(uri, http.MethodPatch, commentPayload, nil)
	return err
}

func (c *Client) AddProjectWebhook(projectID model.ForgeRemoteID, webhookOpts RepoWebhook) error {
	uri := fmt.Sprintf(pathProjectWebhooks, c.base+apiPath+apiV1Path, projectID)
	_, err := c.do(uri, http.MethodPost, webhookOpts, nil)
	return err
}

func (c *Client) RemoveProjectWebhook(projectID model.ForgeRemoteID, url *string) error {
	uri := fmt.Sprintf(pathProjectWebhooks, c.base+apiPath+apiV1Path, projectID)
	if url != nil {
		uri = uri + "?url=" + *url
	}
	_, err := c.do(uri, http.MethodDelete, nil, nil)
	return err
}

func (c *Client) GetProjectWebhook(projectID model.ForgeRemoteID, link string) (*string, error) {
	uri := fmt.Sprintf(pathProjectWebhooks, c.base+apiPath+apiV1Path, projectID)
	var repoHooks []RepoWebhook
	_, err := c.do(uri, http.MethodGet, nil, &repoHooks)
	if err != nil {
		return nil, err
	}
	for _, repoHook := range repoHooks {
		if strings.HasPrefix(repoHook.URL, link) {
			return &repoHook.URL, nil
		}
	}
	return nil, nil
}

func (c *Client) do(rawurl, method string, in, out interface{}) (*string, error) {
	uri, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}

	// if we are posting or putting data, we need to
	// write it to the body of the request.
	var buf io.ReadWriter
	if in != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(in)
		if err != nil {
			return nil, err
		}
	}
	// creates a new http request to radicle httpd.
	req, err := http.NewRequestWithContext(c.ctx, method, uri.String(), buf)
	if err != nil {
		return nil, err
	}
	if in != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if len(c.token) > 0 {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// if an error is encountered, parse and return the
	// error response.
	if resp.StatusCode >= http.StatusBadRequest {
		err := Error{}
		_ = json.NewDecoder(resp.Body).Decode(&err)
		err.Status = resp.StatusCode
		return nil, err
	}

	// if a json response is expected, parse and return
	// the json response.
	if out != nil {
		return nil, json.NewDecoder(resp.Body).Decode(out)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	bodyString := string(bodyBytes)
	return &bodyString, nil
}
