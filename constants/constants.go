package constants

// Define constants in this file

//goland:noinspection GoSnakeCaseUsage
const (
	BINARY_NAME = "devd"
)

//goland:noinspection GoSnakeCaseUsage
const (
	REQUIRE_NETRC_FILE_PERMISSION = 0o600

	REQUIRE_SSH_DIR_PERMISSION                  = 0o750
	REQUIRE_SSH_AUTHORIZED_KEYS_FILE_PERMISSION = 0o644
	REQUIRE_SSH_PRIVATE_KEY_FILE_PERMISSION     = 0o400
	ACCEPTABLE_SSH_PRIVATE_KEY_FILE_PERMISSION  = 0o600

	REQUIRE_SECRET_DIR_PERMISSION  = 0o700
	REQUIRE_SECRET_FILE_PERMISSION = 0o400

	RECOMMEND_PUBLIC_KEY_FILE_PERMISSION  = 0o644
	RECOMMEND_KNOWN_HOSTS_FILE_PERMISSION = 0o600
	RECOMMEND_SSH_CONFIG_FILE_PERMISSION  = 0o600
)

//goland:noinspection GoSnakeCaseUsage
const (
	ENV_RSYNC_PASSWORD = "RSYNC_PASSWORD"
	ENV_SSHPASS        = "SSHPASS"
)

//goland:noinspection GoSnakeCaseUsage
const (
	PREDEFINED_ALIAS_FILE_NAME = ".devd_alias"
)

//goland:noinspection GoSnakeCaseUsage
const (
	FLAG_REQUIRE_WORKING_USERNAME = "require-username"

	FLAG_USE_WORKING_USERNAME = "iam"
)

//goland:noinspection GoSnakeCaseUsage
const (
	SECRETS_DIR_NAME = ".secrets"
)
