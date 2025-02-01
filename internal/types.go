package internal

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

const (
	CreatePatchCommentType = "revision.comment"
	EventTypeHeaderKey     = "x-radicle-event-type"
	SignatureHashKey       = "x-radicle-signature"
	EventTypePush          = "push"
	EventTypePatch         = "patch"
)

type ListOpts struct {
	Page    int
	PerPage int
}

type NodeInfo struct {
	ID     string `json:"id"`
	Config Node   `json:"config"`
}

type SessionInfo struct {
	SessionId string `json:"sessionId"`
	Status    string `json:"status"`
	PublicKey string `json:"publicKey"`
	Alias     string `json:"alias"`
	IssuedAt  int64  `json:"issuedAt"`
	ExpiresAt int64  `json:"expiresAt"`
}

type Node struct {
	ID    string `json:"id"`
	Alias string `json:"alias"`
}

type Repository struct {
	ID             string         `json:"id"`
	Name           string         `json:"name"`
	Description    string         `json:"description"`
	Visibility     RepoVisibility `json:"visibility"`
	DefaultBranch  string         `json:"defaultBranch"`
	Default_Branch string         `json:"default_branch"`
	URL            string         `json:"url"`
	CloneURL       string         `json:"clone_url"`
	Delegates      []Delegates    `json:"delegates"`
	Head           string         `json:"head"`
}

type HookRepository struct {
	ID             string         `json:"id"`
	Name           string         `json:"name"`
	Description    string         `json:"description"`
	Visibility     RepoVisibility `json:"visibility"`
	DefaultBranch  string         `json:"defaultBranch"`
	Default_Branch string         `json:"default_branch"`
	URL            string         `json:"url"`
	CloneURL       string         `json:"clone_url"`
	Delegates      []string       `json:"delegates"`
	Head           string         `json:"head"`
}

type RepoVisibility struct {
	Type string `json:"type"`
}

type Delegates struct {
	ID    string `json:"id"`
	Alias string `json:"alias"`
}

type RepositoryCommits []*RepositoryCommit

type RepositoryCommit struct {
	ID      string   `json:"id"`
	Parents []string `json:"parents"`
}

type Commit struct {
	ID        string       `json:"id"`
	Title     string       `json:"title"`
	Message   string       `json:"message"`
	Timestamp time.Time    `json:"timestamp"`
	URL       string       `json:"url"`
	Author    CommitAuthor `json:"author"`
	Added     []string     `json:"added"`
	Modified  []string     `json:"modified"`
	Removed   []string     `json:"removed"`
}

type CommitAuthor struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type ProjectFile struct {
	Binary  bool   `json:"binary"`
	Name    string `json:"name"`
	Content string `json:"content"`
	Path    string `json:"path"`
}

type FileTree struct {
	Entries []FileTreeEntries `json:"entries"`
}

type FileTreeEntries struct {
	Path string `json:"path"`
	Name string `json:"name"`
	Kind string `json:"kind"`
}

type PatchState struct {
	Status string `json:"status"`
}

type CreatePatchComment struct {
	Type     string `json:"type"`
	Body     string `json:"body"`
	Revision string `json:"revision"`
}

type RepoWebhook struct {
	RepoID      string `json:"repo_id"`
	URL         string `json:"url"`
	Secret      string `json:"secret"`
	ContentType string `json:"content_type"`
}

type ΗοοκPushPayload struct {
	Author     Node           `json:"author"`
	Before     string         `json:"before"`
	After      string         `json:"after"`
	Commits    []Commit       `json:"commits"`
	Repository HookRepository `json:"repository"`
}

type ΗοοκPatchPayload struct {
	Action     string         `json:"action"`
	Patch      PatchHook      `json:"patch"`
	Repository HookRepository `json:"repository"`
}

type PatchHook struct {
	ID        string              `json:"id"`
	Author    Node                `json:"author"`
	Title     string              `json:"title"`
	State     State               `json:"state"`
	Before    string              `json:"before"`
	After     string              `json:"after"`
	Commits   []Commit            `json:"commits"`
	URL       string              `json:"url"`
	Target    string              `json:"target"`
	Labels    []string            `json:"labels"`
	Assignees []string            `json:"assignees"`
	Revisions []PatchRevisionHook `json:"revisions"`
}

type PatchRevisionHook struct {
	ID          string    `json:"id"`
	Author      Node      `json:"author"`
	Description string    `json:"description"`
	Base        string    `json:"base"`
	Oid         string    `json:"oid"`
	Timestamp   time.Time `json:"timestamp"`
}

type Patch struct {
	ID        string          `json:"id"`
	Author    Node            `json:"author"`
	Title     string          `json:"title"`
	State     State           `json:"state"`
	Before    string          `json:"before"`
	After     string          `json:"after"`
	Commits   []Commit        `json:"commits"`
	URL       string          `json:"url"`
	Target    string          `json:"target"`
	Labels    []string        `json:"labels"`
	Assignees []string        `json:"assignees"`
	Revisions []PatchRevision `json:"revisions"`
}

type PatchRevision struct {
	ID          string   `json:"id"`
	Author      Node     `json:"author"`
	Description string   `json:"description"`
	Base        string   `json:"base"`
	Oid         string   `json:"oid"`
	Timestamp   UnixTime `json:"timestamp"`
}

type Draft struct {
	Status string `json:"status"`
}

type State struct {
	Status    string      `json:"status"`
	Conflicts []Conflicts `json:"conflicts"`
}

type Archived struct {
	Status string `json:"status"`
}

type Merged struct {
	Status   string `json:"status"`
	Revision string `json:"revision"`
	Commit   string `json:"commit"`
}

type Conflicts struct {
	RevisionID string `json:"revision_id"`
	Oid        string `json:"oid"`
}

func (o *ListOpts) Encode() string {
	params := url.Values{}
	if o.Page > 0 {
		// Radicle's pagination starts from 0 but woodpecker's from 1
		page := o.Page - 1
		params.Set("page", strconv.Itoa(page))
	}
	if o.PerPage != 0 {
		params.Set("perPage", strconv.Itoa(o.PerPage))
	}
	return params.Encode()
}

type Error struct {
	Status int
	Body   struct {
		Message string `json:"message"`
	} `json:"error"`
}

func (e Error) Error() string {
	if len(e.Body.Message) > 0 {
		return e.Body.Message
	}
	return "an unknown error occurred with the forge"
}

type UnixTime struct {
	time.Time
}

func (u *UnixTime) UnmarshalJSON(b []byte) error {
	var timestamp int64
	err := json.Unmarshal(b, &timestamp)
	if err != nil {
		return err
	}
	u.Time = time.Unix(timestamp, 0)
	return nil
}

func (u UnixTime) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%d", (u.Time.Unix()))), nil
}
