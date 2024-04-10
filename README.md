## devd

### Install

```bash
go install -v github.com/bcdevtools/devd/cmd/devd@latest
```

### Query tools

Lazy RPC setting
```bash
export DEVD_EVM_RPC='https://api.securerpc.com/v1'
```
_By setting this environment variable, you don't need to pass --rpc flag everytime for non-localhost EVM Json-RPC_

#### Query account balance

```bash
devd query balance [account_address] [optional_erc20_contract_addresses...] [--rpc http://localhost:8545]
# devd q b 0xAccount
# devd q b ethm1account
# devd q b 0xAccount 0xErc20Contract
# devd q b ethm1account 0xErc20Contract 0xErc20Contract
```

#### Query ERC20 token information

```bash
devd query erc20 [erc20_contract_address] [optional_account_address] [--rpc http://localhost:8545]
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

### Convert tools

#### Convert address between different formats

```bash
devd convert address [address] [optional_bech32]
# devd c a 0xAccount ethm
# devd c a ethm1account
# devd c a ethm1account xyz
```
***WARN: DO NOT use this command to convert address across chains with different HD-Path! (eg: Ethermint 60 and Cosmos 118)***

#### Encode string into ABI or decode ABI into string

```bash
devd convert abi_string [string or ABI encoded string]
# devd c abi_string 000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000045553444300000000000000000000000000000000000000000000000000000000
# devd c abi_string USDC Token
```

#### Convert hexadecimal to decimal and vice versa

```bash
devd convert hex_2_dec [hexadecimal or decimal]
# devd c h2d 0x16a
# devd c h2d 362
# devd c h2d 16a
```

#### Convert Solidity event/method signature into hashed signature

```bash
devd convert solc_sig [event/method signature]
# devd c solc_sig 'transfer(address,uint256)'
# devd c solc_sig 'function transfer(address recipient, uint256 amount) external returns (bool);'
# devd c solc_sig 'event Transfer(address indexed from, address indexed to, uint256 value);'
```

### Debug tools

#### Compute EVM transaction intrinsic gas

```bash
devd debug intrinsic_gas [0xdata]
# devd d intrinsic_gas 0xdata
```
