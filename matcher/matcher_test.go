package matcher_test

import (
	"github.com/homedepot/github-webhook/matcher"
	"github.com/homedepot/github-webhook/configuration"
	"github.com/homedepot/github-webhook/events"
	"github.com/stretchr/testify/assert"
	"testing"
	"encoding/json"
)


////ID string
////Type int
////Name string
//Data map[string]interface{}


func TestMatcherEventNoRules(t *testing.T) {

	trigger := configuration.Trigger{
		Name: "trigger",
		Event: "cookie",
	}

	event := events.Event{
		Name: "cookie",
	}
	assert.True(t, matcher.Matches(trigger, event))


}

func TestMatcherDifferentEvent(t *testing.T) {

	trigger := configuration.Trigger{
		Name: "trigger",
		Event: "cookie",
	}

	event := events.Event{
		Name: "cookie_monster",
	}
	assert.False(t, matcher.Matches(trigger, event))


}

func TestMatcherRef(t *testing.T) {

	bytes := []byte(`{
	"ref": "refs/heads/master"
}`)

	var data map[string]interface{}
	if err := json.Unmarshal(bytes, &data); err != nil {
		assert.FailNow(t, err.Error())
	}

	rules := make(map[string]string)
	rules["ref"] = "refs/heads/master"

	trigger := configuration.Trigger{
		Name: "trigger",
		Event: "cookie",
		Rules: rules,
	}

	event := events.Event{
		Name: "cookie",
		Data: data,
	}
	assert.True(t, matcher.Matches(trigger, event))


}

func TestMatcherBool(t *testing.T) {

	bytes := []byte(`{
	"ref": "refs/heads/master",
	"base_ref": null,
	"distinct": true
}`)

	var data map[string]interface{}
	if err := json.Unmarshal(bytes, &data); err != nil {
		assert.FailNow(t, err.Error())
	}

	rules := make(map[string]string)
	rules["distinct"] = "true"

	trigger := configuration.Trigger{
		Name: "trigger",
		Event: "cookie",
		Rules: rules,
	}

	event := events.Event{
		Name: "cookie",
		Data: data,
	}
	assert.True(t, matcher.Matches(trigger, event))


}

func TestMatcherMultipleRules(t *testing.T) {

	bytes := []byte(`{
	"ref": "refs/heads/master",
	"base_ref": null,
	"distinct": true
}`)

	var data map[string]interface{}
	if err := json.Unmarshal(bytes, &data); err != nil {
		assert.FailNow(t, err.Error())
	}

	rules := make(map[string]string)
	rules["distinct"] = "true"
	rules["ref"] = "refs/heads/master"

	trigger := configuration.Trigger{
		Name: "trigger",
		Event: "cookie",
		Rules: rules,
	}

	event := events.Event{
		Name: "cookie",
		Data: data,
	}
	assert.True(t, matcher.Matches(trigger, event))


}

func TestMatcherNested(t *testing.T) {

	bytes := []byte(`{
	"ref": "refs/heads/master",
	"base_ref": null,
	"commits": [{
	"message": "test4",
	"timestamp": "2016-09-30T22:31:27-04:00",
	"author": {
	"name": "Mendez, Marcos",
	"email": "spam@somewhere.com",
	"username": "username"
}}]}`)

	var data map[string]interface{}
	if err := json.Unmarshal(bytes, &data); err != nil {
		assert.FailNow(t, err.Error())
	}

	rules := make(map[string]string)
	rules["commits.author.username"] = "username"

	trigger := configuration.Trigger{
		Name: "trigger",
		Event: "cookie",
		Rules: rules,
	}

	event := events.Event{
		Name: "cookie",
		Data: data,
	}
	assert.True(t, matcher.Matches(trigger, event))

}

func TestMatcherNestedStringList(t *testing.T) {

	bytes := []byte(`{
	"ref": "refs/heads/master",
	"base_ref": null,
	"commits": [{
	"message": "test4",
	"timestamp": "2016-09-30T22:31:27-04:00",
	"author": {
	"name": "Mendez, Marcos",
	"email": "spam@somewhere.com",
	"username": "username"
	},
 	"committer": {
 	"name": "GitHub Enterprise",
 	"email": "noreply@somewhere.com"
 	},
 	"added": [],
 	"removed": [],
 	"modified": ["README.md"]}]}`)

	var data map[string]interface{}
	if err := json.Unmarshal(bytes, &data); err != nil {
		assert.FailNow(t, err.Error())
	}

	rules := make(map[string]string)
	rules["commits.modified"] = "README.md"

	trigger := configuration.Trigger{
		Name: "trigger",
		Event: "cookie",
		Rules: rules,
	}

	event := events.Event{
		Name: "cookie",
		Data: data,
	}
	assert.True(t, matcher.Matches(trigger, event))

}

func TestMatcherFullEvent(t *testing.T) {

	bytes := []byte(` {
 	"ref": "refs/heads/master",
 	"before": "2e6fc73ecbf39516e4e53aa7582a5f2db2992f15",
 	"after": "a469370482a7079a145bcd2f4f735de48a1ac84c",
 	"created": false,
 	"deleted": false,
 	"forced": false,
 	"base_ref": null,
 	"compare": "https://somewhere.com/username/test/compare/2e6fc73ecbf3...a469370482a7",
 	"commits": [{
 	"id": "a469370482a7079a145bcd2f4f735de48a1ac84c",
 	"tree_id": "90deb4da79aa11d972b8ae7d65b95e2b8c44ec11",
 	"distinct": true,
 	"message": "test4",
 	"timestamp": "2016-09-30T22:31:27-04:00",
 	"url": "https://somewhere.com/username/test/commit/a469370482a7079a145bcd2f4f735de48a1ac84c",
 	"author": {
 	"name": "Mendez, Marcos",
 	"email": "spam@somewhere.com",
 	"username": "username"
 	},
 	"committer": {
 	"name": "GitHub Enterprise",
 	"email": "noreply@somewhere.com"
 	},
 	"added": [],
 	"removed": [],
 	"modified": ["README.md"]
 	}],
 	"head_commit": {
 	"id": "a469370482a7079a145bcd2f4f735de48a1ac84c",
 	"tree_id": "90deb4da79aa11d972b8ae7d65b95e2b8c44ec11",
 	"distinct": true,
 	"message": "test4",
 	"timestamp": "2016-09-30T22:31:27-04:00",
 	"url": "https://somewhere.com/username/test/commit/a469370482a7079a145bcd2f4f735de48a1ac84c",
 	"author": {
 	"name": "Mendez, Marcos",
 	"email": "spam@somewhere.com",
 	"username": "username"
 	},
 	"committer": {
 	"name": "GitHub Enterprise",
 	"email": "noreply@somewhere.com"
 	},
 	"added": [],
 	"removed": [],
 	"modified": ["README.md"]
 	},
 	"repository": {
 	"id": 13066,
 	"name": "test",
 	"full_name": "username/test",
 	"owner": {
 	"name": "username",
 	"email": "spam@somewhere.com"
 	},
 	"private": false,
 	"html_url": "https://somewhere.com/username/test",
 	"description": "test",
 	"fork": false,
 	"url": "https://somewhere.com/username/test",
 	"forks_url": "https://somewhere.com/api/v3/repos/username/test/forks",
 	"keys_url": "https://somewhere.com/api/v3/repos/username/test/keys{/key_id}",
 	"collaborators_url": "https://somewhere.com/api/v3/repos/username/test/collaborators{/collaborator}",
 	"teams_url": "https://somewhere.com/api/v3/repos/username/test/teams",
 	"hooks_url": "https://somewhere.com/api/v3/repos/username/test/hooks",
 	"issue_events_url": "https://somewhere.com/api/v3/repos/username/test/issues/events{/number}",
 	"events_url": "https://somewhere.com/api/v3/repos/username/test/events",
 	"assignees_url": "https://somewhere.com/api/v3/repos/username/test/assignees{/user}",
 	"branches_url": "https://somewhere.com/api/v3/repos/username/test/branches{/branch}",
 	"tags_url": "https://somewhere.com/api/v3/repos/username/test/tags",
 	"blobs_url": "https://somewhere.com/api/v3/repos/username/test/git/blobs{/sha}",
 	"git_tags_url": "https://somewhere.com/api/v3/repos/username/test/git/tags{/sha}",
 	"git_refs_url": "https://somewhere.com/api/v3/repos/username/test/git/refs{/sha}",
 	"trees_url": "https://somewhere.com/api/v3/repos/username/test/git/trees{/sha}",
 	"statuses_url": "https://somewhere.com/api/v3/repos/username/test/statuses/{sha}",
 	"languages_url": "https://somewhere.com/api/v3/repos/username/test/languages",
 	"stargazers_url": "https://somewhere.com/api/v3/repos/username/test/stargazers",
 	"contributors_url": "https://somewhere.com/api/v3/repos/username/test/contributors",
 	"subscribers_url": "https://somewhere.com/api/v3/repos/username/test/subscribers",
 	"subscription_url": "https://somewhere.com/api/v3/repos/username/test/subscription",
 	"commits_url": "https://somewhere.com/api/v3/repos/username/test/commits{/sha}",
 	"git_commits_url": "https://somewhere.com/api/v3/repos/username/test/git/commits{/sha}",
 	"comments_url": "https://somewhere.com/api/v3/repos/username/test/comments{/number}",
 	"issue_comment_url": "https://somewhere.com/api/v3/repos/username/test/issues/comments{/number}",
 	"contents_url": "https://somewhere.com/api/v3/repos/username/test/contents/{+path}",
 	"compare_url": "https://somewhere.com/api/v3/repos/username/test/compare/{base}...{head}",
 	"merges_url": "https://somewhere.com/api/v3/repos/username/test/merges",
 	"archive_url": "https://somewhere.com/api/v3/repos/username/test/{archive_format}{/ref}",
 	"downloads_url": "https://somewhere.com/api/v3/repos/username/test/downloads",
 	"issues_url": "https://somewhere.com/api/v3/repos/username/test/issues{/number}",
 	"pulls_url": "https://somewhere.com/api/v3/repos/username/test/pulls{/number}",
 	"milestones_url": "https://somewhere.com/api/v3/repos/username/test/milestones{/number}",
 	"notifications_url": "https://somewhere.com/api/v3/repos/username/test/notifications{?since,all,participating}",
 	"labels_url": "https://somewhere.com/api/v3/repos/username/test/labels{/name}",
 	"releases_url": "https://somewhere.com/api/v3/repos/username/test/releases{/id}",
 	"deployments_url": "https://somewhere.com/api/v3/repos/username/test/deployments",
 	"created_at": 1475287106,
 	"updated_at": "2016-10-01T01:58:26Z",
 	"pushed_at": 1475289092,
 	"git_url": "git://somewhere.com/username/test.git",
 	"ssh_url": "git@somewhere.com:username/test.git",
 	"clone_url": "https://somewhere.com/username/test.git",
 	"svn_url": "https://somewhere.com/username/test",
 	"homepage": null,
 	"size": 0,
 	"stargazers_count": 0,
 	"watchers_count": 0,
 	"language": null,
 	"has_issues": true,
 	"has_downloads": true,
 	"has_wiki": true,
 	"has_pages": false,
 	"forks_count": 0,
 	"mirror_url": null,
 	"open_issues_count": 0,
 	"forks": 0,
 	"open_issues": 0,
 	"watchers": 0,
 	"default_branch": "master",
 	"stargazers": 0,
 	"master_branch": "master"
 	},
 	"pusher": {
 	"name": "username",
 	"email": "spam@somewhere.com"
 	},
 	"sender": {
 	"login": "username",
 	"id": 146,
 	"avatar_url": "https://somewhere.com/avatars/u/146?",
 	"gravatar_id": "",
 	"url": "https://somewhere.com/api/v3/users/username",
 	"html_url": "https://somewhere.com/username",
 	"followers_url": "https://somewhere.com/api/v3/users/username/followers",
 	"following_url": "https://somewhere.com/api/v3/users/username/following{/other_user}",
 	"gists_url": "https://somewhere.com/api/v3/users/username/gists{/gist_id}",
 	"starred_url": "https://somewhere.com/api/v3/users/username/starred{/owner}{/repo}",
 	"subscriptions_url": "https://somewhere.com/api/v3/users/username/subscriptions",
 	"organizations_url": "https://somewhere.com/api/v3/users/username/orgs",
 	"repos_url": "https://somewhere.com/api/v3/users/username/repos",
 	"events_url": "https://somewhere.com/api/v3/users/username/events{/privacy}",
 	"received_events_url": "https://somewhere.com/api/v3/users/username/received_events",
 	"type": "User",
 	"site_admin": false,
 	"ldap_dn": "DC=somewhere,DC=com"
 	}
 }`)

	var data map[string]interface{}
	if err := json.Unmarshal(bytes, &data); err != nil {
		assert.FailNow(t, err.Error())
	}

	rules := make(map[string]string)
	rules["sender.ldap_dn"] = "DC=somewhere,DC=com"

	trigger := configuration.Trigger{
		Name: "trigger",
		Event: "cookie",
		Rules: rules,
	}

	event := events.Event{
		Name: "cookie",
		Data: data,
	}
	assert.True(t, matcher.Matches(trigger, event))

}