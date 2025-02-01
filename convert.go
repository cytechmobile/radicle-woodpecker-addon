package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	forge_types "go.woodpecker-ci.org/woodpecker/v3/server/forge/types"
	"go.woodpecker-ci.org/woodpecker/v3/server/model"
	"radicle-woodpecker-addon/internal"
	"strings"
)

// As radicle does not support user avatars, use radicle's logo as avatar
const RADICLE_IMAGE = "data:image/png;base64," +
	"iVBORw0KGgoAAAANSUhEUgAAACwAAAAsBAMAAADsqkcyAAAAElBMVEUAAAAzM91VVf/09PT/Vf////+iehdrAAAAAXRSTlMAQObYZgAAAAFiS0dEBfhv6ccAAABTSURBVCjPY2AAAiUlBjhAYlNFGMFBEaSasKAgBFNbWAkFUFMYxDQGAgRNbWEXFwRNLWFBrIA6wgxg8yGRC8JQoSEhDEuUaGmKSsKQZIOSpigXBgAOHTr5ND3M6gAAAABJRU5ErkJggg=="

// convertUser is a helper function used to convert a Radicle Node Info structure
// to the Woodpecker User structure.
func convertUser(rad *radicle) *model.User {
	return &model.User{
		ForgeRemoteID: model.ForgeRemoteID(rad.nodeID),
		Login:         rad.alias,
		Avatar:        RADICLE_IMAGE,
		AccessToken:   rad.sessionToken,
	}
}

// convertProject is a helper function used to convert a Radicle Project structure
// to the Woodpecker Repo structure.
func convertProject(project *internal.Repository, user *model.User, rad *radicle) *model.Repo {
	projectID := strings.TrimPrefix(project.ID, "rad:")
	defaultBranch := project.DefaultBranch
	if len(defaultBranch) == 0 {
		defaultBranch = project.Default_Branch
	}
	return &model.Repo{
		ForgeRemoteID: model.ForgeRemoteID(projectID),
		Name:          fmt.Sprintf("%s (%s)", project.Name, project.ID),
		FullName:      fmt.Sprintf("%s (%s)", project.Name, project.ID),
		ForgeURL:      fmt.Sprintf("%s/%s", rad.URL(), projectID),
		Clone:         fmt.Sprintf("%s/%s.git", rad.URL(), projectID),
		Hash:          project.ID,
		Avatar:        RADICLE_IMAGE,
		CloneSSH:      "",
		Branch:        defaultBranch,
		Perm: &model.Perm{
			Pull:  true,
			Push:  true,
			Admin: true,
		},
		Owner:     rad.Name(),
		PREnabled: true,
	}
}

// convertHookProject is a helper function used to convert a Radicle Project structure
// to the Woodpecker Repo structure from the hook info.
func convertHookProject(project *internal.HookRepository, user *model.User, rad *radicle) *model.Repo {
	projectID := strings.TrimPrefix(project.ID, "rad:")
	defaultBranch := project.DefaultBranch
	if len(defaultBranch) == 0 {
		defaultBranch = project.Default_Branch
	}
	return &model.Repo{
		ForgeRemoteID: model.ForgeRemoteID(projectID),
		Name:          fmt.Sprintf("%s (%s)", project.Name, project.ID),
		FullName:      fmt.Sprintf("%s (%s)", project.Name, project.ID),
		ForgeURL:      fmt.Sprintf("%s/%s", rad.URL(), projectID),
		Clone:         fmt.Sprintf("%s/%s.git", rad.URL(), projectID),
		Hash:          project.ID,
		Avatar:        RADICLE_IMAGE,
		CloneSSH:      "",
		Branch:        defaultBranch,
		Perm: &model.Perm{
			Pull:  true,
			Push:  true,
			Admin: true,
		},
		Owner:     rad.Name(),
		PREnabled: true,
	}
}

// convertProjectFileToContent is a helper function used to convert a Radicle Project File content
// to the Woodpecker file content structure.
func convertProjectFileToContent(projectFile *internal.ProjectFile) ([]byte, error) {
	return []byte(projectFile.Content), nil
}

// convertFileContent is a helper function used to convert a Radicle Project file Contents
// to the Woodpecker File Meta structure.
func convertFileContent(fileContentEntries internal.FileTreeEntries, fileContent []byte) *forge_types.FileMeta {
	return &forge_types.FileMeta{
		Name: fileContentEntries.Path,
		Data: fileContent,
	}
}

// convertProjectPatch is a helper function used to convert a Radicle patch
// to the Woodpecker Pull Request structure.
func convertProjectPatch(patch *internal.Patch) *model.PullRequest {
	return &model.PullRequest{
		Index: model.ForgeRemoteID(patch.ID),
		Title: patch.Title,
	}
}

// generateHmacSignature generates an hmac signature with the given key and message
func generateHmacSignature(key string, msg []byte) string {
	hmacSha256 := hmac.New(sha256.New, []byte(key))
	hmacSha256.Write(msg)
	res := hex.EncodeToString(hmacSha256.Sum(nil))
	return res
}
