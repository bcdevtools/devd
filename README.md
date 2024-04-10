## devd

### Install

```bash
go install -v github.com/bcdevtools/devd/cmd/devd@latest
```

### Tools

#### Query ERC20 token information

```bash
devd query erc20 [contract_address] [optional_account_address] [--rpc http://localhost:8545]
# devd q erc20 0x12..89
# devd q erc20 0x12..89 0x34..FF
# devd q erc20 0x12..89 ethm1...zz
```

#### Convert address between different formats

```bash
devd convert address [address] [optional_bech32]
# devd c a 0x12..89 ethm
# devd c a ethm1...zz
# devd c a ethm1...zz xyz
```

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

#### Get EVM transaction information

```bash
devd query eth_getTransactionByHash [0xhash] [--rpc http://localhost:8545]
# devd q tx 0xAA..FF
```

#### Get EVM transaction receipt

```bash
devd query eth_getTransactionReceipt [0xhash] [--rpc http://localhost:8545]
# devd q receipt 0xAA..FF
```