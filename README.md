## devd
#### Install
```bash
go install -v github.com/bcdevtools/devd/cmd/devd@latest
```

### Sync files between servers:
```bash
devd f rsync --help
```

```bash
RSYNC_PASSWORD=1234567 devd files rsync /var/log/nginx/access.log backup@192.168.0.2:/mnt/md0/backup/nginx-logs \
  --local-to-remote
```

```bash
devd files rsync /var/log/nginx/access.log backup@192.168.0.2:/mnt/md0/backup/nginx-logs \
  --local-to-remote --password-file ~/password.txt
```

```bash
SSHPASS=1234567 devd files rsync /var/log/nginx/access.log backup-server:/mnt/md0/backup/nginx-logs \
  --local-to-remote --passphrase
```

Notes:
- This use rsync
- When either source or destination is remote machine:
  - Either environment variable RSYNC_PASSWORD or ENV_SSHPASS or flag --password-file is required (priority flag)
  - Environment variables RSYNC_PASSWORD and ENV_SSHPASS are treated similar thus either needed. If both provided, must be identical
  - You must connect to that remote server at least one time before to perform host key verification (one time action) because the transfer will be performed via ssh.

### Secure delete files and directory:
```bash
devd files rm [file/dir] [--delete]
```

### Command aliases
```bash
devd a --help

# Listing aliases
devd a

# Invoke alias
devd a <alias>
```

_Defined your own aliases by create a TSV `~/.devd_alias`_

### Download file
```bash
devd dl --help
```

```bash
devd dl https://example.com/images/favicon.ico
```

```bash
devd download https://example.com/images/favicon.ico \
  --output-file favicon.ico --working-directory ~/Downloads --concurrent 4 
```

```bash
devd dl https://example.com/images/favicon.ico \
  -o logo.svg -D ~/Downloads -c 4
```

Notes:
- Priority download tools: aria2c > wget > curl
- The flag `--concurrent` (`-c`) is only used when aria2c is used as download tool
- Default concurrent download is 4 for speed up download process

### Generate SSH key-pair using ssh-keygen and ed25519
```bash
devd gen ssh-key file_name email@example.com
```

### Generate UFW rules
```bash
devd gen ufw --help
```

```bash
# Allow connect to :port from anywhere
devd gen ufw-allow [port]

# Allow connection to any port from a specified IP
devd gen ufw-allow [ip]

# Allow connects to a specified port from a specified IP
devd gen ufw-allow [ip] [port]
```

_Pre-defined ports: http=80, https=443, ssh=22, db=5432, grpc=9090, rpc=26657, evm=8545, p2p=26656, rest=1317_

### Checking tools used by Devd
```bash
devd verify-tools
```

### Security check
```bash
sudo devd security-check
```
