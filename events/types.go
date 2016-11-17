package events

const (
	PING = iota
	PUSH
	PULL_REQUEST
	COMMIT_COMMENT
	CREATE
	DELETE
	DEPLOYMENT
	DEPLOYMENT_STATUS
	FORK
	GOLLUM
	ISSUE_COMMENT
	ISSUES
	MEMBER
	MEMBERSHIP
	PAGE_BUILD
	PUBLIC
	PULL_REQUEST_REVIEW_COMMENT
	PULL_REQUEST_REVIEW
	REPOSITORY
	RELEASE
	STATUS
	TEAM_ADD
	WATCH
)

var EventTypes = map[string]int{

	// Github Webhook ping
	"ping": PING,

	// Any Git push to a Repository, including editing tags or branches. Commits via API actions that update references
	// are also counted. This is the default event.
	"push": PUSH,

	// Any time a Pull Request is assigned, unassigned, labeled, unlabeled, opened, edited, closed, reopened, or
	// synchronized (updated due to a new push in the branch that the pull request is tracking).
	"pull_request": PULL_REQUEST,

	// Any time a Commit is commented on.
	"commit_comment": COMMIT_COMMENT,

	// Any time a Branch or Tag is created.
	"create": CREATE,

	// Any time a Branch or Tag is deleted.
	"delete": DELETE,

	// Any time a Repository has a new deployment created from the API.
	"deployment": DEPLOYMENT,

	// Any time a deployment for a Repository has a status update from the API.
	"deployment_status": DEPLOYMENT_STATUS,

	// Any time a Repository is forked.
	"fork": FORK,

	//Any time a Wiki page is updated.
	"gollum": GOLLUM,

	// Any time a comment on an issue is created, edited, or deleted.
	"issue_comment": ISSUE_COMMENT,

	//Any time an Issue is assigned, unassigned, labeled, unlabeled, opened, edited, closed, or reopened.
	"issues": ISSUES,

	// Any time a User is added as a collaborator to a Repository.
	"member": MEMBER,

	// Any time a User is added or removed from a team. Organization hooks only.
	"membership": MEMBERSHIP,

	// Any time a Pages site is built or results in a failed build.
	"page_build": PAGE_BUILD,

	// Any time a Repository changes from private to public.
	"public": PUBLIC,

	// Any time a comment on a Pull Request's unified diff is created, edited, or deleted (in the Files Changed tab).
	"pull_request_review_comment": PULL_REQUEST_REVIEW_COMMENT,

	// Any time a Pull Request Review is submitted.
	"pull_request_review": PULL_REQUEST_REVIEW,

	// Any time a Repository is created, deleted, made public, or made private.
	"repository": REPOSITORY,

	// Any time a Release is published in a Repository.
	"release": RELEASE,

	// Any time a Repository has a status update from the API
	"status": STATUS,

	// Any time a team is added or modified on a Repository.
	"team_add": TEAM_ADD,

	// Any time a User stars a Repository.
	"watch": WATCH,
}

type Event struct {
	ID   string
	Type int
	Name string
	Data map[string]interface{}
}