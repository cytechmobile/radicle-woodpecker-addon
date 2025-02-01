package fixtures

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Handler returns an http.Handler that is capable of handling a variety of mock
// Radicle requests and returning mock responses.
func Handler() http.Handler {
	gin.SetMode(gin.TestMode)

	e := gin.New()
	e.GET("/api/v1/node", getNodeInfo)
	e.GET("/api/v1/sessions/:session_id", getSession)
	e.GET("/api/v1/projects", getProjects)
	e.GET("/api/v1/projects/:project_id", getProject)
	e.GET("/api/v1/projects/:project_id/blob/:commit/:file", getProjectCommitFile)
	e.GET("/api/v1/projects/:project_id/tree/:commit/", getProjectCommitTree)
	e.GET("/api/v1/projects/:project_id/commits", getProjectCommits)
	e.GET("/api/v1/projects/:project_id/patches", getProjectPatches)
	e.PATCH("/api/v1/projects/:project_id/patches/:patch_id", addProjectPatchComment)
	e.POST("/api/v1/projects/:project_id/webhooks", addProjectWebhook)
	e.DELETE("/api/v1/projects/:project_id/webhooks", removeProjectWebhook)
	return e
}

func getNodeInfo(c *gin.Context) {
	c.String(200, nodePayload)
}

func getSession(c *gin.Context) {
	switch c.Param("session_id") {
	case "not_found":
		c.String(404, notFound)
	case "unauthed_sess_id":
		c.String(200, sessionUnauthorizedPayload)
	default:
		c.String(200, sessionPayload)
	}
}

func getProjects(c *gin.Context) {
	switch c.Query("page") {
	case "0":
		c.String(200, projectsPayloadPage0)
	case "1":
		c.String(200, projectsPayloadPage1)
	default:
		c.String(200, emptyPayload)
	}
}

func getProject(c *gin.Context) {
	switch c.Param("project_id") {
	case "not_found":
		c.String(404, notFound)
	default:
		c.String(200, projectPayload)
	}
}

func getProjectCommitFile(c *gin.Context) {
	if c.Param("commit") != "the_commit_id" {
		c.String(404, notFound)
		return
	}
	switch c.Param("project_id") {
	case "not_found":
		c.String(404, notFound)
	default:
		c.String(200, projectCommitFilePayload)
	}
}

func getProjectCommitTree(c *gin.Context) {
	if c.Param("commit") != "the_commit_id" {
		c.String(404, notFound)
		return
	}
	switch c.Param("project_id") {
	case "not_found":
		c.String(404, notFound)
	default:
		c.String(200, projectCommitTreePayload)
	}
}

func getProjectCommits(c *gin.Context) {
	switch c.Param("project_id") {
	case "not_found":
		c.String(404, notFound)
	default:
		c.String(200, projectCommitsPayload)
	}
}

func getProjectPatches(c *gin.Context) {
	switch c.Param("project_id") {
	case "not_found":
		c.String(404, notFound)
	default:
		if c.Query("page") == "0" {
			c.String(200, projectPullRequestsPayload)
		} else {
			c.String(200, emptyPayload)
		}
	}
}

func addProjectPatchComment(c *gin.Context) {
	switch c.Param("project_id") {
	case "invalid":
		c.String(400, invalid)
	case "not_found":
		c.String(404, notFound)
	default:
		switch c.Param("patch_id") {
		case "patchID":
			c.String(200, emptyPayload)
		default:
			c.String(400, invalid)
		}
	}
}

func addProjectWebhook(c *gin.Context) {
	switch c.Param("project_id") {
	case "invalid":
		c.String(400, invalid)
	case "not_found":
		c.String(404, notFound)
	default:
		c.String(200, emptyPayload)
	}
}

func removeProjectWebhook(c *gin.Context) {
	switch c.Param("project_id") {
	case "invalid":
		c.String(400, invalid)
	case "not_found":
		c.String(404, notFound)
	default:
		c.String(200, emptyPayload)
	}
}

const nodePayload = `
{
	"id": "someid",
	"config": {
		"alias": "myalias"
	}
}
`
const sessionPayload = `
{
	"sessionId": "session_id",
	"status": "authorized",
	"publicKey": "a_pub_key",
	"alias": "myalias",
	"issuedAt": 1234567890,
	"expiresAt": 1234567891
}
`

const sessionUnauthorizedPayload = `
{
	"sessionId": "session_id",
	"status": "unauthorized",
	"publicKey": "a_pub_key",
	"alias": "myalias",
	"issuedAt": 1234567890,
	"expiresAt": 1234567891
}
`
const emptyPayload = `[]`

const projectsPayloadPage0 = `
[
	{
		"name": "a-project",
		"description": "a description",
		"defaultBranch": "main",
		"delegates": [
			{
				"id":"did:key:z6MksFqXN3Yhqk8pTJdUGLwATkRfQvwZXPqR2qMEhbS9wzpT",
				"alias":"cloudhead"
			},
			{
				"id":"did:key:z6MksFqXN3Yhqk8pTJdUGLwATkRfQvwZXPqR2qMEhbS9waSd",
				"alias":"michalis"
			}
		],
		"head": "00bfa9b18be32001481334126c311c4a327dff2e",
		"id": "rad:a_project"
	},
	{
		"name": "b-project",
		"description": "b description",
		"defaultBranch": "master",
		"delegates": [
			{
				"id":"did:key:z6MksFqXN3Yhqk8pTJdUGLwATkRfQvwZXPqR2qMEhbS9wzpT",
				"alias":"cloudhead"
			},
			{
				"id":"did:key:z6MksFqXN3Yhqk8pTJdUGLwATkRfQvwZXPqR2qMEhbS9waSd",
				"alias":"michalis"
			}
		],
		"head": "00bfa9b18be32001481334126c311c4a327dff2f",
		"id": "rad:b_project"
	}
]
`

const projectsPayloadPage1 = `
[
	{
		"name": "c-project",
		"description": "c description",
		"defaultBranch": "main",
		"delegates": [
			{
				"id":"did:key:z6MksFqXN3Yhqk8pTJdUGLwATkRfQvwZXPqR2qMEhbS9wzpT",
				"alias":"cloudhead"
			},
			{
				"id":"did:key:z6MksFqXN3Yhqk8pTJdUGLwATkRfQvwZXPqR2qMEhbS9waSd",
				"alias":"michalis"
			}
		],
		"head": "00bfa9b18be32001481334126c311c4a327dff2e",
		"id": "rad:c_project"
	},
	{
		"name": "d-project",
		"description": "d description",
		"defaultBranch": "master",
		"delegates": [
			{
				"id":"did:key:z6MksFqXN3Yhqk8pTJdUGLwATkRfQvwZXPqR2qMEhbS9wzpT",
				"alias":"cloudhead"
			},
			{
				"id":"did:key:z6MksFqXN3Yhqk8pTJdUGLwATkRfQvwZXPqR2qbS9wQweRty",
				"alias":"kostas"
			}
		],
		"head": "00bfa9b18be32001481334126c311c4a327dff2f",
		"id": "rad:d_project"
	}
]
`

const projectPayload = `
{
	"name": "a-project",
	"description": "a description",
	"defaultBranch": "main",
	"delegates": [
		{
			"id":"did:key:z6MksFqXN3Yhqk8pTJdUGLwATkRfQvwZXPqR2qMEhbS9wzpT",
			"alias":"cloudhead"
		},
		{
			"id":"did:key:z6MksFqXN3Yhqk8pTJdUGLwATkRfQvwZXPqR2qMEhbS9waSd",
			"alias":"michalis"
		}
	],
	"head": "00bfa9b18be32001481334126c311c4a327dff2e",
	"id": "rad:valid_project_id"
}
`

const projectCommitFilePayload = `
{
	"binary": false,
	"name": "file_name.md",
	"content": "file content",
	"path": "file_path/file_name.md"
}
`

const projectCommitTreePayload = `
{
	"entries": [
		{
			"path": "Readme.md",
			"name": "source",
			"kind": "blob"
		},
		{
			"path": "cargo-checksum.json",
			"name": "cargo-checksum.json",
			"kind": "blob"
		},
		{
			"path": "debian",
			"name": "build-deb",
			"kind": "tree"
		}
	]
}
`

const projectCommitsPayload = `
[
	{
		"id": "00bfa9b18be32001481334126c311c4a327dff2e",
		"parents": [
			"5bb95551460527ce7c24640683d4c0d5cd55a52e"
		]
	},
	{
		"id": "00bfa9b18be32001481334126c311c4a327dff2f",
		"parents": [
			"5bb95551460527ce7c24640683d4c0d5cd55a52f"
		]
	}
]
`

const projectPullRequestsPayload = `
	[
   {
      "id":"c7eee5122d0467aec5e71c228c958f9c79fe17c9",
      "author":{
         "id":"did:key:z6MksFqXN3Yhqk8pTJdUGLwATkRfQvwZXPqR2qMEhbS9wzpT",
         "alias":"cloudhead"
      },
      "title":"Use repository consistently over project",
      "state":{
         "status":"open",
				 "conflicts": []
      },
			"before": "beforeSHA1",
			"after": "afterSHA1",
      "target":"delegates",
      "labels":[

      ],
      "merges":[

      ],
      "assignees":[

      ],
      "revisions":[
         {
            "id":"c7eee5122d0467aec5e71c228c958f9c79fe17c9",
            "author":{
               "id":"did:key:z6MksFqXN3Yhqk8pTJdUGLwATkRfQvwZXPqR2qMEhbS9wzpT",
               "alias":"cloudhead"
            },
            "description":"We were using the two pretty interchangeably. \"Project\" should only be\nused to refer to the repository payload. \"Repository\" should be used\nwhen referring to the resource that is fetched, synced, cloned, checked\nout etc.",
            "edits":[
               {
                  "author":{
                     "id":"did:key:z6MksFqXN3Yhqk8pTJdUGLwATkRfQvwZXPqR2qMEhbS9wzpT",
                     "alias":"cloudhead"
                  },
                  "body":"We were using the two pretty interchangeably. \"Project\" should only be\nused to refer to the repository payload. \"Repository\" should be used\nwhen referring to the resource that is fetched, synced, cloned, checked\nout etc.",
                  "timestamp":1705680429,
                  "embeds":[

                  ]
               }
            ],
            "base":"7b3d380ceb5e268b28c2ada97dde0652d7ecb35b",
            "oid":"a859b04341cb793ef0725010acf3720c0c7a5acc",
            "refs":[

            ],
            "discussions":[

            ],
            "timestamp":1705680429,
            "reviews":[

            ]
         },
         {
            "id":"6c56a5fd9fe54bb61b11427b0ab7e98d0e2a92c8",
            "author":{
               "id":"did:key:z6MksFqXN3Yhqk8pTJdUGLwATkRfQvwZXPqR2qMEhbS9wzpT",
               "alias":"cloudhead"
            },
            "description":"Rebase.",
            "edits":[
               {
                  "author":{
                     "id":"did:key:z6MksFqXN3Yhqk8pTJdUGLwATkRfQvwZXPqR2qMEhbS9wzpT",
                     "alias":"cloudhead"
                  },
                  "body":"Rebase.",
                  "timestamp":1705680585,
                  "embeds":[

                  ]
               }
            ],
            "base":"dc8561847d734e558738974be13015b1981aa1cb",
            "oid":"b2ebce45d83697b12fd7b780e2a0e4b4f14e0544",
            "refs":[
               "refs/heads/patches/c7eee5122d0467aec5e71c228c958f9c79fe17c9"
            ],
            "discussions":[

            ],
            "timestamp":1705680585,
            "reviews":[

            ]
         }
      ]
   },
   {
      "id":"c969b0642425b6fdd34ade8df866cc0d330e8cc5",
      "author":{
         "id":"did:key:z6MkkfM3tPXNPrPevKr3uSiQtHPuwnNhu2yUVjgd2jXVsVz5",
         "alias":"sebastinez"
      },
      "title":"httpd: Add project sync status endpoint",
      "state":{
         "status":"open"
      },
			"before": "beforeSHA2",
			"after": "afterSHA2",
      "target":"delegates",
      "labels":[

      ],
      "merges":[

      ],
      "assignees":[

      ],
      "revisions":[
         {
            "id":"c969b0642425b6fdd34ade8df866cc0d330e8cc5",
            "author":{
               "id":"did:key:z6MkkfM3tPXNPrPevKr3uSiQtHPuwnNhu2yUVjgd2jXVsVz5",
               "alias":"sebastinez"
            },
            "description":"",
            "edits":[
               {
                  "author":{
                     "id":"did:key:z6MkkfM3tPXNPrPevKr3uSiQtHPuwnNhu2yUVjgd2jXVsVz5",
                     "alias":"sebastinez"
                  },
                  "body":"",
                  "timestamp":1705579427,
                  "embeds":[

                  ]
               }
            ],
            "base":"d139762f4dbbb4ac72e72535bafaaa05e8284125",
            "oid":"fa4e151fef24c4673e4aab24f23414510c12fd9f",
            "refs":[
               "refs/heads/patches/b4d6ad7dbb949ebb7ebc66b2d9a4c130ae9f5c68",
               "refs/heads/patches/c969b0642425b6fdd34ade8df866cc0d330e8cc5"
            ],
            "discussions":[

            ],
            "timestamp":1705579427,
            "reviews":[

            ]
         }
      ]
   }
]
`

const HookPushPayloadInvalid = `
{
	"author": {
		"id": "did:key:z6MkkfM3tPXNPrPevKr3uSiQtHPuwnNhu2yUVjgd2jXVsVz5",
		"alias": "sebastinez"
	},
	"before": "6f3905801e6aeffb116c6e629d693c09f6622491",
	"after": "ab6b2a2d318bf214d02f5427d541bbbf8140ab55",
	"commits": [],
	"repository": {
		"id": "rad:z3gqcJUoA1n9HaHKufZs5FCSGazv5",
		"name": "heartwood",
		"description": "Radicle Heartwood Protocol & Stack",
		"private": false,
		"default_branch": "master",
		"url": "rad:z3gqcJUoA1n9HaHKufZs5FCSGazv5",
		"clone_url": "http://127.0.0.1:8080/rad:z3gqcJUoA1n9HaHKufZs5FCSGazv5.git",
		"delegates": ["did:key:z6MksFqXN3Yhqk8pTJdUGLwATkRfQvwZXPqR2qMEhbS9wzpT"]
	}
}
`

const HookPushPayload = `
{
	"author": {
		"id": "did:key:z6MkkfM3tPXNPrPevKr3uSiQtHPuwnNhu2yUVjgd2jXVsVz5",
		"alias": "seb"
	},
	"before": "6f3905801e6aeffb116c6e629d693c09f6622491",
	"after": "ab6b2a2d318bf214d02f5427d541bbbf8140ab55",
	"commits": [{
		"id": "ab6b2a2d318bf214d02f5427d541bbbf8140ab55",
		"title": "Update signed refs",
		"message": "",
		"timestamp": "2024-12-12T12:24:30Z",
		"url": "rad:z3gqcJUoA1n9HaHKufZs5FCSGazv5/commits/ab6b2a2d318bf214d02f5427d541bbbf8140ab55",
		"author": {
			"name": "radicle",
			"email": "radicle@localhost"
		},
		"added": [],
		"modified": ["refs"],
		"removed": []
	}, {
		"id": "1e7fa3584457f5894bfaed3b65918ec9d6668a4e",
		"title": "Update signed refs",
		"message": "",
		"timestamp": "2024-12-12T12:24:30Z",
		"url": "rad:z3gqcJUoA1n9HaHKufZs5FCSGazv5/commits/1e7fa3584457f5894bfaed3b65918ec9d6668a4e",
		"author": {
			"name": "radicle",
			"email": "radicle@localhost"
		},
		"added": [],
		"modified": [],
		"removed": ["signature"]
	}],
	"repository": {
		"id": "rad:z3gqcJUoA1n9HaHKufZs5FCSGazv5",
		"name": "heartwood",
		"description": "Radicle Heartwood Protocol & Stack",
		"private": false,
		"default_branch": "master",
		"url": "rad:z3gqcJUoA1n9HaHKufZs5FCSGazv5",
		"clone_url": "http://127.0.0.1:8080/rad:z3gqcJUoA1n9HaHKufZs5FCSGazv5.git",
		"delegates": ["did:key:z6MksFqXN3Yhqk8pTJdUGLwATkRfQvwZXPqR2qMEhbS9wzpT"]
	}
}
`

const HookPatchPayloadInvalid = `
{
	"action": "created",
	"patch": {
		"id": "ed1fb3dea5e2db7d520664ecaf416ff0b6c72181",
		"author": {
			"id": "did:key:z6MksMpnzPF48pk4XAnqVotKmfs2SE3bxA57UA8KL9DnWnY3",
			"alias": "michalis_server"
		},
		"title": "Woodpecker pipeline fix",
		"state": {
			"status": "open",
			"conflicts": []
		},
		"before": "ef25208520566bfb96fb00b16ea7c8bd98ffeb8e",
		"after": "274ac829adec365bb8a84b3673d8abff4a0ec1b6",
		"commits": [{
			"id": "274ac829adec365bb8a84b3673d8abff4a0ec1b6",
			"title": "Fix pipeline",
			"message": "",
			"timestamp": 1705650791,
			"url": "rad:z32iyJDyFLqvPFzwHm8YadK4HQ2EY/commits/274ac829adec365bb8a84b3673d8abff4a0ec1b6",
			"author": {
				"name": "Michalis Zampetakis",
				"email": "mzampetakis@gmail.com"
			},
			"added": [],
			"modified": [".woodpecker.yaml"],
			"removed": []
		}],
		"url": "rad:z32iyJDyFLqvPFzwHm8YadK4HQ2EY/patches/ed1fb3dea5e2db7d520664ecaf416ff0b6c72181",
		"target": "ef25208520566bfb96fb00b16ea7c8bd98ffeb8e",
		"labels": [],
		"assignees": [],
		"revisions": []
	},
	"repository": {
		"id": "rad:z32iyJDyFLqvPFzwHm8YadK4HQ2EY",
		"name": "mz-ci",
		"description": "",
		"private": false,
		"default_branch": "master",
		"url": "rad:z32iyJDyFLqvPFzwHm8YadK4HQ2EY",
		"clone_url": "http://127.0.0.1:8080/rad:z32iyJDyFLqvPFzwHm8YadK4HQ2EY.git",
		"delegates": ["did:key:z6MksMpnzPF48pk4XAnqVotKmfs2SE3bxA57UA8KL9DnWnY3"]
	}
}
`

const HookPatchPayload = `
{
	"action": "created",
	"patch": {
		"id": "ed1fb3dea5e2db7d520664ecaf416ff0b6c72181",
		"author": {
			"id": "did:key:z6MksMpnzPF48pk4XAnqVotKmfs2SE3bxA57UA8KL9DnWnY3",
			"alias": "michalis_server"
		},
		"title": "Woodpecker pipeline fix",
		"state": {
			"status": "open",
			"conflicts": []
		},
		"before": "ef25208520566bfb96fb00b16ea7c8bd98ffeb8e",
		"after": "274ac829adec365bb8a84b3673d8abff4a0ec1b6",
		"commits": [{
			"id": "274ac829adec365bb8a84b3673d8abff4a0ec1b6",
			"title": "Fix pipeline",
			"message": "",
			"timestamp": "2024-12-12T12:24:30Z",
			"url": "rad:z32iyJDyFLqvPFzwHm8YadK4HQ2EY/commits/274ac829adec365bb8a84b3673d8abff4a0ec1b6",
			"author": {
				"name": "Michalis Zampetakis",
				"email": "mzampetakis@gmail.com"
			},
			"added": [],
			"modified": [".woodpecker.yaml"],
			"removed": []
		}],
		"url": "rad:z32iyJDyFLqvPFzwHm8YadK4HQ2EY/patches/ed1fb3dea5e2db7d520664ecaf416ff0b6c72181",
		"target": "ef25208520566bfb96fb00b16ea7c8bd98ffeb8e",
		"labels": [],
		"assignees": [],
		"revisions": [{
			"id": "ed1fb3dea5e2db7d520664ecaf416ff0b6c72181",
			"author": {
				"id": "did:key:z6MksMpnzPF48pk4XAnqVotKmfs2SE3bxA57UA8KL9DnWnY3",
				"alias": "michalis_server"
			},
			"description": "",
			"base": "ef25208520566bfb96fb00b16ea7c8bd98ffeb8e",
			"oid": "274ac829adec365bb8a84b3673d8abff4a0ec1b6",
			"timestamp": "2024-12-12T12:24:30Z"
		}]
	},
	"repository": {
		"id": "rad:z32iyJDyFLqvPFzwHm8YadK4HQ2EY",
		"name": "mz-ci",
		"description": "",
		"private": false,
		"default_branch": "master",
		"url": "rad:z32iyJDyFLqvPFzwHm8YadK4HQ2EY",
		"clone_url": "http://127.0.0.1:8080/rad:z32iyJDyFLqvPFzwHm8YadK4HQ2EY.git",
		"delegates": ["did:key:z6MksMpnzPF48pk4XAnqVotKmfs2SE3bxA57UA8KL9DnWnY3"]
	}
}
`

const notFound = `
{
	error: "Not Found",
	code: 404
}
`

const invalid = `
{
	error: "Invalid",
	code: 400
}
`
