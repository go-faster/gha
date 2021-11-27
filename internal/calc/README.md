# Calc

Archive contains jsonlines stream of multiple types of events.

## Events

* [Event types](https://docs.github.com/en/developers/webhooks-and-events/events/github-event-types)
* [Common properties](https://docs.github.com/en/developers/webhooks-and-events/events/github-event-types#event-object-common-properties)

```
WatchEvent
IssueCommentEvent
PushEvent
CreateEvent
ForkEvent
PullRequestReviewEvent
IssuesEvent
PullRequestEvent
PullRequestReviewCommentEvent
DeleteEvent
CommitCommentEvent
ReleaseEvent
GollumEvent
MemberEvent
PublicEvent
```


### Watch

`WatchEvent`

When someone **stars** a repository. 
The type of activity is specified in the action property of the payload object. 

https://docs.github.com/en/rest/reference/activity#starring

```json
{
  "id": "13723948817",
  "type": "WatchEvent",
  "actor": {
    "id": 72303679,
    "login": "neiltcliu",
    "display_login": "neiltcliu",
    "gravatar_id": "",
    "url": "https://api.github.com/users/neiltcliu",
    "avatar_url": "https://avatars.githubusercontent.com/u/72303679?"
  },
  "repo": {
    "id": 50880060,
    "name": "HiddenRamblings/TagMo",
    "url": "https://api.github.com/repos/HiddenRamblings/TagMo"
  },
  "payload": {
    "action": "started"
  },
  "public": true,
  "created_at": "2020-10-03T12:00:00Z",
  "org": {
    "id": 30602223,
    "login": "HiddenRamblings",
    "gravatar_id": "",
    "url": "https://api.github.com/orgs/HiddenRamblings",
    "avatar_url": "https://avatars.githubusercontent.com/u/30602223?"
  }
}
```

## Issue comment

`IssueCommentEvent`

```json
{
  "id": "13723948818",
  "type": "IssueCommentEvent",
  "actor": {
    "id": 12938238,
    "login": "lukeb2e",
    "display_login": "lukeb2e",
    "gravatar_id": "",
    "url": "https://api.github.com/users/lukeb2e",
    "avatar_url": "https://avatars.githubusercontent.com/u/12938238?"
  },
  "repo": {
    "id": 168479288,
    "name": "ScoopInstaller/Main",
    "url": "https://api.github.com/repos/ScoopInstaller/Main"
  },
  "payload": {
    "action": "created",
    "issue": {
      "url": "https://api.github.com/repos/ScoopInstaller/Main/issues/826",
      "repository_url": "https://api.github.com/repos/ScoopInstaller/Main",
      "labels_url": "https://api.github.com/repos/ScoopInstaller/Main/issues/826/labels{/name}",
      "comments_url": "https://api.github.com/repos/ScoopInstaller/Main/issues/826/comments",
      "events_url": "https://api.github.com/repos/ScoopInstaller/Main/issues/826/events",
      "html_url": "https://github.com/ScoopInstaller/Main/pull/826",
      "id": 568143321,
      "node_id": "MDExOlB1bGxSZXF1ZXN0Mzc3NjI2ODAy",
      "number": 826,
      "title": "vlang: Add version 0.1.28.1",
      "user": {
        "login": "lukeb2e",
        "id": 12938238,
        "node_id": "MDQ6VXNlcjEyOTM4MjM4",
        "avatar_url": "https://avatars2.githubusercontent.com/u/12938238?v=4",
        "gravatar_id": "",
        "url": "https://api.github.com/users/lukeb2e",
        "html_url": "https://github.com/lukeb2e",
        "followers_url": "https://api.github.com/users/lukeb2e/followers",
        "following_url": "https://api.github.com/users/lukeb2e/following{/other_user}",
        "gists_url": "https://api.github.com/users/lukeb2e/gists{/gist_id}",
        "starred_url": "https://api.github.com/users/lukeb2e/starred{/owner}{/repo}",
        "subscriptions_url": "https://api.github.com/users/lukeb2e/subscriptions",
        "organizations_url": "https://api.github.com/users/lukeb2e/orgs",
        "repos_url": "https://api.github.com/users/lukeb2e/repos",
        "events_url": "https://api.github.com/users/lukeb2e/events{/privacy}",
        "received_events_url": "https://api.github.com/users/lukeb2e/received_events",
        "type": "User",
        "site_admin": false
      },
      "labels": [
        {
          "id": 1681995974,
          "node_id": "MDU6TGFiZWwxNjgxOTk1OTc0",
          "url": "https://api.github.com/repos/ScoopInstaller/Main/labels/please-review",
          "name": "please-review",
          "color": "a573d1",
          "default": false,
          "description": ""
        }
      ],
      "state": "open",
      "locked": false,
      "assignee": null,
      "assignees": [],
      "milestone": null,
      "comments": 1,
      "created_at": "2020-02-20T08:51:39Z",
      "updated_at": "2020-10-03T11:59:59Z",
      "closed_at": null,
      "author_association": "NONE",
      "active_lock_reason": null,
      "pull_request": {
        "url": "https://api.github.com/repos/ScoopInstaller/Main/pulls/826",
        "html_url": "https://github.com/ScoopInstaller/Main/pull/826",
        "diff_url": "https://github.com/ScoopInstaller/Main/pull/826.diff",
        "patch_url": "https://github.com/ScoopInstaller/Main/pull/826.patch"
      },
      "body": "Hi,\r\n\r\nI would like to add https://vlang.io to the bucket.\r\n\r\nThere are still some questions with this PR that need to be verified:\r\n\r\n- [ ] Name: v is very short. I followed the scheme used by golang here, but vlang might be a saner option (especially while searching).\r\n- [ ] Bucket: is this the correct bucket or should this be added to the extras bucket.\r\n\r\nLet me know if one of the above should be changed.\r\n\r\nKind regards",
      "performed_via_github_app": null
    },
    "comment": {
      "url": "https://api.github.com/repos/ScoopInstaller/Main/issues/comments/703092341",
      "html_url": "https://github.com/ScoopInstaller/Main/pull/826#issuecomment-703092341",
      "issue_url": "https://api.github.com/repos/ScoopInstaller/Main/issues/826",
      "id": 703092341,
      "node_id": "MDEyOklzc3VlQ29tbWVudDcwMzA5MjM0MQ==",
      "user": {
        "login": "lukeb2e",
        "id": 12938238,
        "node_id": "MDQ6VXNlcjEyOTM4MjM4",
        "avatar_url": "https://avatars2.githubusercontent.com/u/12938238?v=4",
        "gravatar_id": "",
        "url": "https://api.github.com/users/lukeb2e",
        "html_url": "https://github.com/lukeb2e",
        "followers_url": "https://api.github.com/users/lukeb2e/followers",
        "following_url": "https://api.github.com/users/lukeb2e/following{/other_user}",
        "gists_url": "https://api.github.com/users/lukeb2e/gists{/gist_id}",
        "starred_url": "https://api.github.com/users/lukeb2e/starred{/owner}{/repo}",
        "subscriptions_url": "https://api.github.com/users/lukeb2e/subscriptions",
        "organizations_url": "https://api.github.com/users/lukeb2e/orgs",
        "repos_url": "https://api.github.com/users/lukeb2e/repos",
        "events_url": "https://api.github.com/users/lukeb2e/events{/privacy}",
        "received_events_url": "https://api.github.com/users/lukeb2e/received_events",
        "type": "User",
        "site_admin": false
      },
      "created_at": "2020-10-03T11:59:59Z",
      "updated_at": "2020-10-03T11:59:59Z",
      "author_association": "NONE",
      "body": "Changed the file name as recommended.",
      "performed_via_github_app": null
    }
  },
  "public": true,
  "created_at": "2020-10-03T12:00:00Z",
  "org": {
    "id": 16618068,
    "login": "ScoopInstaller",
    "gravatar_id": "",
    "url": "https://api.github.com/orgs/ScoopInstaller",
    "avatar_url": "https://avatars.githubusercontent.com/u/16618068?"
  }
}
```

### Push

`PushEvent`

```json
{
  "id": "13723948819",
  "type": "PushEvent",
  "actor": {
    "id": 8517910,
    "login": "LombiqBot",
    "display_login": "LombiqBot",
    "gravatar_id": "",
    "url": "https://api.github.com/users/LombiqBot",
    "avatar_url": "https://avatars.githubusercontent.com/u/8517910?"
  },
  "repo": {
    "id": 264190944,
    "name": "Lombiq/Orchard.AngularJS",
    "url": "https://api.github.com/repos/Lombiq/Orchard.AngularJS"
  },
  "payload": {
    "push_id": 5792516651,
    "size": 0,
    "distinct_size": 0,
    "ref": "refs/heads/gh-pages",
    "head": "ae26d4621cec90d357177d8bd385db4599342c1d",
    "before": "ae26d4621cec90d357177d8bd385db4599342c1d",
    "commits": []
  },
  "public": true,
  "created_at": "2020-10-03T12:00:00Z",
  "org": {
    "id": 8158177,
    "login": "Lombiq",
    "gravatar_id": "",
    "url": "https://api.github.com/orgs/Lombiq",
    "avatar_url": "https://avatars.githubusercontent.com/u/8158177?"
  }
}
```

### Create

`CreateEvent`

```json
{
  "id": "13723948826",
  "type": "CreateEvent",
  "actor": {
    "id": 47211690,
    "login": "quanton314",
    "display_login": "quanton314",
    "gravatar_id": "",
    "url": "https://api.github.com/users/quanton314",
    "avatar_url": "https://avatars.githubusercontent.com/u/47211690?"
  },
  "repo": {
    "id": 300870023,
    "name": "quanton314/quantontest.github.io",
    "url": "https://api.github.com/repos/quanton314/quantontest.github.io"
  },
  "payload": {
    "ref": "main",
    "ref_type": "branch",
    "master_branch": "main",
    "description": null,
    "pusher_type": "user"
  },
  "public": true,
  "created_at": "2020-10-03T12:00:00Z"
}
```

### Fork

`ForkEvent`

```json
{
  "id": "13723948831",
  "type": "ForkEvent",
  "actor": {
    "id": 23740251,
    "login": "anaghkanungo7",
    "display_login": "anaghkanungo7",
    "gravatar_id": "",
    "url": "https://api.github.com/users/anaghkanungo7",
    "avatar_url": "https://avatars.githubusercontent.com/u/23740251?"
  },
  "repo": {
    "id": 299953825,
    "name": "eddiejaoude/Hacktoberfest-FirstPR",
    "url": "https://api.github.com/repos/eddiejaoude/Hacktoberfest-FirstPR"
  },
  "payload": {
    "forkee": {
      "id": 300870036,
      "node_id": "MDEwOlJlcG9zaXRvcnkzMDA4NzAwMzY=",
      "name": "Hacktoberfest-FirstPR",
      "full_name": "anaghkanungo7/Hacktoberfest-FirstPR",
      "private": false,
      "owner": {
        "login": "anaghkanungo7",
        "id": 23740251,
        "node_id": "MDQ6VXNlcjIzNzQwMjUx",
        "avatar_url": "https://avatars1.githubusercontent.com/u/23740251?v=4",
        "gravatar_id": "",
        "url": "https://api.github.com/users/anaghkanungo7",
        "html_url": "https://github.com/anaghkanungo7",
        "followers_url": "https://api.github.com/users/anaghkanungo7/followers",
        "following_url": "https://api.github.com/users/anaghkanungo7/following{/other_user}",
        "gists_url": "https://api.github.com/users/anaghkanungo7/gists{/gist_id}",
        "starred_url": "https://api.github.com/users/anaghkanungo7/starred{/owner}{/repo}",
        "subscriptions_url": "https://api.github.com/users/anaghkanungo7/subscriptions",
        "organizations_url": "https://api.github.com/users/anaghkanungo7/orgs",
        "repos_url": "https://api.github.com/users/anaghkanungo7/repos",
        "events_url": "https://api.github.com/users/anaghkanungo7/events{/privacy}",
        "received_events_url": "https://api.github.com/users/anaghkanungo7/received_events",
        "type": "User",
        "site_admin": false
      },
      "html_url": "https://github.com/anaghkanungo7/Hacktoberfest-FirstPR",
      "description": "Make your first PR for hacktoberfest! Beginner Friendly!",
      "fork": true,
      "url": "https://api.github.com/repos/anaghkanungo7/Hacktoberfest-FirstPR",
      "forks_url": "https://api.github.com/repos/anaghkanungo7/Hacktoberfest-FirstPR/forks",
      "keys_url": "https://api.github.com/repos/anaghkanungo7/Hacktoberfest-FirstPR/keys{/key_id}",
      "collaborators_url": "https://api.github.com/repos/anaghkanungo7/Hacktoberfest-FirstPR/collaborators{/collaborator}",
      "teams_url": "https://api.github.com/repos/anaghkanungo7/Hacktoberfest-FirstPR/teams",
      "hooks_url": "https://api.github.com/repos/anaghkanungo7/Hacktoberfest-FirstPR/hooks",
      "issue_events_url": "https://api.github.com/repos/anaghkanungo7/Hacktoberfest-FirstPR/issues/events{/number}",
      "events_url": "https://api.github.com/repos/anaghkanungo7/Hacktoberfest-FirstPR/events",
      "assignees_url": "https://api.github.com/repos/anaghkanungo7/Hacktoberfest-FirstPR/assignees{/user}",
      "branches_url": "https://api.github.com/repos/anaghkanungo7/Hacktoberfest-FirstPR/branches{/branch}",
      "tags_url": "https://api.github.com/repos/anaghkanungo7/Hacktoberfest-FirstPR/tags",
      "blobs_url": "https://api.github.com/repos/anaghkanungo7/Hacktoberfest-FirstPR/git/blobs{/sha}",
      "git_tags_url": "https://api.github.com/repos/anaghkanungo7/Hacktoberfest-FirstPR/git/tags{/sha}",
      "git_refs_url": "https://api.github.com/repos/anaghkanungo7/Hacktoberfest-FirstPR/git/refs{/sha}",
      "trees_url": "https://api.github.com/repos/anaghkanungo7/Hacktoberfest-FirstPR/git/trees{/sha}",
      "statuses_url": "https://api.github.com/repos/anaghkanungo7/Hacktoberfest-FirstPR/statuses/{sha}",
      "languages_url": "https://api.github.com/repos/anaghkanungo7/Hacktoberfest-FirstPR/languages",
      "stargazers_url": "https://api.github.com/repos/anaghkanungo7/Hacktoberfest-FirstPR/stargazers",
      "contributors_url": "https://api.github.com/repos/anaghkanungo7/Hacktoberfest-FirstPR/contributors",
      "subscribers_url": "https://api.github.com/repos/anaghkanungo7/Hacktoberfest-FirstPR/subscribers",
      "subscription_url": "https://api.github.com/repos/anaghkanungo7/Hacktoberfest-FirstPR/subscription",
      "commits_url": "https://api.github.com/repos/anaghkanungo7/Hacktoberfest-FirstPR/commits{/sha}",
      "git_commits_url": "https://api.github.com/repos/anaghkanungo7/Hacktoberfest-FirstPR/git/commits{/sha}",
      "comments_url": "https://api.github.com/repos/anaghkanungo7/Hacktoberfest-FirstPR/comments{/number}",
      "issue_comment_url": "https://api.github.com/repos/anaghkanungo7/Hacktoberfest-FirstPR/issues/comments{/number}",
      "contents_url": "https://api.github.com/repos/anaghkanungo7/Hacktoberfest-FirstPR/contents/{+path}",
      "compare_url": "https://api.github.com/repos/anaghkanungo7/Hacktoberfest-FirstPR/compare/{base}...{head}",
      "merges_url": "https://api.github.com/repos/anaghkanungo7/Hacktoberfest-FirstPR/merges",
      "archive_url": "https://api.github.com/repos/anaghkanungo7/Hacktoberfest-FirstPR/{archive_format}{/ref}",
      "downloads_url": "https://api.github.com/repos/anaghkanungo7/Hacktoberfest-FirstPR/downloads",
      "issues_url": "https://api.github.com/repos/anaghkanungo7/Hacktoberfest-FirstPR/issues{/number}",
      "pulls_url": "https://api.github.com/repos/anaghkanungo7/Hacktoberfest-FirstPR/pulls{/number}",
      "milestones_url": "https://api.github.com/repos/anaghkanungo7/Hacktoberfest-FirstPR/milestones{/number}",
      "notifications_url": "https://api.github.com/repos/anaghkanungo7/Hacktoberfest-FirstPR/notifications{?since,all,participating}",
      "labels_url": "https://api.github.com/repos/anaghkanungo7/Hacktoberfest-FirstPR/labels{/name}",
      "releases_url": "https://api.github.com/repos/anaghkanungo7/Hacktoberfest-FirstPR/releases{/id}",
      "deployments_url": "https://api.github.com/repos/anaghkanungo7/Hacktoberfest-FirstPR/deployments",
      "created_at": "2020-10-03T11:59:59Z",
      "updated_at": "2020-09-30T14:52:35Z",
      "pushed_at": "2020-10-03T09:52:27Z",
      "git_url": "git://github.com/anaghkanungo7/Hacktoberfest-FirstPR.git",
      "ssh_url": "git@github.com:anaghkanungo7/Hacktoberfest-FirstPR.git",
      "clone_url": "https://github.com/anaghkanungo7/Hacktoberfest-FirstPR.git",
      "svn_url": "https://github.com/anaghkanungo7/Hacktoberfest-FirstPR",
      "homepage": null,
      "size": 4,
      "stargazers_count": 0,
      "watchers_count": 0,
      "language": null,
      "has_issues": false,
      "has_projects": true,
      "has_downloads": true,
      "has_wiki": true,
      "has_pages": false,
      "forks_count": 0,
      "mirror_url": null,
      "archived": false,
      "disabled": false,
      "open_issues_count": 0,
      "license": null,
      "forks": 0,
      "open_issues": 0,
      "watchers": 0,
      "default_branch": "main",
      "public": true
    }
  },
  "public": true,
  "created_at": "2020-10-03T12:00:00Z"
}
```
