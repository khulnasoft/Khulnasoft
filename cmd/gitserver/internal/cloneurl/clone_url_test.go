package cloneurl

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"

	"github.com/khulnasoft/khulnasoft/internal/database/dbmocks"
	"github.com/khulnasoft/khulnasoft/internal/extsvc/awscodecommit"
	"github.com/khulnasoft/khulnasoft/internal/extsvc/azuredevops"
	"github.com/khulnasoft/khulnasoft/internal/extsvc/bitbucketcloud"
	"github.com/khulnasoft/khulnasoft/internal/extsvc/bitbucketserver"
	"github.com/khulnasoft/khulnasoft/internal/extsvc/gerrit"
	"github.com/khulnasoft/khulnasoft/internal/extsvc/github"
	"github.com/khulnasoft/khulnasoft/internal/extsvc/gitlab"
	"github.com/khulnasoft/khulnasoft/internal/extsvc/perforce"
	"github.com/khulnasoft/khulnasoft/internal/extsvc/phabricator"
	"github.com/khulnasoft/khulnasoft/internal/timeutil"
	"github.com/khulnasoft/khulnasoft/schema"
)

func TestAWSCodeCloneURLs(t *testing.T) {
	clock := timeutil.NewFakeClock(time.Now(), 0)
	now := clock.Now()

	repo := &awscodecommit.Repository{
		ARN:          "arn:aws:codecommit:us-west-1:999999999999:stripe-go",
		AccountID:    "999999999999",
		ID:           "f001337a-3450-46fd-b7d2-650c0EXAMPLE",
		Name:         "stripe-go",
		Description:  "The stripe-go lib",
		HTTPCloneURL: "https://git-codecommit.us-west-1.amazonaws.com/v1/repos/stripe-go",
		LastModified: &now,
	}

	cfg := schema.AWSCodeCommitConnection{
		GitCredentials: schema.AWSCodeCommitGitCredentials{
			Username: "username",
			Password: "password",
		},
	}

	got, err := awsCodeCloneURL(repo, &cfg)
	require.NoError(t, err)
	want := "https://username:password@git-codecommit.us-west-1.amazonaws.com/v1/repos/stripe-go"
	if got != want {
		t.Fatalf("wrong cloneURL, got: %q, want: %q", got, want)
	}
}

func TestAzureDevOpsCloneURL(t *testing.T) {
	cfg := schema.AzureDevOpsConnection{
		// the remote url used for clone has the username attached,
		// so we double-check that it gets replaced properly.
		Url:      "https://admin@dev.azure.com",
		Username: "admin",
		Token:    "pa$$word",
	}

	repo := &azuredevops.Repository{
		ID:        "test-project",
		RemoteURL: "https://sgtestazure@dev.azure.com/sgtestazure/sgtestazure/_git/sgtestazure",
	}

	got, err := azureDevOpsCloneURL(repo, &cfg)
	require.NoError(t, err)
	want := "https://admin:pa$$word@dev.azure.com/sgtestazure/sgtestazure/_git/sgtestazure"
	if got != want {
		t.Fatalf("wrong cloneURL, got: %q, want: %q", got, want)
	}
}

func TestBitbucketServerCloneURLs(t *testing.T) {
	repo := &bitbucketserver.Repo{
		ID:   1,
		Slug: "bar",
		Project: &bitbucketserver.Project{
			Key: "foo",
		},
	}

	cfg := schema.BitbucketServerConnection{
		Token:    "abc",
		Username: "username",
		Password: "password",
	}

	t.Run("ssh", func(t *testing.T) {
		repo.Links.Clone = []bitbucketserver.Link{
			// even if the first link is http, ssh should prevail
			{Name: "http", Href: "https://asdine@bitbucket.example.com/scm/sg/sourcegraph.git"},
			{Name: "ssh", Href: "ssh://git@bitbucket.example.com:7999/sg/sourcegraph.git"},
		}

		cfg.GitURLType = "ssh" // use ssh in the config as well

		got, err := bitbucketServerCloneURL(repo, &cfg)
		require.NoError(t, err)
		want := "ssh://git@bitbucket.example.com:7999/sg/sourcegraph.git"
		if got != want {
			t.Fatalf("wrong cloneURL, got: %q, want: %q", got, want)
		}
	})

	t.Run("http", func(t *testing.T) {
		// Second test: http
		repo.Links.Clone = []bitbucketserver.Link{
			{Name: "http", Href: "https://asdine@bitbucket.example.com/scm/sg/sourcegraph.git"},
		}

		got, err := bitbucketServerCloneURL(repo, &cfg)
		require.NoError(t, err)
		want := "https://username:abc@bitbucket.example.com/scm/sg/sourcegraph.git"
		if got != want {
			t.Fatalf("wrong cloneURL, got: %q, want: %q", got, want)
		}
	})

	t.Run("no token", func(t *testing.T) {
		// Third test: no token
		cfg.Token = ""

		got, err := bitbucketServerCloneURL(repo, &cfg)
		require.NoError(t, err)
		want := "https://username:password@bitbucket.example.com/scm/sg/sourcegraph.git"
		if got != want {
			t.Fatalf("wrong cloneURL, got: %q, want: %q", got, want)
		}
	})
}

func TestBitbucketCloudCloneURLs(t *testing.T) {
	repo := &bitbucketcloud.Repo{
		FullName: "sg/khulnasoft",
	}

	repo.Links.Clone = []bitbucketcloud.Link{
		{Name: "https", Href: "https://asdine@bitbucket.org/sg/sourcegraph.git"},
		{Name: "ssh", Href: "git@bitbucket.org/sg/sourcegraph.git"},
	}

	cfg := schema.BitbucketCloudConnection{
		Url:         "bitbucket.org",
		Username:    "username",
		AppPassword: "password",
	}

	t.Run("ssh", func(t *testing.T) {
		cfg.GitURLType = "ssh"

		got, err := bitbucketCloudCloneURL(repo, &cfg)
		require.NoError(t, err)
		want := "git@bitbucket.org:sg/sourcegraph.git"
		if got != want {
			t.Fatalf("wrong cloneURL, got: %q, want: %q", got, want)
		}
	})

	t.Run("http", func(t *testing.T) {
		cfg.GitURLType = "http"

		got, err := bitbucketCloudCloneURL(repo, &cfg)
		require.NoError(t, err)
		want := "https://username:password@bitbucket.org/sg/sourcegraph.git"
		if got != want {
			t.Fatalf("wrong cloneURL, got: %q, want: %q", got, want)
		}
	})
}

func TestGitHubCloneURLs(t *testing.T) {
	t.Run("empty repo.URL", func(t *testing.T) {
		_, err := githubCloneURL(context.Background(), dbmocks.NewMockDB(), &github.Repository{}, &schema.GitHubConnection{})
		got := fmt.Sprintf("%v", err)
		want := "empty repo.URL"
		if diff := cmp.Diff(want, got); diff != "" {
			t.Fatalf("Mismatch (-want +got):\n%s", diff)
		}
	})

	var repo github.Repository
	repo.NameWithOwner = "foo/bar"

	tests := []struct {
		InstanceUrl string
		RepoURL     string
		Token       string
		GitURLType  string
		Want        string
	}{
		{"https://github.com", "https://github.com/foo/bar", "", "", "https://github.com/foo/bar"},
		{"https://github.com", "https://github.com/foo/bar", "abcd", "", "https://oauth2:abcd@github.com/foo/bar"},
		{"https://github.com", "https://github.com/foo/bar", "abcd", "ssh", "git@github.com:foo/bar.git"},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("URL(%q) / Token(%q) / URLType(%q)", test.InstanceUrl, test.Token, test.GitURLType), func(t *testing.T) {
			cfg := schema.GitHubConnection{
				Url:        test.InstanceUrl,
				Token:      test.Token,
				GitURLType: test.GitURLType,
			}

			repo.URL = test.RepoURL

			got, err := githubCloneURL(context.Background(), dbmocks.NewMockDB(), &repo, &cfg)
			if err != nil {
				t.Fatal(err)
			}
			if got != test.Want {
				t.Fatalf("wrong cloneURL, got: %q, want: %q", got, test.Want)
			}
		})
	}
}

func TestGitLabCloneURLs(t *testing.T) {
	repo := &gitlab.Project{
		ProjectCommon: gitlab.ProjectCommon{
			ID:                1,
			PathWithNamespace: "foo/bar",
			SSHURLToRepo:      "git@gitlab.com:gitlab-org/gitaly.git",
			HTTPURLToRepo:     "https://gitlab.com/gitlab-org/gitaly.git",
		},
	}

	tests := []struct {
		Token      string
		GitURLType string
		TokenType  string
		Want       string
	}{
		{Want: "https://gitlab.com/gitlab-org/gitaly.git"},
		{Token: "abcd", Want: "https://git:abcd@gitlab.com/gitlab-org/gitaly.git"},
		{Token: "abcd", TokenType: "oauth", Want: "https://oauth2:abcd@gitlab.com/gitlab-org/gitaly.git"},
		{Token: "abcd", GitURLType: "ssh", Want: "git@gitlab.com:gitlab-org/gitaly.git"},
		{Token: "abcd", GitURLType: "ssh", Want: "git@gitlab.com:gitlab-org/gitaly.git"},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Token(%q) / URLType(%q)", test.Token, test.GitURLType), func(t *testing.T) {
			cfg := schema.GitLabConnection{
				Token:      test.Token,
				TokenType:  test.TokenType,
				GitURLType: test.GitURLType,
			}

			got, err := gitlabCloneURL(repo, &cfg)
			require.NoError(t, err)
			if got != test.Want {
				t.Fatalf("wrong cloneURL, got: %q, want: %q", got, test.Want)
			}
		})
	}
}

func TestGerritCloneURL(t *testing.T) {
	cfg := schema.GerritConnection{
		Url:      "https://gerrit.com",
		Username: "admin",
		Password: "pa$$word",
	}

	project := &gerrit.Project{
		ID:            "test-project",
		HTTPURLToRepo: "https://gerrit.com/a/test-project",
		SSHURLToRepo:  "ssh://gerrit-ssh.com:29418/a/test-project",
	}

	t.Run("HTTP", func(t *testing.T) {
		got, err := gerritCloneURL(project, &cfg)
		require.NoError(t, err)
		want := "https://admin:pa$$word@gerrit.com/a/test-project"
		if got != want {
			t.Fatalf("wrong cloneURL, got: %q, want: %q", got, want)
		}
	})
	t.Run("Legacy HTTP", func(t *testing.T) {
		project := &gerrit.Project{
			ID: "test-project",
			// Those values are not yet populated. This will only happen on the
			// next sync after we introduced them. This function should still work.
			HTTPURLToRepo: "",
			SSHURLToRepo:  "",
		}
		got, err := gerritCloneURL(project, &cfg)
		require.NoError(t, err)
		want := "https://admin:pa$$word@gerrit.com/a/test-project"
		if got != want {
			t.Fatalf("wrong cloneURL, got: %q, want: %q", got, want)
		}
	})
	t.Run("SSH", func(t *testing.T) {
		cfg.GitURLType = "ssh"
		got, err := gerritCloneURL(project, &cfg)
		require.NoError(t, err)
		want := "ssh://gerrit-ssh.com:29418/a/test-project"
		if got != want {
			t.Fatalf("wrong cloneURL, got: %q, want: %q", got, want)
		}
	})
}

func TestPerforceCloneURL(t *testing.T) {
	cfg := schema.PerforceConnection{
		P4Port:   "ssl:111.222.333.444:1666",
		P4User:   "admin",
		P4Passwd: "pa$$word",
	}

	repo := &perforce.Depot{
		Depot: "//Sourcegraph/",
	}

	got := perforceCloneURL(repo, &cfg)
	want := "perforce://admin:pa$$word@ssl:111.222.333.444:1666//Sourcegraph/"
	if got != want {
		t.Fatalf("wrong cloneURL, got: %q, want: %q", got, want)
	}
}

func TestPhabricatorCloneURL(t *testing.T) {
	meta := `
{
    "ID": 8,
    "VCS": "git",
    "Name": "testing",
    "PHID": "PHID-REPO-vl3v7n7jkzf5pjozoxuy",
    "URIs": [
        {
            "ID": "78",
            "PHID": "PHID-RURI-kmdhjr2u4ugjgaaatp4k",
            "Display": "git@gitolite.sgdev.org:testing",
            "Disabled": false,
            "Effective": "git@gitolite.sgdev.org:testing",
            "Normalized": "gitolite.sgdev.org/testing",
            "DateCreated": "2019-05-03T11:16:27Z",
            "DateModified": "0001-01-01T00:00:00Z",
            "BuiltinProtocol": "",
            "BuiltinIdentifier": ""
        },
        {
            "ID": "71",
            "PHID": "PHID-RURI-xu54xqjhvxwyxxzjoz63",
            "Display": "ssh://git@phabricator.sgdev.org/diffusion/8/test.git",
            "Disabled": false,
            "Effective": "ssh://git@phabricator.sgdev.org/diffusion/8/test.git",
            "Normalized": "phabricator.sgdev.org/diffusion/8",
            "DateCreated": "2019-05-03T11:16:06Z",
            "DateModified": "0001-01-01T00:00:00Z",
            "BuiltinProtocol": "ssh",
            "BuiltinIdentifier": "id"
        },
        {
            "ID": "70",
            "PHID": "PHID-RURI-3pstu43sbjncekq6rwqt",
            "Display": "ssh://git@phabricator.sgdev.org/source/test.git",
            "Disabled": false,
            "Effective": "ssh://git@phabricator.sgdev.org/source/test.git",
            "Normalized": "phabricator.sgdev.org/source/test",
            "DateCreated": "2019-05-03T11:16:06Z",
            "DateModified": "0001-01-01T00:00:00Z",
            "BuiltinProtocol": "ssh",
            "BuiltinIdentifier": "shortname"
        },
        {
            "ID": "69",
            "PHID": "PHID-RURI-5qh22baoby6u445k3nx5",
            "Display": "ssh://git@phabricator.sgdev.org/diffusion/TESTING/test.git",
            "Disabled": false,
            "Effective": "ssh://git@phabricator.sgdev.org/diffusion/TESTING/test.git",
            "Normalized": "phabricator.sgdev.org/diffusion/TESTING",
            "DateCreated": "2019-05-03T11:16:06Z",
            "DateModified": "0001-01-01T00:00:00Z",
            "BuiltinProtocol": "ssh",
            "BuiltinIdentifier": "callsign"
        }
    ],
    "Status": "active",
    "Callsign": "TESTING",
    "Shortname": "test",
    "EditPolicy": "admin",
    "ViewPolicy": "users",
    "DateCreated": "2019-05-03T11:16:06Z",
    "DateModified": "2019-08-08T14:45:57Z"
}
`

	repo := &phabricator.Repo{}
	err := json.Unmarshal([]byte(meta), repo)
	if err != nil {
		t.Fatal(err)
	}

	got, err := phabricatorCloneURL(repo, nil)
	require.NoError(t, err)
	want := "ssh://git@phabricator.sgdev.org/diffusion/8/test.git"

	if want != got {
		t.Fatalf("Want %q, got %q", want, got)
	}
}
