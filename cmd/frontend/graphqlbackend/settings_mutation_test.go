package graphqlbackend

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/khulnasoft/khulnasoft/internal/actor"
	"github.com/khulnasoft/khulnasoft/internal/api"
	"github.com/khulnasoft/khulnasoft/internal/auth"
	"github.com/khulnasoft/khulnasoft/internal/database/dbmocks"
	"github.com/khulnasoft/khulnasoft/internal/dotcom"
	"github.com/khulnasoft/khulnasoft/internal/gitserver"
	"github.com/khulnasoft/khulnasoft/internal/types"
)

func TestSettingsMutation_EditSettings(t *testing.T) {
	users := dbmocks.NewMockUserStore()
	users.GetByIDFunc.SetDefaultReturn(&types.User{ID: 1}, nil)
	users.GetByCurrentAuthUserFunc.SetDefaultReturn(&types.User{ID: 1, SiteAdmin: false}, nil)

	settings := dbmocks.NewMockSettingsStore()
	settings.GetLatestFunc.SetDefaultReturn(&api.Settings{ID: 1, Contents: "{}"}, nil)
	settings.CreateIfUpToDateFunc.SetDefaultHook(func(ctx context.Context, subject api.SettingsSubject, lastID, authorUserID *int32, contents string) (*api.Settings, error) {
		if want := `{
  "p": {
    "x": 123
  }
}`; contents != want {
			t.Errorf("got %q, want %q", contents, want)
		}
		return &api.Settings{ID: 2, Contents: contents}, nil
	})

	db := dbmocks.NewMockDB()
	db.UsersFunc.SetDefaultReturn(users)
	db.SettingsFunc.SetDefaultReturn(settings)

	RunTests(t, []*Test{
		{
			Context: actor.WithActor(context.Background(), &actor.Actor{UID: 1}),
			Schema:  mustParseGraphQLSchema(t, db),
			Query: `
				mutation($value: JSONValue) {
					settingsMutation(input: {subject: "VXNlcjox", lastID: 1}) {
						editSettings(edit: {keyPath: [{property: "p"}], value: $value}) {
							empty {
								alwaysNil
							}
						}
					}
				}
			`,
			Variables: map[string]any{"value": map[string]int{"x": 123}},
			ExpectedResult: `
				{
					"settingsMutation": {
						"editSettings": {
							"empty": null
						}
					}
				}
			`,
		},
	})
}

func TestSettingsMutation_OverwriteSettings(t *testing.T) {
	users := dbmocks.NewMockUserStore()
	users.GetByIDFunc.SetDefaultReturn(&types.User{ID: 1}, nil)
	users.GetByCurrentAuthUserFunc.SetDefaultReturn(&types.User{ID: 1, SiteAdmin: false}, nil)

	settings := dbmocks.NewMockSettingsStore()
	settings.GetLatestFunc.SetDefaultReturn(&api.Settings{ID: 1, Contents: "{}"}, nil)
	settings.CreateIfUpToDateFunc.SetDefaultHook(func(ctx context.Context, subject api.SettingsSubject, lastID, authorUserID *int32, contents string) (*api.Settings, error) {
		if want := `x`; contents != want {
			t.Errorf("got %q, want %q", contents, want)
		}
		return &api.Settings{ID: 2, Contents: contents}, nil
	})

	db := dbmocks.NewMockDB()
	db.UsersFunc.SetDefaultReturn(users)
	db.SettingsFunc.SetDefaultReturn(settings)

	RunTests(t, []*Test{
		{
			Context: actor.WithActor(context.Background(), &actor.Actor{UID: 1}),
			Schema:  mustParseGraphQLSchema(t, db),
			Query: `
				mutation($contents: String!) {
					settingsMutation(input: {subject: "VXNlcjox", lastID: 1}) {
						overwriteSettings(contents: $contents) {
							empty {
								alwaysNil
							}
						}
					}
				}
			`,
			Variables: map[string]any{"contents": "x"},
			ExpectedResult: `
				{
					"settingsMutation": {
						"overwriteSettings": {
							"empty": null
						}
					}
				}
			`,
		},
	})
}

func TestSettingsMutation(t *testing.T) {
	db := dbmocks.NewMockDB()
	t.Run("only allowed by authenticated user on Sourcegraph.com", func(t *testing.T) {
		users := dbmocks.NewMockUserStore()
		db.UsersFunc.SetDefaultReturn(users)

		dotcom.MockSourcegraphDotComMode(t, true)

		tests := []struct {
			name  string
			ctx   context.Context
			setup func()
		}{
			{
				name: "unauthenticated",
				ctx:  context.Background(),
				setup: func() {
					users.GetByIDFunc.SetDefaultReturn(&types.User{ID: 1}, nil)
				},
			},
			{
				name: "another user",
				ctx:  actor.WithActor(context.Background(), &actor.Actor{UID: 2}),
				setup: func() {
					users.GetByIDFunc.SetDefaultHook(func(ctx context.Context, id int32) (*types.User, error) {
						return &types.User{ID: id}, nil
					})
				},
			},
		}
		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				test.setup()

				_, err := newSchemaResolver(db, gitserver.NewTestClient(t), nil).SettingsMutation(
					test.ctx,
					&settingsMutationArgs{
						Input: &settingsMutationGroupInput{
							Subject: MarshalUserID(1),
						},
					},
				)
				got := fmt.Sprintf("%v", err)
				want := auth.ErrMustBeSiteAdminOrSameUser.Error()
				assert.Equal(t, want, got)
			})
		}
	})
}
