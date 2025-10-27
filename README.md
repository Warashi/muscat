# muscat

muscat is a CLI that bridges a remote shell and your local workstation. It allows
you to reuse the local clipboard and browser from a remote session, and to expose
locally running TCP services to your remote environment through SSH.

## Highlights
- Remote clipboard integration (`muscat copy`, `muscat paste`)
- Local browser triggering (`muscat open`)
- Reverse port exposure via `muscat expose-port`, including automatic discovery
  of newly opened ports with `--auto`
- Compatibility with the original forward-only `muscat port-forward`

## Quick Start

### SSH configuration
```conf:~/.ssh/config
Host workbench
  Hostname workbench.remote.local
  RemoteForward /home/<remote_user>/.muscat.socket /Users/<local_user>/.muscat.socket
  ExitOnForwardFailure yes
```

### Start the muscat server on the remote host
```sh
$ muscat server
```

### Expose a local service to the remote host
```sh
# Expose localhost:3000 and request the same remote port.
$ muscat expose-port 3000
```

The command stays attached and streams logs. When the remote side connects,
muscat bridges the connection back to your local service.

To expose multiple ports explicitly:
```sh
$ muscat expose-port 3000 5432:15432
```

### Auto mode
```sh
# Watch for newly opened local ports and expose them automatically.
$ muscat expose-port --auto \
    --remote-policy next-free \
    --auto-exclude 22 \
    --auto-exclude-process ssh
```

muscat scans for listening TCP sockets on a fixed interval (default 5s) and
creates exposures for new ports. Manual ports supplied on the command line stay
active alongside auto-managed ports.

## Reverse Port Exposure Overview
```
+------------------+     SSH (unix socket forward)    +------------------+
| Local workstation | <=============================> | Remote host       |
|  muscat expose    |                                 |  muscat server    |
+------------------+                                 +------------------+
        |                                                       ^
        | 1. Local service listens on 127.0.0.1:3000            |
        |                                                       |
        v                                                       |
 Remote clients connect to remotehost:3000 ---------------------+
```

`muscat expose-port` runs a bidirectional stream with the server. The server
creates TCP listeners on demand and forwards each accepted connection to the
client. Data is multiplexed over a single gRPC stream, so multiple remote
connections can share one session efficiently.

### Key flags
- `--bind-address` – address used by the server when exposing ports
- `--public` – shortcut for `--bind-address 0.0.0.0`
- `--remote-policy` – conflict handling policy (`fail`, `next-free`, `skip`)
- `--auto-interval` – scan interval for `--auto`
- `--auto-exclude`, `--auto-exclude-process` – filters for the watcher

See `docs/expose-port.md` for workflow details, examples, and troubleshooting.

## Security Notes
- Binding to `0.0.0.0` (via `--public` or explicit `--bind-address`) exposes the
  service to every host that can reach the remote machine. Review firewall and
  authentication settings before enabling it.
- Only TCP is supported. UDP and IPv6 listeners are currently out of scope.
- Some operating systems require elevated privileges to enumerate listening
  ports for `--auto`. If muscat cannot list ports it will log the failure and
  skip auto registrations.

## Legacy `port-forward`
The classic forward-only `muscat port-forward` command remains available. Use it
when you need to reach services that reside on the remote network from your
local machine.

## Clipboard and Browser helpers
```sh
$ echo hello | muscat copy
$ muscat paste
$ muscat open 'https://example.com'
```

muscat can behave like `pbcopy`, `pbpaste`, and `xdg-open`. Create symbolic
links with those names to avoid retraining your muscle memory.
