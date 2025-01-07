package store

import (
	"context"
	"github.com/keegancsmith/sqlf"
	"github.com/khulnasoft/khulnasoft/internal/database"
	ghtypes "github.com/khulnasoft/khulnasoft/internal/github_apps/types"
	"testing"

	"github.com/google/go-cmp/cmp"

	bt "github.com/khulnasoft/khulnasoft/internal/batches/testing"
	btypes "github.com/khulnasoft/khulnasoft/internal/batches/types"
	et "github.com/khulnasoft/khulnasoft/internal/encryption/testing"
	"github.com/khulnasoft/khulnasoft/internal/extsvc"
	"github.com/khulnasoft/khulnasoft/internal/extsvc/auth"
)

func testStoreSiteCredentials(t *testing.T, ctx context.Context, s *Store, _ bt.Clock) {
	appID := 1
	kind := ghtypes.SiteCredentialGitHubAppKind

	query := sqlf.Sprintf(
		"INSERT INTO github_apps (app_id, name, slug, base_url, client_id, client_secret, private_key, encryption_key_id, app_url, domain, kind) VALUES (%d, %s, %s, %s, %s, %s, %s, %s, DEFAULT, DEFAULT, %s)",
		appID, "name", "slug", "url", "clientID", "clientSecret", "privateKey", "encryptionKeyID", kind,
	)
	if err := s.Store.Exec(ctx, query); err != nil {
		t.Fatal(err)
	}

	var credentials []*btypes.SiteCredential

	creds := []struct {
		externalServiceID   string
		externalServiceType string
		githubAppID         int
	}{
		{
			externalServiceType: extsvc.TypeBitbucketServer,
			externalServiceID:   "https://someurl.test",
		},
		{
			externalServiceType: extsvc.TypeGitHub,
			externalServiceID:   "https://second.someurl.test",
			githubAppID:         appID,
		},
		{
			externalServiceType: extsvc.TypeGitHub,
			externalServiceID:   "https://someurl.test",
		},
		{
			externalServiceType: extsvc.TypeGitLab,
			externalServiceID:   "https://someurl.test",
		},
	}

	t.Run("Create", func(t *testing.T) {
		for _, c := range creds {
			sc := &btypes.SiteCredential{
				ExternalServiceType: c.externalServiceType,
				ExternalServiceID:   c.externalServiceID,
				GitHubAppID:         c.githubAppID,
			}

			token := &auth.OAuthBearerToken{Token: "123"}
			if err := s.CreateSiteCredential(ctx, sc, token); err != nil {
				t.Fatal(err)
			}

			if sc.ID == 0 {
				t.Fatal("id should not be zero")
			}
			if sc.CreatedAt.IsZero() {
				t.Fatal("CreatedAt should be set")
			}
			if sc.UpdatedAt.IsZero() {
				t.Fatal("UpdatedAt should be set")
			}
			credentials = append(credentials, sc)
		}
	})

	t.Run("Get", func(t *testing.T) {
		t.Run("ByID", func(t *testing.T) {
			want := credentials[0]
			opts := GetSiteCredentialOpts{ID: want.ID}

			have, err := s.GetSiteCredential(ctx, opts)
			if err != nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(have, want, et.CompareEncryptable); diff != "" {
				t.Fatal(diff)
			}
		})

		t.Run("ByKind-URL", func(t *testing.T) {
			want := credentials[0]
			opts := GetSiteCredentialOpts{
				ExternalServiceType: want.ExternalServiceType,
				ExternalServiceID:   want.ExternalServiceID,
			}

			have, err := s.GetSiteCredential(ctx, opts)
			if err != nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(have, want, et.CompareEncryptable); diff != "" {
				t.Fatal(diff)
			}
		})

		t.Run("NoResults", func(t *testing.T) {
			opts := GetSiteCredentialOpts{ID: 0xdeadbeef}

			_, have := s.GetSiteCredential(ctx, opts)
			want := ErrNoResults

			if have != want {
				t.Fatalf("have err %v, want %v", have, want)
			}
		})
	})

	t.Run("List", func(t *testing.T) {
		t.Run("All", func(t *testing.T) {
			cs, next, err := s.ListSiteCredentials(ctx, ListSiteCredentialsOpts{})
			if err != nil {
				t.Fatal(err)
			}
			if have, want := next, int64(0); have != want {
				t.Fatalf("have next %d, want %d", have, want)
			}

			have, want := cs, credentials
			if len(have) != len(want) {
				t.Fatalf("listed %d site credentials, want: %d", len(have), len(want))
			}

			if diff := cmp.Diff(have, want, et.CompareEncryptable); diff != "" {
				t.Fatal(diff)
			}
		})

		t.Run("WithLimit", func(t *testing.T) {
			for i := 1; i <= len(credentials); i++ {
				cs, next, err := s.ListSiteCredentials(ctx, ListSiteCredentialsOpts{LimitOpts: LimitOpts{Limit: i}})
				if err != nil {
					t.Fatal(err)
				}

				{
					have, want := next, int64(0)
					if i < len(credentials) {
						want = credentials[i].ID
					}

					if have != want {
						t.Fatalf("limit: %v: have next %v, want %v", i, have, want)
					}
				}

				{
					have, want := cs, credentials[:i]
					if len(have) != len(want) {
						t.Fatalf("listed %d site credentials, want: %d", len(have), len(want))
					}

					if diff := cmp.Diff(have, want, et.CompareEncryptable); diff != "" {
						t.Fatal(diff)
					}
				}
			}
		})
	})

	t.Run("Update", func(t *testing.T) {
		t.Run("Found", func(t *testing.T) {
			for _, cred := range credentials {
				if err := cred.SetAuthenticator(ctx, &auth.BasicAuthWithSSH{
					BasicAuth: auth.BasicAuth{
						Username: "foo",
						Password: "bar",
					},
					PrivateKey: "so private",
					PublicKey:  "so public",
					Passphrase: "probably written on a post-it",
				}); err != nil {
					t.Fatal(err)
				}

				if err := s.UpdateSiteCredential(ctx, cred); err != nil {
					t.Errorf("unexpected error: %+v", err)
				}

				if have, err := s.GetSiteCredential(ctx, GetSiteCredentialOpts{
					ID: cred.ID,
				}); err != nil {
					t.Errorf("error retrieving credential: %+v", err)
				} else if diff := cmp.Diff(have, cred, et.CompareEncryptable); diff != "" {
					t.Errorf("unexpected difference in credentials (-have +want):\n%s", diff)
				}
			}
		})
		t.Run("NotFound", func(t *testing.T) {
			cred := &btypes.SiteCredential{
				ID:         0xdeadbeef,
				Credential: database.NewEmptyCredential(),
			}
			if err := s.UpdateSiteCredential(ctx, cred); err == nil {
				t.Errorf("unexpected nil error")
			} else if err != ErrNoResults {
				t.Errorf("unexpected error: have=%v want=%v", err, ErrNoResults)
			}
		})
	})

	t.Run("Delete", func(t *testing.T) {
		t.Run("ByID", func(t *testing.T) {
			for _, cred := range credentials {
				if err := s.DeleteSiteCredential(ctx, cred.ID); err != nil {
					t.Fatal(err)
				}
			}
		})
		t.Run("NotFound", func(t *testing.T) {
			if err := s.DeleteSiteCredential(ctx, 0xdeadbeef); err == nil {
				t.Fatal("expected err but got nil")
			} else if err != ErrNoResults {
				t.Fatalf("invalid error %+v", err)
			}
		})
	})
}
