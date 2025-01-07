## devd

### Install

```bash
go install -v github.com/bcdevtools/devd/v2/cmd/devd@latest
```

### Query tools

Lazy RPC setting
```bash
export DEVD_EVM_RPC='https://evm.example.com:8545'
```
_By setting this environment variable, you don't need to pass `--rpc` flag everytime for non-localhost EVM Json-RPC_
___
Lazy Rest API setting
```bash
export DEVD_COSMOS_REST='https://ethermint-rest.example.com:1317'
```
_By setting this environment variable, you don't need to pass `--rest` flag everytime for non-localhost Rest API_

#### Query account balance

```bash
devd query balance [account addr] [optional ERC20 addr..] [--erc20] [--rpc http://localhost:8545]
# devd q b 0xAccount
# devd q b ethm1account
# devd q b 0xAccount 0xErc20Contract
# devd q b ethm1account 0xErc20Contract1 0xErc20Contract2
# devd q b 0xAccount --erc20 [--rest http://localhost:1317]
```
_`--erc20` flag, if provided, will attempt to fetch user balance of contracts on `x/erc20` module and virtual frontier bank contracts. This request additional Rest-API endpoint provided, or use default 1317._

#### Query block/tx events
```bash
devd query events [height/tx hash] [--filter one] [--filter two] [--tm-rpc http://localhost:26657]
# devd q events COS...MOS
# devd q events 0x...evm
# devd q events 10000
```

#### Query ERC20 token information

```bash
devd query erc20 [ERC20 addr] [optional account] [--rpc http://localhost:8545]
# devd q erc20 0xErc20Contract
# devd q erc20 0xErc20Contract 0xAccount
# devd q erc20 0xErc20Contract ethm1account
```

#### Get EVM transaction information

```bash
devd query eth_getTransactionByHash [0xHash] [--rpc http://localhost:8545]
# devd q tx 0xHash
```

#### Get EVM transaction receipt

```bash
devd query eth_getTransactionReceipt [0xHash] [--rpc http://localhost:8545]
# devd q receipt 0xHash
```

#### Get EVM block by number

```bash
devd query eth_getBlockByNumber [hex or dec block no] [--full] [--rpc http://localhost:8545]
# devd q block 0xF
# devd q block 16 --full
```

#### Trace EVM transaction

```bash
devd query debug_traceTransaction [0xHash] [--tracer callTracer] [--rpc http://localhost:8545]
# devd q trace 0xHash
# devd q trace 0xHash --tracer callTracer
```

### Tx tools

#### Send EVM transaction

```bash
# Transfer native coin
devd tx send [to] [amount]
# Transfer ERC-20 token
devd tx send [to] [amount] [--erc20 contract_address]
```

#### Deploy EVM contract

```bash
# Deploy contract with deployment bytecode
devd tx deploy-contract [deployment bytecode]
# Deploy ERC-20 contract with pre-defined bytecode
devd tx deploy-contract erc20
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
devd convert abi_string [string or ABI encoded string]
# devd c abi_string 000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000045553444300000000000000000000000000000000000000000000000000000000
# devd c abi_string USDC Token
# echo 'USDC Token' | devd c abi_string
```

#### Convert hexadecimal to decimal and vice versa

***Support pipe***
```bash
devd convert hex_2_dec [hexadecimal]
# devd c h2d 0x16a
# devd c h2d 16a
# echo 16a | devd c h2d
devd convert dec_2_hex [decimal]
# devd c d2h 170
# echo 170 | devd c d2h
```

#### Convert Solidity event/method signature into hashed signature

```bash
devd convert solc_sig [event/method signature]
# devd c solc_sig 'transfer(address,uint256)'
# devd c solc_sig 'function transfer(address recipient, uint256 amount) external returns (bool);'
# devd c solc_sig 'event Transfer(address indexed from, address indexed to, uint256 value);'
```

#### Convert input into upper/lower case

***Support pipe***
```bash
devd convert to_lower_case [input]
# devd c lowercase AA
# echo AA | devd c lowercase
devd convert to_upper_case [input]
# devd c uppercase aa
# echo aa | devd c uppercase
```

#### Encode/Decode base64

***Support pipe***
```bash
devd convert encode_base64 [input]
# devd c base64 123
# echo 123 | devd c base64
devd convert decode_base64 [base64]
# devd c decode_base64 TVRJeg==
# echo TVRJeg== | devd c decode_base64
```

#### Convert raw balance into display balance and vice versa

```bash
devd convert display_balance [raw balance] [exponent]
# devd c dbal 10011100 6
# > 10.0111
```

```bash
devd convert raw_balance [display balance] [exponent] [--decimals-point , or .]
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

#### Compute EVM transaction intrinsic gas

```bash
devd debug intrinsic_gas [0xCallData]
# devd d intrinsic_gas 0xCallData
```
_Assumption: no access list, not contract creation, Homestead, EIP-2028 (Istanbul). If contract creation, plus 32,000 into the output._

### Notes:

- Output messages are printed via stdout, while messages with prefixes `INF:` `WARN:` and `ERR:` are printed via stderr. So for integration with other tools, to omit stderr, forward stdout only.
  > Eg: `devd c a cosmos1... 1> /tmp/output.txt`
- When passing arguments into command via both argument and pipe, the argument will be used.
  > Eg: `echo 123 | devd c d2h 456` will convert `456` to hexadecimal, not `123`.
