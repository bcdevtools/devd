## devd

### Install

```bash
go install -v github.com/bcdevtools/devd/cmd/devd@latest
```

### Tools

#### Query ERC20 token information

```bash
devd query erc20 [contract_address] [optional_account_address] [--host ...]
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
dev convert abi_string [string or ABI encoded string]
# dev c abi_string 000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000045553444300000000000000000000000000000000000000000000000000000000
# dev c abi_string USDC Token
```

#### Convert hexadecimal to decimal and vice versa

```bash
dev convert hex_2_dec [hexadecimal or decimal]
# dev c h2d 0x16a
# dev c h2d 362
# dev c h2d 16a
```
