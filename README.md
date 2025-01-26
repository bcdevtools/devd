## devd

### Install

```bash
go install -v github.com/bcdevtools/devd/v3/cmd/devd@latest
```

### Query tools

Lazy EVM-RPC setting
```bash
export DEVD_EVM_RPC='https://evm.example.com:8545'
```
_By setting this environment variable, you don't need to pass `--evm-rpc` flag everytime for non-localhost EVM Json-RPC_
___
Lazy TM-RPC setting
```bash
export DEVD_TM_RPC='https://rpc.example.com:26657'
```
_By setting this environment variable, you don't need to pass `--tm-rpc` flag everytime for non-localhost Tendermint RPC_
___
Lazy Rest API setting
```bash
export DEVD_COSMOS_REST='https://cosmos-rest.example.com:1317'
```
_By setting this environment variable, you don't need to pass `--rest` flag everytime for non-localhost Rest API_
___

#### Query account balance

```bash
devd query balance [account addr] [optional ERC20 addr..] [--erc20] [--evm-rpc http://localhost:8545]
# devd q b 0xAccount
# devd q b ethm1account
# devd q b 0xAccount 0xErc20Contract
# devd q b ethm1account 0xErc20Contract1 0xErc20Contract2
# devd q b 0xAccount --erc20 [--rest http://localhost:1317]
```
_`--erc20` flag, if provided, will attempt to fetch user balance of contracts on `x/erc20` module and virtual frontier bank contracts. This request additional Rest-API endpoint provided, or use default 1317._

#### Query block/tx events

```bash
devd query events [height/tx hash] [--filter one] [--filter of_] [--filter these] [--tm-rpc http://localhost:26657]
# devd q events COS...MOS -f sig -f seq_
# devd q events 0x...evm -f txHash
# devd q events 10000
```
_`--filter` flags, if provided, will accept events those contain at least one provided criteria_

#### Query ERC20 token information

```bash
devd query erc20 [ERC20 addr] [optional account] [--evm-rpc http://localhost:8545]
# devd q erc20 0xErc20Contract
# devd q erc20 0xErc20Contract 0xAccount
# devd q erc20 0xErc20Contract ethm1account
```

#### Get EVM transaction information

```bash
devd query eth_getTransactionByHash [0xHash] [--evm-rpc http://localhost:8545]
# devd q tx 0xHash
```

#### Get EVM transaction receipt

```bash
devd query eth_getTransactionReceipt [0xHash] [--evm-rpc http://localhost:8545]
# devd q receipt 0xHash
```

#### Get EVM block by number

```bash
devd query eth_getBlockByNumber [hex or dec block no] [--full] [--evm-rpc http://localhost:8545]
# devd q block 0xF
# devd q block 16 --full
```

#### Trace EVM transaction

```bash
devd query debug_traceTransaction [0xHash] [--tracer callTracer] [--evm-rpc http://localhost:8545]
# devd q trace 0xHash
# devd q trace 0xHash --tracer callTracer
```

### Tx tools

#### Send EVM transaction

```bash
# Transfer native coin
devd tx send [to] [amount] [--raw-tx]
# Transfer ERC-20 token
devd tx send [to] [amount] [--erc20 contract_address]
# Use `--raw-tx` flag to see raw RLP-encoded EVM tx
```

_support short int (2e18, 5bb,...): `devd tx send [to] [1e18/1bb]`_

#### Deploy EVM contract

```bash
# Deploy contract with deployment bytecode
devd tx deploy-contract [deployment bytecode] [--gas 4m] [--gas-prices 20b] [--raw-tx]
# Deploy ERC-20 contract with pre-defined bytecode
devd tx deploy-contract erc20
# Use `--raw-tx` flag to see raw RLP-encoded EVM tx
```

### Convert tools

#### Convert address between different formats

```bash
devd convert address [address] [optional bech32 hrp]
# devd c a 0xAccount ethm
# devd c a ethm1account
# devd c a ethm1account xyz
```
***WARN: DO NOT use this command to convert address across chains with different HD-Path! (eg: Ethermint 60 and Cosmos 118)***

#### Encode string into ABI or decode ABI into string

***Support pipe***
```bash
devd convert abi-string [string or ABI encoded string]
# devd c abi-string 000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000045553444300000000000000000000000000000000000000000000000000000000
# devd c abi-string USDC Token
# echo 'USDC Token' | devd c abi-string
```

#### Convert hexadecimal to decimal and vice versa

***Support pipe***
```bash
devd convert hexadecimal [hexadecimal or decimal]
# devd c hex 0x16a
# echo 0x16a | devd c hex
# devd c hex 362
# echo 362 | devd c hex
```

#### Convert Solidity event/method signature into hashed signature

```bash
devd convert solc-sig [event/method signature]
# devd c solc-sig 'transfer(address,uint256)'
# devd c solc-sig 'function transfer(address recipient, uint256 amount) external returns (bool);'
# devd c solc-sig 'event Transfer(address indexed from, address indexed to, uint256 value);'
```

#### Convert input into lower/upper case

***Support pipe***
```bash
devd convert case [input]
# devd c case AA
# echo AA | devd c case
# > aa
devd convert case [input] --upper
# devd c case aa
# echo aa | devd c case --upper
# > AA
```

#### Encode/Decode base64

***Support pipe***
```bash
devd convert base64 [input]
# devd c base64 123
# echo 123 | devd c base64
devd convert base64 [base64] --decode
# devd c base64 TVRJeg== --decode
# echo TVRJeg== | devd c base64 --decode
```

#### Convert raw balance into display balance and vice versa

```bash
devd convert display-balance [raw balance] [exponent]
# devd c dbal 10011100 6
# > 10.0111
# Support short int:
#  devd c dbal 20bb 18
#  > 20.0
```

```bash
devd convert raw-balance [display balance] [exponent] [--decimals-point , or .]
# devd c rbal 10.0111 6
# > 10011100
# devd c rbal 10,0111 6 -d ,
# > 10011100
```

### Hashing tools

***Support pipe***
```bash
devd hash md5 [input]
# devd hash md5 123
# cat file.txt | devd hash md5
devd hash keccak256 [input]
# devd hash keccak256 123
# cat file.txt | devd hash keccak256
devd hash keccak512 [input]
# devd hash keccak512 123
# cat file.txt | devd hash keccak512
```

### Check tools

#### Listing ports in use and check port holding by process

```bash
# listing
devd check port

# check specific port
devd check port [port]
```

### Debug tools

#### Decode raw RLP-encoded EVM tx into tx object

```bash
devd debug raw-tx [raw RLP-encoded EVM tx hex]
# devd debug raw-tx 0x02f8af82271c3b83112a8883aba95082be209480b5a32e4f032b2a058b4f29ec95eefeeb87adcd80b844a9059cbb000000000000000000000000bfcfe6d5ad56aa831313856949e98656d46f9248000000000000000000000000000000000000000000000000002386f26fc10000c001a088907374a796ed70a5a2bdc51b50010b68dcc4d2ed12d94abc607bb0a90271b6a0167d3e031b70ec511b67416d9ad8334caee7013d95ff8275721b23798c5c3602
```
_to view inner tx information, including sender address_

#### Compute EVM transaction intrinsic gas

```bash
devd debug intrinsic-gas [0xCallData]
# devd d intrinsic-gas 0xCallData
```
_Assumption: no access list, not contract creation, Homestead, EIP-2028 (Istanbul). If contract creation, plus 32,000 into the output._

### Notes:

- Output messages are printed via stdout, while messages with prefixes `INF:` `WARN:` and `ERR:` are printed via stderr. So for integration with other tools, to omit stderr, forward stdout only.
  > Eg: `devd c a cosmos1... 1> /tmp/output.txt`
- When passing arguments into command via both argument and pipe, the argument will be used.
  > Eg: `echo 123 | devd c hex 456` will convert `456` to hexadecimal, not `123`.
- For commands those marked `support short int`, you can pass number with format like:
  - `2e18` = 2 x 10^18
  - `2k` = 2,000
  - `2.7m` = 2,700,000
  - `3.08b` = 3,080,000,000
  - `4kb` = 4,000,000,000,000
  - `5.555mb` = 5,555,000,000,000,000
  - `6bb` = 6,000,000,000,000,000,000 = 6e18
- Some queries will try to decode some fields in response data into human-readable format and inject back into the response data with `_` prefix like EVM tx, receipt, block, trace.