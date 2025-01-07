package cody

import (
	"context"
	"sync"
	"time"

	"github.com/khulnasoft/khulnasoft/cmd/frontend/internal/ssc"
	"github.com/khulnasoft/khulnasoft/internal/database"
	"github.com/khulnasoft/khulnasoft/internal/types"
	"github.com/khulnasoft/khulnasoft/lib/errors"
)

type UserSubscriptionPlan string

const (
	UserSubscriptionPlanFree UserSubscriptionPlan = "FREE"
	UserSubscriptionPlanPro  UserSubscriptionPlan = "PRO"
)

type UserSubscription struct {
	// Status is the current status of the subscription. "pending" means the user has no Cody Pro subscription.
	// "pending" subscription will be removed post Feb 15, 2024. It is required to support users who have opted for
	// Cody Pro Trial on dotcom, but have not entered payment information on SSC.
	// (So they don't have an SSC backed subscription, but we need to act like they do, until 2/15.)
	Status ssc.SubscriptionStatus
	// Plan is the plan the user is subscribed to.
	Plan UserSubscriptionPlan
	// ApplyProRateLimits indicates the user should be given higher rate limits
	// for Cody and related functionality. (Use this value instead of checking
	// the subscription status for simplicity.)
	ApplyProRateLimits bool
	// CurrentPeriodStartAt is the start date of the current billing cycle.
	CurrentPeriodStartAt time.Time
	// CurrentPeriodEndAt is the end date of the current billing cycle.
	// IMPORTANT: This may be IN THE PAST. e.g. if the subscription was
	// canceled, this will be when the subscription ended.
	CurrentPeriodEndAt time.Time
	// CancelAtPeriodEnd flags whether or not a subscription will automatically cancel at the end
	// of the current billing cycle, or if it will renew.
	CancelAtPeriodEnd bool
}

// shouldHaveCodyPro returns whether or not the user should have access to Cody Pro functionality
// based on their subscription status.
func shouldHaveCodyPro(status ssc.SubscriptionStatus) bool {
	switch status {
	case ssc.SubscriptionStatusActive, ssc.SubscriptionStatusPastDue, ssc.SubscriptionStatusTrialing:
		// Active is the regular state for a valid subscription.
		// PastDue is when there is some form of payment problem, but is within the grace period before
		//     the subscription gets canceled.
		// Trialing is when the user is on a free trial of Cody Pro.
		return true
	default:
		return false
	}
}

// consolidateSubscriptionDetails handles the raw subscription data from SSC.
func consolidateSubscriptionDetails(ctx context.Context, user types.User, subscription *ssc.Subscription) (*UserSubscription, error) {
	if subscription != nil && shouldHaveCodyPro(subscription.Status) {
		currentPeriodStart, err := time.Parse(time.RFC3339, subscription.CurrentPeriodStart)
		if err != nil {
			return nil, err
		}

		currentPeriodEnd, err := time.Parse(time.RFC3339, subscription.CurrentPeriodEnd)
		if err != nil {
			return nil, err
		}

		return &UserSubscription{
			Status:               subscription.Status,
			Plan:                 UserSubscriptionPlanPro,
			ApplyProRateLimits:   true,
			CurrentPeriodStartAt: currentPeriodStart,
			CurrentPeriodEndAt:   currentPeriodEnd,
			CancelAtPeriodEnd:    subscription.CancelAtPeriodEnd,
		}, nil
	}

	// If the user doesn't have a subscription in the SSC backend or it is cancelled, then we need
	// find current period date range using the user.createdAt.
	currentPeriodStartAt, currentPeriodEndAt, err := freeUserCurrentPeriodDateRange(ctx, user, subscription)
	if err != nil {
		return nil, err
	}

	return &UserSubscription{
		Status:               ssc.SubscriptionStatusActive,
		Plan:                 UserSubscriptionPlanFree,
		ApplyProRateLimits:   false,
		CurrentPeriodStartAt: *currentPeriodStartAt,
		CurrentPeriodEndAt:   *currentPeriodEndAt,
		CancelAtPeriodEnd:    false,
	}, nil
}

// getSAMSAccountIDsForUser returns the user's SAMS account ID if available.
//
// If the user has not associated a SAMS identity with their dotcom user account,
// will return ("", nil). After we migrate all dotcom user accounts to SAMS, that
// should no longer be possible.
func getSAMSAccountIDsForUser(ctx context.Context, db database.DB, dotcomUserID int32) ([]string, error) {
	oidcAccounts, err := db.UserExternalAccounts().List(ctx, database.ExternalAccountsListOptions{
		UserID:      dotcomUserID,
		ServiceType: "openidconnect",
		ServiceID:   ssc.GetSAMSServiceID(),
	})
	if err != nil {
		return []string{}, errors.Wrap(err, "listing external accounts")
	}

	var ids []string

	for _, account := range oidcAccounts {
		ids = append(ids, account.AccountID)
	}

	return ids, nil
}

// SubscriptionForUser returns the user's Cody subscription details.
func SubscriptionForUser(ctx context.Context, db database.DB, user types.User) (*UserSubscription, error) {
	samsAccountIDs, err := getSAMSAccountIDsForUser(ctx, db, user.ID)
	if err != nil {
		return nil, errors.Wrap(err, "fetching user's SAMS account ID")
	}

	if len(samsAccountIDs) == 0 {
		return consolidateSubscriptionDetails(ctx, user, nil)
	}

	sscClient, err := SSCClientFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	// NOTE(naman): As part of #inc-284-plg-users-paying-for-and-being-billed-for-pro-without-being-upgrade
	// it is noted that a user can have multiple SAMS accounts associated with their dotcom account.
	// Originally we only fetchted subscription data from the first SAMS account associated with the user.
	// But we are now fetching subscription data from all SAMS accounts associated with the user, and using
	// the one with active subscription. This is only a TEMPORARY solution to mititgate the incident.
	// MUST be removed/fixed once we have a proper solution in place as part of #60912.
	var subscriptions []*ssc.Subscription
	for _, samsAccountID := range samsAccountIDs {
		// While developing the SSC backend, we only fetch subscription data for users based on a flag.
		var subscription *ssc.Subscription
		if samsAccountID != "" {
			subscription, err = sscClient.FetchSubscriptionBySAMSAccountID(ctx, samsAccountID)
			if err != nil {
				return nil, errors.Wrap(err, "fetching subscription from SSC")
			}
		}

		if subscription != nil {
			subscriptions = append(subscriptions, subscription)

			// Pick the first one we find in a good state, enabling the user to access Cody Pro features.
			if shouldHaveCodyPro(subscription.Status) {
				return consolidateSubscriptionDetails(ctx, user, subscription)
			}
		}
	}

	if len(subscriptions) == 0 {
		return consolidateSubscriptionDetails(ctx, user, nil)
	}
	// If we didn't find an "working" subscription, then we just return the details
	// of the first one we found. (For this to be stable, we require the list of external
	// accounts to be returned from the database in a sorted order, which it is.)
	return consolidateSubscriptionDetails(ctx, user, subscriptions[0])
}

// SubscriptionForSAMSAccountID returns the user's Cody subscription details associated with the given SAMS account ID.
func SubscriptionForSAMSAccountID(ctx context.Context, db database.DB, user types.User, samsAccountID string) (*UserSubscription, error) {
	if samsAccountID == "" {
		return consolidateSubscriptionDetails(ctx, user, nil)
	}

	var subscription *ssc.Subscription
	sscClient, err := SSCClientFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	subscription, err = sscClient.FetchSubscriptionBySAMSAccountID(ctx, samsAccountID)
	if err != nil {
		return nil, errors.Wrap(err, "fetching subscription from SSC")
	}

	return consolidateSubscriptionDetails(ctx, user, subscription)
}

// getSSCClient returns a self-service Cody API client. We only do this once so that the stateless client
// can persist in memory longer, so we can benefit from the underlying HTTP client only needing to reissue
// SAMS access tokens when needed, rather than minting a new token for every request.
//
// BUG: If the SAMS configuration is added or changed during the lifetime of the this process, the returned
// client will be invalid. (As it would be using the original SAMS client configuration data.) The process
// will need to be restarted to correct this situation.
var getSSCClient = sync.OnceValues[ssc.Client, error](func() (ssc.Client, error) {
	return ssc.NewClient()
})

type MockSSCValue struct {
	Subscription  *ssc.Subscription
	SAMSAccountID string
}

type MockSSCClient struct {
	MockSSCValue   []MockSSCValue
	ShouldBeCalled bool
}

func (m MockSSCClient) FetchSubscriptionBySAMSAccountID(
	ctx context.Context, samsAccountID string) (*ssc.Subscription, error) {
	if !m.ShouldBeCalled {
		return nil, errors.Errorf("FetchSubscriptionBySAMSAccountID should not have be called")
	}

	for _, v := range m.MockSSCValue {
		if v.SAMSAccountID == samsAccountID {
			return v.Subscription, nil
		}
	}

	return nil, errors.Errorf("FetchSubscriptionBySAMSAccountID should not have be called with the given samsAccountID: %s", samsAccountID)
}

type sscClientCtxKey int

const mockSSCClientCtxKey sscClientCtxKey = iota

func SSCClientFromCtx(ctx context.Context) (ssc.Client, error) {
	c, ok := ctx.Value(mockSSCClientCtxKey).(ssc.Client)
	if !ok || c == nil {
		return getSSCClient()
	}
	return c, nil
}

func WithMockSSCClient(ctx context.Context, c ssc.Client) context.Context {
	return context.WithValue(ctx, mockSSCClientCtxKey, c)
}
