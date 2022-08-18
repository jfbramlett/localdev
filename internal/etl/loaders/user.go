package loaders

import (
	"context"
	"github.com/splice/platform/localdev/v2/localdev/internal/etl/datbase"

	splicelogger "github.com/splice/platform/infra/libs/golang/logger"
)

func NewUserLoader(srcDB datbase.Database, tgtDB datbase.Database) Loader {
	return &UserLoader{
		srcDB: srcDB,
		tgtDB: tgtDB,
	}
}

type UserLoader struct {
	srcDB datbase.Database
	tgtDB datbase.Database
}

func (c *UserLoader) Load(ctx context.Context) error {
	logger, _ := splicelogger.FromContext(ctx)
	logger.Info("starting user load")
	defer logger.Info("completed user load")

	logger.Info("deleting user tables")
	if err := DeleteTables(ctx, c.tgtDB, "frontend_releases", "users", "user_subscriptions", "users_tos",
		"user_credits", "user_emails", "user_frontend_releases", "payment_methods", "purchases", "credit_awards",
		"sounds_user_states", "entitlements"); err != nil {
		return err
	}

	logger.Info("loading frontend_releases")
	stmt := `
INSERT INTO frontend_releases (sha, created_at, updated_at)
VALUES
	('5091f700b5b6a2d60ae202db38ab2f13d06b7856', '2021-01-13 15:02:47', '2021-01-13 15:02:47');
`
	if err := c.tgtDB.Exec(ctx, stmt); err != nil {
		return err
	}

	logger.Info("loading users")
	stmt = `
INSERT INTO users (created_at, updated_at, email, encrypted_password, reset_password_token, reset_password_sent_at, sign_in_count, current_sign_in_at, last_sign_in_at, current_sign_in_ip, last_sign_in_ip, confirmation_token, confirmed_at, confirmation_sent_at, unconfirmed_email, authentication_token, username, name, avatar, test_account, bio, location, signup_type, tos, prior_username, subscriptions, customer_id, tfa, tfa_enabled_at, tfa_secret, phone_verified, phone_number, sounds, braintree_customer_id, role, recurly_uuid, uuid, taxable_at)
VALUES
	('2020-08-04 00:01:11', '2021-01-13 15:22:32', 'splice-user@splice.com', '$2a$10$o4nY3eGhmAvnIHR4fneuYuDEksX2ZRL/0deJpEHBF6/s/d3jJ1WBq', NULL, NULL, 6, '2021-01-13 15:22:50', '2021-01-13 15:22:32', '74.64.208.135', '74.64.208.135', NULL, NULL, NULL, NULL, 'MrzIOf1etbmsWDiHhi5sQC4V', 'spliceuser', 'Splice User', NULL, 0, NULL, NULL, 'home', 1, NULL, 9223372036854775807, NULL, 0, NULL, NULL, 0, NULL, 1, NULL, 0, '56680920-086e-0780-776e-da53b361c3c1e629271', 'e1639c13-0427-9680-4a0b-0041b25b05cdfe9b506', '2020-11-25 18:37:37.140');
`
	if err := c.tgtDB.Exec(ctx, stmt); err != nil {
		return err
	}

	stmt = `
INSERT INTO users (created_at, updated_at, email, encrypted_password, reset_password_token, reset_password_sent_at, sign_in_count, current_sign_in_at, last_sign_in_at, current_sign_in_ip, last_sign_in_ip, confirmation_token, confirmed_at, confirmation_sent_at, unconfirmed_email, authentication_token, username, name, avatar, test_account, bio, location, signup_type, tos, prior_username, subscriptions, customer_id, tfa, tfa_enabled_at, tfa_secret, phone_verified, phone_number, sounds, braintree_customer_id, role, recurly_uuid, uuid, taxable_at)
VALUES
	('2020-08-04 00:01:11', '2021-01-13 15:22:32', 'splice-admin@splice.com', '$2a$10$o4nY3eGhmAvnIHR4fneuYuDEksX2ZRL/0deJpEHBF6/s/d3jJ1WBq', NULL, NULL, 6, '2021-01-13 15:22:50', '2021-01-13 15:22:32', '74.64.208.135', '74.64.208.135', NULL, NULL, NULL, NULL, 'cNfhiqLSr9HjzYJkx_h1', 'spliceadmin', 'Splice Admin', NULL, 0, NULL, NULL, 'home', 1, NULL, 9223372036854775807, NULL, 0, NULL, NULL, 0, NULL, 1, NULL, 15, '96680920-086e-0780-776e-da53b361c3c1e6292ef', 'f2639c13-0427-9680-4a0b-0041b25b05cdfe9b556', '2020-11-25 18:37:37.140');
`
	if err := c.tgtDB.Exec(ctx, stmt); err != nil {
		return err
	}

	logger.Info("loading user_subscriptions")
	stmt = `
INSERT INTO user_subscriptions (uuid, plan_uuid, user_id, created_at, updated_at, expires_at, provider_id, renews_at, downgrades_at, activates_at, upgraded_at, cancellation_requested_at, payment_method_id, provider_type, parent_subscription_id, gift_card_redemption_uuid, trial_expires_at, current, internal_trial_period, payment_provider_switched_at, pause_at, paused_at, resume_at, resumed_at, retry_resumption_at, resumption_retries_count, hold_scheduled_at, payment_method_switched_at, trial_converted_at, resumption_payment_method_uuid, resumption_payment_method_type, resumption_plan_uuid, cancelled_at, cancellation_completed_at, manually_resumed_at, limboed_at, is_paid)
VALUES
	('b4c61613-0b4a-6a80-d02e-4f6ccc586a1f0878761', (select uuid from plans where code='sounds-lite' and provider_id='sounds-lite-0002'), (select id from users where username='spliceuser'), '2020-11-25 18:37:52', '2020-12-25 18:40:36', '2021-01-25 18:39:18', '576f7dcfd2c5aedb479c8946e8bc3cae', '2021-01-25 18:39:18', NULL, '2020-11-25 18:37:51', NULL, NULL, 1185691, 'recurly', NULL, NULL, NULL, 1, 0, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, '2020-11-25 18:38:18', NULL, NULL, NULL, NULL, NULL, NULL, NULL, 1);
`
	if err := c.tgtDB.Exec(ctx, stmt); err != nil {
		return err
	}

	logger.Info("loading user_tos")
	stmt = `
INSERT INTO users_tos (user_id, accepted_at)
VALUES
	((select id from users where username='spliceuser'), '2020-08-04 00:01:11');
`
	if err := c.tgtDB.Exec(ctx, stmt); err != nil {
		return err
	}

	logger.Info("loading user_credits")
	stmt = `
INSERT INTO user_credits (user_id, credit_count, created_at, updated_at)
VALUES
	((select id from users where username='spliceuser'), 191, '2020-11-25 18:38:18', '2020-12-25 18:40:35');
`
	if err := c.tgtDB.Exec(ctx, stmt); err != nil {
		return err
	}

	logger.Info("loading user_emails")
	stmt = `
INSERT INTO user_emails (user_id, email, status, created_at, updated_at)
VALUES
	((select id from users where username='spliceuser'), 'john.bramlett@splice.com', 'valid', '2021-01-13 15:22:49', NULL);
`
	if err := c.tgtDB.Exec(ctx, stmt); err != nil {
		return err
	}

	logger.Info("loading user_frontend_releases")
	stmt = `
INSERT INTO user_frontend_releases (user_id, frontend_release_id, created_at, updated_at)
VALUES
	((select id from users where username='spliceuser'), (select id from frontend_releases where sha='5091f700b5b6a2d60ae202db38ab2f13d06b7856'), '2020-08-04 09:28:43', '2021-01-13 15:32:17');
`
	if err := c.tgtDB.Exec(ctx, stmt); err != nil {
		return err
	}

	logger.Info("loading payment_methods")
	stmt = `
INSERT INTO payment_methods (customer_id, token, type, user_id, uuid, created_at, updated_at, last4, paypal_email, exp_month, exp_year, provider_type, one_time_purchase, brand, postal_code, address1, address2, city, state, country, first_name, last_name)
VALUES
	('56680920-086e-0780-776e-da53b361c3c1e629271', '', 'credit-card', (select id from users where username='spliceuser'), '87368849-0eb6-c680-a7a0-f520b56d0b3e83c924e', '2020-11-25 18:37:41', NULL, '1111', NULL, 12, 2025, 'recurly', 0, 'Visa', '10023', '19 W 69th St', NULL, 'New York', 'NY', 'US', 'John', 'Bramlett');
`
	if err := c.tgtDB.Exec(ctx, stmt); err != nil {
		return err
	}

	logger.Info("loading purchases")
	stmt = `
INSERT INTO purchases (uuid, user_id, credit_card_uuid, provider_charge_id, amount, deliverable_uuid, deliverable_type, created_at, product, description, provider_type, payment_method_uuid, subscription_uuid, subtotal, tax, total)
VALUES
	('f804ed19-08c1-fd80-e5a0-cb6c9341264a3118bea', (select id from users where username='spliceuser'), NULL, '1432129', 799, '', 'subscription', '2020-11-25 19:04:00', 'Sounds', 'Sounds Subscription', 'recurly', '87368849-0eb6-c680-a7a0-f520b56d0b3e83c924e', 'b4c61613-0b4a-6a80-d02e-4f6ccc586a1f0878761', 799, 0, 799),
	('b13dbf23-0bdd-df80-88d1-e76f3a8006ee97a20c2', (select id from users where username='spliceuser'), NULL, '1487630', 799, '', 'subscription', '2020-12-25 18:40:36', 'Sounds', 'Sounds Subscription', 'recurly', '87368849-0eb6-c680-a7a0-f520b56d0b3e83c924e', 'b4c61613-0b4a-6a80-d02e-4f6ccc586a1f0878761', 799, 0, 799);
`
	if err := c.tgtDB.Exec(ctx, stmt); err != nil {
		return err
	}

	logger.Info("loading credit_awards")
	stmt = `
INSERT INTO credit_awards (user_id, user_subscription_uuid, reason, invoice_id, created_at, credits, plan_uuid)
VALUES
	((select id from users where username='spliceuser'), 'b4c61613-0b4a-6a80-d02e-4f6ccc586a1f0878761', 'recurly sounds purchase: manual conversion', NULL, '2020-11-25 18:38:19', 100, '9e2c7430-054f-3480-6d77-75ac069e0f345483d9d'),
	((select id from users where username='spliceuser'), 'b4c61613-0b4a-6a80-d02e-4f6ccc586a1f0878761', 'user credit refill on charge', '1487630', '2020-12-25 18:40:36', 100, '9e2c7430-054f-3480-6d77-75ac069e0f345483d9d');
`
	if err := c.tgtDB.Exec(ctx, stmt); err != nil {
		return err
	}

	logger.Info("loading sounds_user_states")
	stmt = `
INSERT INTO sounds_user_states (user_id, created_at, updated_at, state_changed_at, state, user_subscription_uuid, previous_state)
VALUES
	((select id from users where username='spliceuser'), '2020-11-25 18:37:52', NULL, '2020-11-25 18:37:52', 'trialing', 'b4c61613-0b4a-6a80-d02e-4f6ccc586a1f0878761', ''),
	((select id from users where username='spliceuser'), '2020-11-25 18:38:19', NULL, '2020-11-25 18:38:19', 'subscribed', 'b4c61613-0b4a-6a80-d02e-4f6ccc586a1f0878761', 'trialing');
`
	if err := c.tgtDB.Exec(ctx, stmt); err != nil {
		return err
	}

	logger.Info("loading entitlements")
	stmt = `
INSERT INTO entitlements (grantee_type, grantee_uuid, grantor, resource, created_at, updated_at, expires_at, deleted_at)
VALUES
	('user', (select uuid from users where username='spliceuser'), 'subscription-manager', 'catalog/view/standard', '2020-11-25 18:37:51.703', '2020-12-25 18:40:35.938', '2021-01-26 00:39:18.000', '9999-01-01 00:00:00.000'),
	('user', (select uuid from users where username='spliceuser'), 'subscription-manager', 'dummy/trial/lite', '2020-11-25 18:37:51.703', '2020-11-25 18:38:18.664', '2020-11-26 18:37:51.536', '2020-11-25 18:38:18.000'),
	('user', (select uuid from users where username='spliceuser'), 'subscription-manager', 'quota/credits/100', '2020-11-25 18:38:18.663', '2020-12-25 18:40:35.938', '2021-01-26 00:39:18.000', '9999-01-01 00:00:00.000'),
	('user', (select uuid from users where username='spliceuser'), 'subscription-manager', 'dummy/subscribed/lite', '2020-11-25 18:38:18.663', '2020-12-25 18:40:35.938', '2021-01-26 00:39:18.000', '9999-01-01 00:00:00.000'),
	('user', (select uuid from users where username='spliceuser'), 'slackbot', 'catalog/view/standard', '2020-12-16 19:05:41.235', '2020-12-16 19:10:24.643', '2120-11-22 19:10:24.640', '9999-01-01 00:00:00.000'),
	('user', (select uuid from users where username='spliceuser'), 'slackbot', 'catalog/view/skills', '2020-12-16 19:05:41.235', '2020-12-16 19:10:24.643', '2120-11-22 19:10:24.640', '9999-01-01 00:00:00.000');
`
	if err := c.tgtDB.Exec(ctx, stmt); err != nil {
		return err
	}

	stmt = `
INSERT INTO entitlements (grantee_type, grantee_uuid, grantor, resource, created_at, updated_at, expires_at, deleted_at)
VALUES
	('user', (select uuid from users where username='spliceadmin'), 'slackbot', 'catalog/view/admin', '2020-12-16 19:05:41.235', '2020-12-16 19:10:24.643', '2120-11-22 19:10:24.640', '9999-01-01 00:00:00.000');
`
	if err := c.tgtDB.Exec(ctx, stmt); err != nil {
		return err
	}

	return nil
}
