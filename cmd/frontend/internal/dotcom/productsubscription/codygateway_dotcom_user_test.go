package productsubscription_test

import (
	"context"
	"testing"
	"time"

	"github.com/hexops/autogold/v2"
	"github.com/sourcegraph/log/logtest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/khulnasoft/khulnasoft/cmd/frontend/graphqlbackend"
	"github.com/khulnasoft/khulnasoft/cmd/frontend/internal/cody"
	"github.com/khulnasoft/khulnasoft/cmd/frontend/internal/dotcom/productsubscription"
	"github.com/khulnasoft/khulnasoft/cmd/frontend/internal/ssc"
	"github.com/khulnasoft/khulnasoft/internal/accesstoken"
	"github.com/khulnasoft/khulnasoft/internal/actor"
	"github.com/khulnasoft/khulnasoft/internal/audit/audittest"
	"github.com/khulnasoft/khulnasoft/internal/authz"
	"github.com/khulnasoft/khulnasoft/internal/conf"
	"github.com/khulnasoft/khulnasoft/internal/database"
	"github.com/khulnasoft/khulnasoft/internal/database/dbtest"
	"github.com/khulnasoft/khulnasoft/internal/extsvc"
	"github.com/khulnasoft/khulnasoft/internal/featureflag"
	"github.com/khulnasoft/khulnasoft/internal/rbac"
	"github.com/khulnasoft/khulnasoft/internal/types"
	"github.com/khulnasoft/khulnasoft/lib/pointers"
	"github.com/khulnasoft/khulnasoft/schema"
)

func TestCodyGatewayDotcomUserResolver(t *testing.T) {
	chatOverrideLimit := 200
	codeOverrideLimit := 400

	tru := true
	cfg := &conf.Unified{
		SiteConfiguration: schema.SiteConfiguration{
			CodyEnabled: &tru,
			LicenseKey:  "asdf",
			Completions: &schema.Completions{
				Provider:                         "khulnasoft",
				PerUserCodeCompletionsDailyLimit: 20,
				PerUserDailyLimit:                10,
			},
		},
	}
	conf.Mock(cfg)
	defer func() {
		conf.Mock(nil)
	}()

	ctx := context.Background()
	db := database.NewDB(logtest.Scoped(t), dbtest.NewDB(t))

	// User with default rate limits
	adminUser, err := db.Users().Create(ctx, database.NewUser{Username: "admin", EmailIsVerified: true, Email: "admin@test.com"})
	require.NoError(t, err)

	// Verified User with default rate limits
	verifiedUser, err := db.Users().Create(ctx, database.NewUser{Username: "verified", EmailIsVerified: true, Email: "verified@test.com"})
	require.NoError(t, err)

	// Unverified User with default rate limits
	unverifiedUser, err := db.Users().Create(ctx, database.NewUser{Username: "unverified", EmailIsVerified: false, Email: "christopher.warwick@sourcegraph.com", EmailVerificationCode: "CODE"})
	require.NoError(t, err)

	// User with rate limit overrides
	overrideUser, err := db.Users().Create(ctx, database.NewUser{Username: "override", EmailIsVerified: true, Email: "override@test.com"})
	require.NoError(t, err)
	err = db.Users().SetChatCompletionsQuota(context.Background(), overrideUser.ID, pointers.Ptr(chatOverrideLimit))
	require.NoError(t, err)
	err = db.Users().SetCodeCompletionsQuota(context.Background(), overrideUser.ID, pointers.Ptr(codeOverrideLimit))
	require.NoError(t, err)

	tests := []struct {
		name        string
		user        *types.User
		wantChat    graphqlbackend.BigInt
		wantCode    graphqlbackend.BigInt
		wantEnabled bool
	}{
		{
			name:        "admin user",
			user:        adminUser,
			wantChat:    graphqlbackend.BigInt(cfg.Completions.PerUserDailyLimit),
			wantCode:    graphqlbackend.BigInt(cfg.Completions.PerUserCodeCompletionsDailyLimit),
			wantEnabled: true,
		},
		{
			name:        "verified user default limits",
			user:        verifiedUser,
			wantChat:    graphqlbackend.BigInt(cfg.Completions.PerUserDailyLimit),
			wantCode:    graphqlbackend.BigInt(cfg.Completions.PerUserCodeCompletionsDailyLimit),
			wantEnabled: true,
		},
		{
			name:        "unverified user",
			user:        unverifiedUser,
			wantChat:    0,
			wantCode:    0,
			wantEnabled: false,
		},
		{
			name:        "override user",
			user:        overrideUser,
			wantChat:    graphqlbackend.BigInt(chatOverrideLimit),
			wantCode:    graphqlbackend.BigInt(codeOverrideLimit),
			wantEnabled: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			// Create an admin context to use for the request
			adminContext := actor.WithActor(context.Background(), actor.FromActualUser(adminUser))

			// Generate a dotcom api Token for the test user
			_, dotcomToken, err := db.AccessTokens().Create(context.Background(), test.user.ID, []string{authz.ScopeUserAll}, test.name, test.user.ID, time.Time{})
			require.NoError(t, err)
			// convert token into a gateway token
			gatewayToken, err := accesstoken.GenerateDotcomUserGatewayAccessToken(dotcomToken)
			require.NoError(t, err)

			logger, exportLogs := logtest.Captured(t)

			// Make request from the admin checking the test user's token
			r := productsubscription.CodyGatewayDotcomUserResolver{Logger: logger, DB: db}
			userResolver, err := r.CodyGatewayDotcomUserByToken(adminContext, &graphqlbackend.CodyGatewayUsersByAccessTokenArgs{Token: gatewayToken})
			require.NoError(t, err)

			chat, err := userResolver.CodyGatewayAccess().ChatCompletionsRateLimit(adminContext)
			require.NoError(t, err)
			if chat != nil {
				require.Equal(t, test.wantChat, chat.Limit())
			} else {
				require.False(t, test.wantEnabled) // If there is no limit make sure it's expected to be disabled
			}

			code, err := userResolver.CodyGatewayAccess().CodeCompletionsRateLimit(adminContext)
			require.NoError(t, err)
			if chat != nil {
				require.Equal(t, test.wantCode, code.Limit())
			} else {
				require.False(t, test.wantEnabled) // If there is no limit make sure it's expected to be disabled
			}

			assert.Equal(t, test.wantEnabled, userResolver.CodyGatewayAccess().Enabled())

			// A user was resolved in this test case, we should have an audit log
			assert.True(t, exportLogs().Contains(func(l logtest.CapturedLog) bool {
				fields, ok := audittest.ExtractAuditFields(l)
				if !ok {
					return ok
				}
				return fields.Entity == "dotcom-codygatewayuser" && fields.Action == "access"
			}))
		})
	}
}

func TestCodyGatewayDotcomUserResolverUserNotFound(t *testing.T) {
	ctx := context.Background()
	db := database.NewDB(logtest.Scoped(t), dbtest.NewDB(t))

	// admin user to make request
	adminUser, err := db.Users().Create(ctx, database.NewUser{Username: "admin", EmailIsVerified: true, Email: "admin@test.com"})
	require.NoError(t, err)

	// Create an admin context to use for the request
	adminContext := actor.WithActor(context.Background(), actor.FromActualUser(adminUser))

	r := productsubscription.CodyGatewayDotcomUserResolver{Logger: logtest.Scoped(t), DB: db}
	_, err = r.CodyGatewayDotcomUserByToken(adminContext, &graphqlbackend.CodyGatewayUsersByAccessTokenArgs{Token: "NOT_A_TOKEN"})

	_, got := err.(productsubscription.ErrDotcomUserNotFound)
	assert.True(t, got, "should have error type ErrDotcomUserNotFound")
}

func TestCodyGatewayDotcomUserResolverRequestAccess(t *testing.T) {
	ctx := context.Background()
	db := database.NewDB(logtest.Scoped(t), dbtest.NewDB(t))

	// Admin
	adminUser, err := db.Users().Create(ctx, database.NewUser{Username: "admin", EmailIsVerified: true, Email: "admin@test.com"})
	require.NoError(t, err)

	// Not Admin with RBAC
	notAdminUser, err := db.Users().Create(ctx, database.NewUser{Username: "verified", EmailIsVerified: true, Email: "verified@test.com"})
	require.NoError(t, err)
	ns, action, err := rbac.ParsePermissionDisplayName(rbac.ProductSubscriptionsReadPermission)
	require.NoError(t, err)
	perm, err := db.Permissions().Create(ctx, database.CreatePermissionOpts{
		Namespace: ns,
		Action:    action,
	})
	require.NoError(t, err)
	role, err := db.Roles().Create(ctx, "SUBSCRIPTIONS_READER", false)
	require.NoError(t, err)
	err = db.RolePermissions().Assign(ctx, database.AssignRolePermissionOpts{
		PermissionID: perm.ID,
		RoleID:       role.ID,
	})
	require.NoError(t, err)
	err = db.UserRoles().Assign(ctx, database.AssignUserRoleOpts{
		UserID: notAdminUser.ID,
		RoleID: role.ID,
	})
	require.NoError(t, err)

	// No admin, no RBAC
	noAccessUser, err := db.Users().Create(ctx, database.NewUser{Username: "nottheone", EmailIsVerified: true, Email: "nottheone@test.com"})
	require.NoError(t, err)

	// cody user
	codyUser, err := db.Users().Create(ctx, database.NewUser{Username: "cody", EmailIsVerified: true, Email: "cody@test.com"})
	require.NoError(t, err)
	// Generate a token for the cody user
	_, codyUserApiToken, err := db.AccessTokens().Create(context.Background(), codyUser.ID, []string{authz.ScopeUserAll}, "cody", codyUser.ID, time.Time{})
	require.NoError(t, err)
	codyUserGatewayToken, err := accesstoken.GenerateDotcomUserGatewayAccessToken(codyUserApiToken)
	require.NoError(t, err)

	tests := []struct {
		name    string
		user    *types.User
		wantErr autogold.Value
	}{
		{
			name:    "admin user",
			user:    adminUser,
			wantErr: nil,
		},
		{
			name:    "RBAC reader role",
			user:    notAdminUser,
			wantErr: nil,
		},
		{
			name:    "not admin or RBAC reader role user",
			user:    noAccessUser,
			wantErr: autogold.Expect("unauthorized"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			// Create a request context from the user
			userContext := actor.WithActor(context.Background(), actor.FromActualUser(test.user))

			// Make request from the test user
			r := productsubscription.CodyGatewayDotcomUserResolver{Logger: logtest.Scoped(t), DB: db}
			_, err := r.CodyGatewayDotcomUserByToken(userContext, &graphqlbackend.CodyGatewayUsersByAccessTokenArgs{Token: codyUserGatewayToken})

			if test.wantErr != nil {
				require.Error(t, err)
				test.wantErr.Equal(t, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCodyGatewayCompletionsRateLimit(t *testing.T) {
	ctx := context.Background()
	db := database.NewDB(logtest.Scoped(t), dbtest.NewDB(t))

	override := 20
	perUserDailyLimit := 30
	perCommunityUserChatMonthlyLLMRequestLimit := 40
	perProUserChatDailyLLMRequestLimit := 50
	oneDayInSeconds := int32(60 * 60 * 24)

	tru := true
	cfg := &conf.Unified{
		SiteConfiguration: schema.SiteConfiguration{
			CodyEnabled: &tru,
			LicenseKey:  "asdf",
			Completions: &schema.Completions{
				Provider: "khulnasoft",
				PerCommunityUserChatMonthlyLLMRequestLimit: perCommunityUserChatMonthlyLLMRequestLimit,
				PerProUserChatDailyLLMRequestLimit:         perProUserChatDailyLLMRequestLimit,
				PerUserDailyLimit:                          perUserDailyLimit,
			},
		},
	}
	conf.Mock(cfg)
	defer func() {
		conf.Mock(nil)
	}()

	// User with an override
	userWithOverrides, err := db.Users().Create(ctx, database.NewUser{Username: "override", EmailIsVerified: true, Email: "override@test.com"})
	require.NoError(t, err)
	err = db.Users().SetChatCompletionsQuota(context.Background(), userWithOverrides.ID, pointers.Ptr(override))
	require.NoError(t, err)

	// Cody SSC - Free user
	sscFreeUser, err := db.Users().Create(ctx, database.NewUser{Username: "ssc-free", EmailIsVerified: true, Email: "ssc-free@test.com"})
	require.NoError(t, err)
	sscFreeUserExternalAccount, err := db.UserExternalAccounts().Insert(ctx, &extsvc.Account{
		UserID: sscFreeUser.ID,
		AccountSpec: extsvc.AccountSpec{
			AccountID:   "123",
			ServiceType: "openidconnect",
			ServiceID:   ssc.GetSAMSServiceID(),
		},
	})
	require.NoError(t, err)

	// Cody SSC - Pro user
	sscProUser, err := db.Users().Create(ctx, database.NewUser{Username: "ssc-pro", EmailIsVerified: true, Email: "ssc-pro@test.com"})
	require.NoError(t, err)
	err = db.Users().ChangeCodyPlan(ctx, sscProUser.ID, true)
	require.NoError(t, err)
	sscProUserExternalAccount, err := db.UserExternalAccounts().Insert(ctx, &extsvc.Account{
		UserID: sscProUser.ID,
		AccountSpec: extsvc.AccountSpec{
			AccountID:   "456",
			ServiceType: "openidconnect",
			ServiceID:   ssc.GetSAMSServiceID(),
		},
	})
	require.NoError(t, err)

	tests := []struct {
		name                            string
		pro                             bool
		user                            *types.User
		externalAccount                 *extsvc.Account
		wantChatLimit                   graphqlbackend.BigInt
		wantChatLimitInterval           int32
		wantCodeCompletionLimit         graphqlbackend.BigInt
		wantCodeCompletionLimitInterval int32
	}{
		{
			name:                            "override",
			user:                            userWithOverrides,
			wantChatLimit:                   graphqlbackend.BigInt(override),
			wantChatLimitInterval:           oneDayInSeconds,
			wantCodeCompletionLimit:         graphqlbackend.BigInt(0),
			wantCodeCompletionLimitInterval: oneDayInSeconds,
		},
		{
			name:                            "ssc-free",
			user:                            sscFreeUser,
			externalAccount:                 sscFreeUserExternalAccount,
			wantChatLimit:                   graphqlbackend.BigInt(perCommunityUserChatMonthlyLLMRequestLimit),
			wantChatLimitInterval:           oneDayInSeconds * 30,
			wantCodeCompletionLimit:         graphqlbackend.BigInt(0),
			wantCodeCompletionLimitInterval: oneDayInSeconds,
		},
		{
			name:                            "ssc-pro",
			user:                            sscProUser,
			externalAccount:                 sscProUserExternalAccount,
			wantChatLimit:                   graphqlbackend.BigInt(perProUserChatDailyLLMRequestLimit),
			wantChatLimitInterval:           oneDayInSeconds,
			wantCodeCompletionLimit:         graphqlbackend.BigInt(0),
			wantCodeCompletionLimitInterval: oneDayInSeconds,
			pro:                             true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Make user an admin
			err = db.Users().SetIsSiteAdmin(ctx, test.user.ID, true)
			user, err := db.Users().GetByID(ctx, test.user.ID)
			require.NoError(t, err)

			// Create resolver and get user
			_, apiToken, err := db.AccessTokens().Create(ctx, user.ID, []string{authz.ScopeUserAll}, "test", user.ID, time.Time{})
			require.NoError(t, err)
			gatewayToken, err := accesstoken.GenerateDotcomUserGatewayAccessToken(apiToken)
			require.NoError(t, err)

			// Create a request context from the user
			userContext := featureflag.WithFlags(actor.WithActor(ctx, actor.FromActualUser(user)), db.FeatureFlags())
			if test.pro {
				userContext = cody.WithMockSSCClient(userContext, cody.MockSSCClient{
					MockSSCValue: []cody.MockSSCValue{{
						Subscription: &ssc.Subscription{
							Status:             ssc.SubscriptionStatusActive,
							BillingInterval:    ssc.BillingIntervalMonthly,
							CancelAtPeriodEnd:  false,
							CurrentPeriodStart: time.Now().Format(time.RFC3339Nano),
							CurrentPeriodEnd:   time.Now().Format(time.RFC3339Nano),
						},
						SAMSAccountID: test.externalAccount.AccountID,
					}},
					ShouldBeCalled: test.externalAccount != nil,
				})
			} else if test.externalAccount != nil {
				userContext = cody.WithMockSSCClient(userContext, cody.MockSSCClient{
					MockSSCValue: []cody.MockSSCValue{{
						SAMSAccountID: test.externalAccount.AccountID,
					}},
					ShouldBeCalled: true,
				})
			} else {
				userContext = cody.WithMockSSCClient(userContext, cody.MockSSCClient{
					MockSSCValue:   []cody.MockSSCValue{},
					ShouldBeCalled: false,
				})
			}

			r := productsubscription.CodyGatewayDotcomUserResolver{Logger: logtest.Scoped(t), DB: db}
			gatewayUser, err := r.CodyGatewayDotcomUserByToken(userContext, &graphqlbackend.CodyGatewayUsersByAccessTokenArgs{Token: gatewayToken})
			require.NoError(t, err)

			access := gatewayUser.CodyGatewayAccess()
			rateLimit, err := access.ChatCompletionsRateLimit(userContext)
			require.NoError(t, err)
			require.Equal(t, test.wantChatLimit, rateLimit.Limit())
			require.Equal(t, test.wantChatLimitInterval, rateLimit.IntervalSeconds())

			rateLimit, err = access.CodeCompletionsRateLimit(userContext)
			require.NoError(t, err)
			require.Equal(t, test.wantCodeCompletionLimit, rateLimit.Limit())
			require.Equal(t, test.wantCodeCompletionLimitInterval, rateLimit.IntervalSeconds())
		})
	}
}
