package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"go.woodpecker-ci.org/woodpecker/v3/server/model"
	"radicle-woodpecker-addon/internal"
)

// parsePushHook parses the push hook payload
func (rad *radicle) parsePushHook(payload []byte) (*model.Repo, *model.Pipeline, error) {
	hook := internal.ΗοοκPushPayload{}
	if err := json.Unmarshal(payload, &hook); err != nil {
		return nil, nil, err
	}
	repo := convertHookProject(&hook.Repository, nil, rad)
	lastCommit := internal.Commit{}
	if len(hook.Commits) == 0 {
		return nil, nil, errors.New("no commits found in push")
	}
	lastCommit = hook.Commits[len(hook.Commits)-1]
	changedFiles := []string{}
	for _, commit := range hook.Commits {
		changedFiles = append(changedFiles, commit.Modified...)
		changedFiles = append(changedFiles, commit.Added...)
		changedFiles = append(changedFiles, commit.Removed...)
	}
	pipeline := model.Pipeline{
		Author:       hook.Author.Alias,
		Event:        model.EventPush,
		Commit:       hook.After,
		Branch:       hook.After,
		Ref:          hook.After,
		Avatar:       RADICLE_IMAGE,
		Message:      lastCommit.Title,
		Timestamp:    lastCommit.Timestamp.Unix(),
		Sender:       lastCommit.Author.Name,
		Email:        lastCommit.Author.Email,
		ForgeURL:     lastCommit.URL,
		ChangedFiles: changedFiles,
	}

	return repo, &pipeline, nil
}

// parsePatchHook parses the patch hook payload
func (rad *radicle) parsePatchHook(payload []byte) (*model.Repo, *model.Pipeline, error) {
	hook := internal.ΗοοκPatchPayload{}
	if err := json.Unmarshal(payload, &hook); err != nil {
		return nil, nil, err
	}
	repo := convertHookProject(&hook.Repository, nil, rad)
	if len(hook.Patch.Revisions) == 0 {
		return nil, nil, errors.New("no revision found in patch")
	}
	lastRevision := hook.Patch.Revisions[len(hook.Patch.Revisions)-1]
	changedFiles := []string{}
	for _, commit := range hook.Patch.Commits {
		changedFiles = append(changedFiles, commit.Modified...)
		changedFiles = append(changedFiles, commit.Added...)
		changedFiles = append(changedFiles, commit.Removed...)
	}
	defaultBranch := hook.Repository.DefaultBranch
	if len(defaultBranch) == 0 {
		defaultBranch = hook.Repository.Default_Branch
	}
	vars := map[string]string{}
	vars["patch_id"] = hook.Patch.ID
	vars["revision_id"] = lastRevision.ID
	pipeline := model.Pipeline{
		Author:              hook.Patch.Author.Alias,
		Event:               model.EventPull,
		Commit:              hook.Patch.After,
		Branch:              hook.Patch.ID,
		Ref:                 hook.Patch.After,
		Refspec:             fmt.Sprintf("refs/patches/%s:refs/heads/%s", hook.Patch.After, defaultBranch),
		Avatar:              RADICLE_IMAGE,
		PullRequestLabels:   hook.Patch.Labels,
		Message:             hook.Patch.Title,
		Timestamp:           lastRevision.Timestamp.Unix(),
		Sender:              hook.Patch.Author.ID,
		Email:               hook.Patch.Author.Alias,
		ForgeURL:            hook.Patch.URL,
		ChangedFiles:        changedFiles,
		AdditionalVariables: vars,
	}

	return repo, &pipeline, nil
}
