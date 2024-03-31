# About

Docker wrapper for [socks5 proxy-server](https://github.com/txthinking/socks5) by TxThinking

# Usage

```bash
docker run -d --name socks5 -p 1080:1080 -e PROXY_USERNAME=<username> -e PROXY_PASSWORD=<password> -e PROXY_IP=<proxy_ip> -e PROXY_PORT=<proxy_port>  nessai1/socks5:latest
```

| ENV variable   |Type| Default | Description                            |
|----------------|----|---------|----------------------------------------|
| PROXY_USERNAME |String| no value | Username that used for proxy auth      |
| PROXY_PASSWORD |String| no value | Password that used for proxy auth      |
| PROXY_IP       |String| 1080    | IP of socks5 proxy                     |
| PROXY_PORT     |String| no value | Port of socks5 proxy (1080 by default) |

**PROXY_IP** are always required.

If proxy use Username/password auth method - **PROXY_USERNAME** and **PROXY_PASSWORD** are both required.