# jump

# Example

edit `~/.ssh/config`

```json
Host *
 UseKeychain yes
 StrictHostKeyChecking no
 UserKnownHostsFile /dev/null

Host {hostName}_{env} # this is remarks
 User root
 HostName 127.0.0.1
```
