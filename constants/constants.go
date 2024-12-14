package constants

// Define constants in this file

//goland:noinspection GoSnakeCaseUsage
const (
	BINARY_NAME = "devd"
)

//goland:noinspection GoSnakeCaseUsage
const (
	ENV_EVM_RPC     = "DEVD_EVM_RPC"
	ENV_COSMOS_REST = "DEVD_COSMOS_REST"
	ENV_TM_RPC      = "DEVD_TM_RPC"

	DEFAULT_EVM_RPC     = "http://localhost:8545"
	DEFAULT_COSMOS_REST = "http://localhost:1317"
	DEFAULT_TM_RPC      = "http://localhost:26657"

	ENV_SECRET_KEY = "DEVD_SECRET_KEY"
)
