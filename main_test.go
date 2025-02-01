package main

import (
	"bytes"
	"github.com/franela/goblin"
	"github.com/gin-gonic/gin"
	forgeTypes "go.woodpecker-ci.org/woodpecker/v3/server/forge/types"
	"go.woodpecker-ci.org/woodpecker/v3/server/model"
	"golang.org/x/net/context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"radicle-woodpecker-addon/fixtures"
	"radicle-woodpecker-addon/internal"
	"testing"
)

var forgeOpts = radicleOpts{
	URL:               "http://node.id",
	WoodpeckerHostURL: "http://woodpecker.host",
}

func Test_radicle(t *testing.T) {
	gin.SetMode(gin.TestMode)
	s := httptest.NewServer(fixtures.Handler())

	g := goblin.Goblin(t)
	g.Describe("Radicle client", func() {
		g.After(func() {
			s.Close()
		})

		forgeOpts.URL = s.URL
		forge, _ := newRadicle(forgeOpts)

		// Test New()
		g.Describe("Creating new Forge", func() {
			g.It("should return an error when missing URL", func() {
				opts := radicleOpts{
					URL:               "",
					WoodpeckerHostURL: "http://woodpecker.host",
				}
				_, err := newRadicle(opts)
				g.Assert(err).IsNotNil()
			})
			g.It("Should return an error when invalid URL", func() {
				opts := radicleOpts{
					URL:               "invalid_%url",
					WoodpeckerHostURL: "http://woodpecker.host",
				}
				_, err := newRadicle(opts)
				g.Assert(err).IsNotNil()
			})
			g.It("Should not return an error when missing Node ID", func() {
				opts := radicleOpts{
					URL:               "http://some.url",
					WoodpeckerHostURL: "http://woodpecker.host",
				}
				_, err := newRadicle(opts)
				g.Assert(err).IsNil("Not expected error")
			})
			g.It("Should return a new Forge with correct data", func() {
				opts := radicleOpts{
					URL:               "http://some.url",
					WoodpeckerHostURL: "http://woodpecker.host",
				}
				forge, err := newRadicle(opts)
				g.Assert(err).IsNil()
				g.Assert(forge.URL()).Equal("http://some.url")
				g.Assert(forge.Name()).Equal("radicle")
			})
			g.It("Should return a new Forge with correct data and trim slashes in URLs", func() {
				opts := radicleOpts{
					URL:               "http://some.url/",
					WoodpeckerHostURL: "http://woodpecker.host",
				}
				forge, err := newRadicle(opts)
				g.Assert(err).IsNil()
				g.Assert(forge.URL()).Equal("http://some.url")
				g.Assert(forge.Name()).Equal("radicle")
			})
		})

		// Test Login()
		g.Describe("When logging in", func() {
			g.Describe("without session ID", func() {
				forgeOpts.URL = s.URL
				forge, err := newRadicle(forgeOpts)
				g.It("Should fail", func() {
					g.Assert(err).IsNil()
					req := &forgeTypes.OAuthRequest{
						Code:  "",
						State: "",
					}
					user, loginURL, err := forge.Login(context.Background(), req)
					g.Assert(err).IsNil()
					g.Assert(user).IsNil()
					g.Assert(loginURL[22:]).Equal("/api/v1/oauth?callback_url=http://woodpecker.host/authorize")
				})
			})
			g.Describe("with unauthorized session ID", func() {
				forgeOpts.URL = s.URL
				forge, err := newRadicle(forgeOpts)
				g.It("Should fail", func() {
					g.Assert(err).IsNil()
					c := gin.Context{}
					c.Request = httptest.NewRequest("POST", "/", nil)
					c.Request.Form = url.Values{}
					c.Request.Form.Add("session_id", "unauthed_sess_id")
					req := &forgeTypes.OAuthRequest{
						Code:  "unauthed_sess_id",
						State: "",
					}
					user, loginURL, err := forge.Login(&c, req)
					g.Assert(err).IsNotNil()
					g.Assert(loginURL[22:]).Equal("/api/v1/oauth?callback_url=http://woodpecker.host/authorize")
					g.Assert(err.Error()).Equal("provided secret token is unauthorized")
					g.Assert(user).IsNil()
				})
			})
			g.Describe("with authorized session ID", func() {
				forgeOpts.URL = s.URL
				forge, err := newRadicle(forgeOpts)
				g.It("Should succeed", func() {
					g.Assert(err).IsNil()
					c := gin.Context{}
					c.Request = httptest.NewRequest("POST", "/", nil)
					c.Request.Form = url.Values{}
					c.Request.Form.Add("session_id", "authed_sess_id")
					req := &forgeTypes.OAuthRequest{
						Code:  "authed_sess_id",
						State: "",
					}
					user, loginURL, err := forge.Login(&c, req)
					g.Assert(err).IsNil()
					g.Assert(user).IsNotNil()
					g.Assert(loginURL).IsNotNil()
					g.Assert(user.Login).Equal("myalias")
					g.Assert(user.AccessToken).Equal("authed_sess_id")
				})
			})
		})

		// Test Repo()
		g.Describe("When requesting a project", func() {
			user := &model.User{
				ForgeRemoteID: "remote_user_id",
				Login:         "user_login",
			}
			g.Describe("with invalid project ID", func() {
				g.It("Should fail", func() {
					repo, err := forge.Repo(context.Background(), user, "not_found", "", "")
					g.Assert(err).IsNotNil()
					g.Assert(repo).IsNil()
				})
			})
			g.Describe("with authorized sessions ID", func() {
				g.It("Should succeed", func() {
					repo, err := forge.Repo(context.Background(), user, "valid_project_id", "", "")
					g.Assert(err).IsNil()
					g.Assert(repo).IsNotNil()
					g.Assert(repo.ForgeRemoteID).Equal(model.ForgeRemoteID("valid_project_id"))
					g.Assert(repo.Name).Equal("a-project (rad:valid_project_id)")
					g.Assert(repo.FullName).Equal("a-project (rad:valid_project_id)")
					g.Assert(repo.ForgeURL).Equal(forge.URL() + "/valid_project_id")
					g.Assert(repo.Clone).Equal(forge.URL() + "/valid_project_id.git")
					g.Assert(repo.CloneSSH).Equal("")
					g.Assert(repo.Branch).Equal("main")
					g.Assert(repo.Perm.Pull).Equal(true)
					g.Assert(repo.Perm.Pull).Equal(true)
					g.Assert(repo.Perm.Admin).Equal(true)
					g.Assert(repo.Owner).Equal("radicle")
				})
			})
		})

		// Test Repos()
		g.Describe("When requesting projects", func() {
			user := &model.User{
				ForgeRemoteID: "remote_user_id",
				Login:         "user_login",
			}
			g.Describe("with many projects", func() {
				g.It("Should succeed", func() {
					repos, err := forge.Repos(context.Background(), user)
					g.Assert(err).IsNil()
					g.Assert(len(repos)).Equal(4)

					g.Assert(repos[0].ForgeRemoteID).Equal(model.ForgeRemoteID("a_project"))
					g.Assert(repos[0].Name).Equal("a-project (rad:a_project)")
					g.Assert(repos[0].FullName).Equal("a-project (rad:a_project)")
					g.Assert(repos[0].ForgeURL).Equal(forge.URL() + "/a_project")
					g.Assert(repos[0].Clone).Equal(forge.URL() + "/a_project.git")
					g.Assert(repos[0].CloneSSH).Equal("")
					g.Assert(repos[0].Branch).Equal("main")
					g.Assert(repos[0].Perm.Pull).Equal(true)
					g.Assert(repos[0].Perm.Pull).Equal(true)
					g.Assert(repos[0].Perm.Admin).Equal(true)
					g.Assert(repos[0].Owner).Equal("radicle")

					g.Assert(repos[1].ForgeRemoteID).Equal(model.ForgeRemoteID("b_project"))
					g.Assert(repos[1].Name).Equal("b-project (rad:b_project)")
					g.Assert(repos[1].FullName).Equal("b-project (rad:b_project)")
					g.Assert(repos[1].ForgeURL).Equal(forge.URL() + "/b_project")
					g.Assert(repos[1].Clone).Equal(forge.URL() + "/b_project.git")
					g.Assert(repos[1].CloneSSH).Equal("")
					g.Assert(repos[1].Branch).Equal("master")
					g.Assert(repos[1].Perm.Pull).Equal(true)
					g.Assert(repos[1].Perm.Pull).Equal(true)
					g.Assert(repos[1].Perm.Admin).Equal(true)
					g.Assert(repos[1].Owner).Equal("radicle")

					g.Assert(repos[2].ForgeRemoteID).Equal(model.ForgeRemoteID("c_project"))
					g.Assert(repos[2].Name).Equal("c-project (rad:c_project)")
					g.Assert(repos[2].FullName).Equal("c-project (rad:c_project)")
					g.Assert(repos[2].ForgeURL).Equal(forge.URL() + "/c_project")
					g.Assert(repos[2].Clone).Equal(forge.URL() + "/c_project.git")
					g.Assert(repos[2].CloneSSH).Equal("")
					g.Assert(repos[2].Branch).Equal("main")
					g.Assert(repos[2].Perm.Pull).Equal(true)
					g.Assert(repos[2].Perm.Pull).Equal(true)
					g.Assert(repos[2].Perm.Admin).Equal(true)
					g.Assert(repos[2].Owner).Equal("radicle")

					g.Assert(repos[3].ForgeRemoteID).Equal(model.ForgeRemoteID("d_project"))
					g.Assert(repos[3].Name).Equal("d-project (rad:d_project)")
					g.Assert(repos[3].FullName).Equal("d-project (rad:d_project)")
					g.Assert(repos[3].ForgeURL).Equal(forge.URL() + "/d_project")
					g.Assert(repos[3].Clone).Equal(forge.URL() + "/d_project.git")
					g.Assert(repos[3].CloneSSH).Equal("")
					g.Assert(repos[3].Branch).Equal("master")
					g.Assert(repos[3].Perm.Pull).Equal(true)
					g.Assert(repos[3].Perm.Pull).Equal(true)
					g.Assert(repos[3].Perm.Admin).Equal(true)
					g.Assert(repos[3].Owner).Equal("radicle")
				})
			})
		})

		// Test File()
		g.Describe("When requesting a project file", func() {
			user := &model.User{
				ForgeRemoteID: "remote_user_id",
				Login:         "user_login",
			}
			pipeline := &model.Pipeline{
				Commit: "the_commit_id",
			}
			g.Describe("with invalid project ID", func() {
				g.It("Should fail", func() {
					repo := &model.Repo{
						ForgeRemoteID: "not_found",
					}
					file, err := forge.File(context.Background(), user, repo, pipeline, "file_name.md")
					g.Assert(err).IsNotNil()
					g.Assert(file).IsNil()
				})
			})
			g.Describe("with valid project ID", func() {
				g.It("Should succeed", func() {
					repo := &model.Repo{
						ForgeRemoteID: "remote_project_id",
					}
					file, err := forge.File(context.Background(), user, repo, pipeline, "file_name.md")
					g.Assert(err).IsNil()
					g.Assert(file).IsNotNil()
					g.Assert(file).Equal([]byte("file content"))
				})
			})
		})

		// Test Tree()
		g.Describe("When requesting a project dir", func() {
			user := &model.User{
				ForgeRemoteID: "remote_user_id",
				Login:         "user_login",
			}
			pipeline := &model.Pipeline{
				Commit: "the_commit_id",
			}
			g.Describe("with invalid project ID", func() {
				g.It("Should fail", func() {
					repo := &model.Repo{
						ForgeRemoteID: "not_found",
					}
					tree, err := forge.Dir(context.Background(), user, repo, pipeline, "")
					g.Assert(err).IsNotNil()
					g.Assert(tree).IsNil()
				})
			})
			g.Describe("with valid project ID", func() {
				g.It("Should succeed", func() {
					repo := &model.Repo{
						ForgeRemoteID: "remote_project_id",
					}
					tree, err := forge.Dir(context.Background(), user, repo, pipeline, "")
					g.Assert(err).IsNil()
					g.Assert(tree).IsNotNil()
					g.Assert(len(tree)).Equal(3)
					g.Assert(tree[0].Name).Equal("Readme.md")
					g.Assert(tree[0].Data).Equal([]byte("file content"))
					g.Assert(tree[1].Name).Equal("cargo-checksum.json")
					g.Assert(tree[1].Data).Equal([]byte("file content"))
					g.Assert(tree[2].Name).Equal("debian")
					g.Assert(tree[2].Data).Equal([]byte(""))
				})
			})
		})

		// Test Status()
		g.Describe("When updating a pipeline status", func() {
			g.Describe("with pending status", func() {
				g.It("Should ignore it", func() {
					repo := &model.Repo{
						ForgeRemoteID: "invalid",
					}
					vars := map[string]string{}
					vars["patch_id"] = "patchID"
					vars["revision_id"] = "revID"
					pipeline := &model.Pipeline{
						Commit:              "the_commit_id",
						Status:              model.StatusPending,
						AdditionalVariables: vars,
					}
					err := forge.Status(context.Background(), nil, repo, pipeline, nil)
					g.Assert(err).IsNil()
				})
			})
			g.Describe("with running status", func() {
				g.It("Should ignore it", func() {
					repo := &model.Repo{
						ForgeRemoteID: "invalid",
					}
					vars := map[string]string{}
					vars["patch_id"] = "patchID"
					vars["revision_id"] = "revID"
					pipeline := &model.Pipeline{
						Commit:              "the_commit_id",
						Status:              model.StatusRunning,
						AdditionalVariables: vars,
					}
					err := forge.Status(context.Background(), nil, repo, pipeline, nil)
					g.Assert(err).IsNil()
				})
			})
			g.Describe("with running status without the patch ID", func() {
				g.It("Should not fail", func() {
					repo := &model.Repo{
						ForgeRemoteID: "invalid",
					}
					vars := map[string]string{}
					vars["revision_id"] = "revID"
					pipeline := &model.Pipeline{
						Commit:              "the_commit_id",
						Status:              model.StatusSuccess,
						AdditionalVariables: vars,
					}
					err := forge.Status(context.Background(), nil, repo, pipeline, nil)
					g.Assert(err).IsNil()
				})
			})
			g.Describe("with failure status on patch", func() {
				g.It("should add patch comment", func() {
					repo := &model.Repo{
						ForgeRemoteID: "repo_id",
					}
					vars := map[string]string{}
					vars["patch_id"] = "patchID"
					vars["revision_id"] = "revID"
					pipeline := &model.Pipeline{
						Commit:              "the_commit_id",
						Status:              model.StatusFailure,
						AdditionalVariables: vars,
					}
					u := model.User{
						AccessToken: "some_AccessToken",
					}
					err := forge.Status(context.Background(), &u, repo, pipeline, nil)
					g.Assert(err).IsNil()
				})
			})

			g.Describe("with failure status on patch with invalid patch ID", func() {
				g.It("should fail", func() {
					repo := &model.Repo{
						ForgeRemoteID: "repo_id",
					}
					vars := map[string]string{}
					vars["patch_id"] = "invalid_patchID"
					vars["revision_id"] = "revID"
					pipeline := &model.Pipeline{
						Commit:              "the_commit_id",
						Status:              model.StatusFailure,
						AdditionalVariables: vars,
					}
					u := model.User{
						AccessToken: "some_token",
					}
					err := forge.Status(context.Background(), &u, repo, pipeline, nil)
					g.Assert(err).IsNotNil()
				})
			})
		})

		//Test Activate
		g.Describe("When activating a repo", func() {
			g.Describe("for a non existing repo ID", func() {
				g.It("Should fail", func() {
					repo := &model.Repo{
						ForgeRemoteID: "not_found",
					}
					u := model.User{
						AccessToken: "some_token",
					}
					err := forge.Activate(context.Background(), &u, repo, "")
					g.Assert(err).IsNotNil()
				})
			})
			g.Describe("for an invalid repo ID", func() {
				g.It("Should fail", func() {
					repo := &model.Repo{
						ForgeRemoteID: "invalid",
					}
					u := model.User{
						AccessToken: "some_token",
					}
					err := forge.Activate(context.Background(), &u, repo, "")
					g.Assert(err).IsNotNil()
				})
			})
			g.Describe("for a valid repo ID", func() {
				g.It("Should succeed", func() {
					repo := &model.Repo{
						ForgeRemoteID: "repo_id",
					}
					u := model.User{
						AccessToken: "some_token",
					}
					err := forge.Activate(context.Background(), &u, repo, "")
					g.Assert(err).IsNil()
				})
			})
		})

		//Test Deactivate
		g.Describe("When deactivating a repo", func() {
			g.Describe("for an invalid repo ID", func() {
				g.It("Should fail", func() {
					repo := &model.Repo{
						ForgeRemoteID: "invalid",
					}
					u := model.User{
						AccessToken: "some_token",
					}
					err := forge.Deactivate(context.Background(), &u, repo, "")
					g.Assert(err).IsNotNil()
				})
			})
			g.Describe("for any repo ID", func() {
				g.It("Should succeed", func() {
					repo := &model.Repo{
						ForgeRemoteID: "some_repo_id",
					}
					u := model.User{
						AccessToken: "some_token",
					}
					err := forge.Deactivate(context.Background(), &u, repo, "")
					g.Assert(err).IsNil()
				})
			})
		})

		//Test Branches
		g.Describe("When a repo's branches", func() {
			g.Describe("for the first page", func() {
				g.It("Should return only the default branch", func() {
					repo := &model.Repo{
						ForgeRemoteID: "some_repo_id",
						Branch:        "the_default_branch",
					}
					listOpts := model.ListOptions{
						Page:    1,
						PerPage: 10,
					}
					u := model.User{
						AccessToken: "some_token",
					}
					branches, err := forge.Branches(context.Background(), &u, repo, &listOpts)
					g.Assert(err).IsNil()
					g.Assert(len(branches)).Equal(1)
					g.Assert(branches[0]).Equal("the_default_branch")
				})
			})
			g.Describe("for the second page", func() {
				g.It("Should return nothing", func() {
					repo := &model.Repo{
						ForgeRemoteID: "some_repo_id",
						Branch:        "the_default_branch",
					}
					listOpts := model.ListOptions{
						Page:    2,
						PerPage: 10,
					}
					u := model.User{
						AccessToken: "some_token",
					}
					branches, err := forge.Branches(context.Background(), &u, repo, &listOpts)
					g.Assert(err).IsNil()
					g.Assert(len(branches)).Equal(0)
				})
			})
		})

		// Test BranchHead
		g.Describe("When requesting a project's BranchHead", func() {
			user := &model.User{
				ForgeRemoteID: "remote_user_id",
				Login:         "user_login",
				AccessToken:   "some_token",
			}

			g.Describe("with invalid project ID", func() {
				g.It("Should fail", func() {
					repo := &model.Repo{
						ForgeRemoteID: "not_found",
						Branch:        "main",
					}
					commit, err := forge.BranchHead(context.Background(), user, repo, "main")
					g.Assert(err).IsNotNil()
					g.Assert(commit).IsNil()
				})
			})
			g.Describe("with valid project ID", func() {
				g.It("Should succeed", func() {
					repo := &model.Repo{
						ForgeRemoteID: "valid_project_id",
						Branch:        "main",
					}
					branchHead, err := forge.BranchHead(context.Background(), user, repo, "main")
					g.Assert(err).IsNil()
					g.Assert(branchHead).IsNotNil()
					g.Assert(branchHead.SHA).Equal("00bfa9b18be32001481334126c311c4a327dff2e")

				})
			})
			g.Describe("with valid project ID and invalid branch", func() {
				g.It("Should fail", func() {
					repo := &model.Repo{
						ForgeRemoteID: "valid_project_id",
						Branch:        "invalid_branch",
					}
					commit, err := forge.BranchHead(context.Background(), user, repo, "main")
					g.Assert(err).IsNotNil()
					g.Assert(err.Error()).Equal("branch does not exist")
					g.Assert(commit).IsNil()
				})
			})
		})

		// Test PullRequests
		g.Describe("When requesting a project's PullRequests", func() {
			g.Describe("with invalid project ID", func() {
				g.It("Should fail", func() {
					repo := &model.Repo{
						ForgeRemoteID: "not_found",
						Branch:        "main",
					}
					listOpts := model.ListOptions{
						Page:    1,
						PerPage: 10,
					}
					u := model.User{
						AccessToken: "some_token",
					}
					pullRequests, err := forge.PullRequests(context.Background(), &u, repo, &listOpts)
					g.Assert(err).IsNotNil()
					g.Assert(len(pullRequests)).Equal(0)
				})
			})

			g.Describe("with valid project ID", func() {
				g.It("Should succeed", func() {
					repo := &model.Repo{
						ForgeRemoteID: "repo_id",
						Branch:        "main",
					}
					listOpts := model.ListOptions{
						Page:    1,
						PerPage: 10,
					}
					u := model.User{
						AccessToken: "some_token",
					}
					pullRequests, err := forge.PullRequests(context.Background(), &u, repo, &listOpts)
					g.Assert(err).IsNil()
					g.Assert(len(pullRequests)).Equal(2)

					g.Assert(pullRequests[0].Index).Equal(model.ForgeRemoteID(
						"c7eee5122d0467aec5e71c228c958f9c79fe17c9"))
					g.Assert(pullRequests[0].Title).Equal("Use repository consistently over project")

					g.Assert(pullRequests[1].Index).Equal(model.ForgeRemoteID(
						"c969b0642425b6fdd34ade8df866cc0d330e8cc5"))
					g.Assert(pullRequests[1].Title).Equal("httpd: Add project sync status endpoint")
				})
			})
		})

		// Test Hook
		g.Describe("When requesting to parse a hook", func() {
			g.Describe("with an unknown event type", func() {
				g.It("Should fail", func() {
					buf := bytes.NewBufferString("")
					req, _ := http.NewRequest("POST", "/hook", buf)
					req.Header = http.Header{}
					req.Header.Set(internal.EventTypeHeaderKey, "an_event")
					_, _, err := forge.Hook(context.Background(), req)
					g.Assert(err).IsNotNil()
				})
			})

			g.Describe("with a valid push event type", func() {
				g.It("Should succeed", func() {
					buf := bytes.NewBufferString(fixtures.HookPushPayload)
					req, _ := http.NewRequest("POST", "/hook", buf)
					req.Header = http.Header{}
					req.Header.Set(internal.EventTypeHeaderKey, "push")
					repo, pipeline, err := forge.Hook(context.Background(), req)
					g.Assert(err).IsNil()
					g.Assert(repo.ForgeRemoteID).Equal(model.ForgeRemoteID("z3gqcJUoA1n9HaHKufZs5FCSGazv5"))
					g.Assert(repo.Name).Equal("heartwood (rad:z3gqcJUoA1n9HaHKufZs5FCSGazv5)")
					g.Assert(repo.FullName).Equal("heartwood (rad:z3gqcJUoA1n9HaHKufZs5FCSGazv5)")
					g.Assert(repo.Branch).Equal("master")
					g.Assert(repo.PREnabled).IsTrue()
					g.Assert(repo.Clone).Equal(forge.URL() + "/z3gqcJUoA1n9HaHKufZs5FCSGazv5.git")
					g.Assert(repo.ForgeURL).Equal(forge.URL() + "/z3gqcJUoA1n9HaHKufZs5FCSGazv5")
					g.Assert(repo.Hash).Equal("rad:z3gqcJUoA1n9HaHKufZs5FCSGazv5")
					g.Assert(repo.Owner).Equal(forge.Name())
					g.Assert(pipeline.Author).Equal("seb")
					g.Assert(pipeline.Event).Equal(model.EventPush)
					g.Assert(pipeline.Commit).Equal("ab6b2a2d318bf214d02f5427d541bbbf8140ab55")
					g.Assert(pipeline.Branch).Equal("ab6b2a2d318bf214d02f5427d541bbbf8140ab55")
					g.Assert(pipeline.Message).Equal("Update signed refs")
					g.Assert(pipeline.Timestamp).Equal(int64(1734006270))
					g.Assert(pipeline.Sender).Equal("radicle")
					g.Assert(pipeline.Email).Equal("radicle@localhost")
					g.Assert(pipeline.ForgeURL).Equal("rad:z3gqcJUoA1n9HaHKufZs5FCSGazv5/commits/1e7fa3584457f5894bfaed3b65918ec9d6668a4e")
					g.Assert(len(pipeline.ChangedFiles)).Equal(2)
				})
			})

			g.Describe("with an invalid push event type", func() {
				g.It("Should fail", func() {
					buf := bytes.NewBufferString(fixtures.HookPushPayloadInvalid)
					req, _ := http.NewRequest("POST", "/hook", buf)
					req.Header = http.Header{}
					req.Header.Set(internal.EventTypeHeaderKey, "push")
					_, _, err := forge.Hook(context.Background(), req)
					g.Assert(err).IsNotNil()
				})
			})

			g.Describe("with a valid patch event type", func() {
				g.It("Should succeed", func() {
					buf := bytes.NewBufferString(fixtures.HookPatchPayload)
					req, _ := http.NewRequest("POST", "/hook", buf)
					req.Header = http.Header{}
					req.Header.Set(internal.EventTypeHeaderKey, "patch")
					repo, pipeline, err := forge.Hook(context.Background(), req)
					g.Assert(err).IsNil()
					g.Assert(repo.ForgeRemoteID).Equal(model.ForgeRemoteID("z32iyJDyFLqvPFzwHm8YadK4HQ2EY"))
					g.Assert(repo.Name).Equal("mz-ci (rad:z32iyJDyFLqvPFzwHm8YadK4HQ2EY)")
					g.Assert(repo.FullName).Equal("mz-ci (rad:z32iyJDyFLqvPFzwHm8YadK4HQ2EY)")
					g.Assert(repo.Branch).Equal("master")
					g.Assert(repo.PREnabled).IsTrue()
					g.Assert(repo.Clone).Equal(forge.URL() + "/z32iyJDyFLqvPFzwHm8YadK4HQ2EY.git")
					g.Assert(repo.ForgeURL).Equal(forge.URL() + "/z32iyJDyFLqvPFzwHm8YadK4HQ2EY")
					g.Assert(repo.Hash).Equal("rad:z32iyJDyFLqvPFzwHm8YadK4HQ2EY")
					g.Assert(repo.Owner).Equal(forge.Name())
					g.Assert(pipeline.Author).Equal("michalis_server")
					g.Assert(pipeline.Event).Equal(model.EventPull)
					g.Assert(pipeline.Commit).Equal("274ac829adec365bb8a84b3673d8abff4a0ec1b6")
					g.Assert(pipeline.Branch).Equal("ed1fb3dea5e2db7d520664ecaf416ff0b6c72181")
					g.Assert(pipeline.Message).Equal("Woodpecker pipeline fix")
					g.Assert(pipeline.Timestamp).Equal(int64(1734006270))
					g.Assert(pipeline.Sender).Equal("did:key:z6MksMpnzPF48pk4XAnqVotKmfs2SE3bxA57UA8KL9DnWnY3")
					g.Assert(pipeline.Email).Equal("michalis_server")
					g.Assert(pipeline.ForgeURL).Equal("rad:z32iyJDyFLqvPFzwHm8YadK4HQ2EY/patches/ed1fb3dea5e2db7d520664ecaf416ff0b6c72181")
					g.Assert(len(pipeline.ChangedFiles)).Equal(1)
				})
			})

			g.Describe("with an invalid patch event type", func() {
				g.It("Should fail", func() {
					buf := bytes.NewBufferString(fixtures.HookPatchPayloadInvalid)
					req, _ := http.NewRequest("POST", "/hook", buf)
					req.Header = http.Header{}
					req.Header.Set(internal.EventTypeHeaderKey, "patch")
					_, _, err := forge.Hook(context.Background(), req)
					g.Assert(err).IsNotNil()
				})
			})
		})

		// Test OrgMembership
		g.Describe("When requesting a user's OrgMembership", func() {
			g.Describe("with some user", func() {
				g.It("Should return false", func() {
					user := &model.User{
						Login: "some-user",
					}
					membership, err := forge.OrgMembership(context.Background(), user, "radicle")
					g.Assert(err).IsNil()
					g.Assert(membership.Member).IsFalse()
					g.Assert(membership.Admin).IsFalse()
				})
			})
			g.Describe("with an admin user", func() {
				g.It("Should return true", func() {
					user := &model.User{
						Login: "radicle",
					}
					membership, err := forge.OrgMembership(context.Background(), user, "radicle")
					g.Assert(err).IsNil()
					g.Assert(membership.Member).IsTrue()
					g.Assert(membership.Admin).IsTrue()
				})
			})
		})

		// Test Org
		g.Describe("When requesting a user's OrgMembership", func() {
			g.Describe("with any user", func() {
				g.It("Should succeed", func() {
					user := &model.User{
						Login: "some-user",
					}
					org, err := forge.Org(context.Background(), user, "radicle")
					g.Assert(err).IsNil()
					g.Assert(org.Name).Equal("radicle")
					g.Assert(org.IsUser).IsTrue()
				})
			})
		})
	})
}
