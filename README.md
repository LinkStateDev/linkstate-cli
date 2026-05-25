# linkstate-cli

Command-line tool for the LinkState learning platform.

## Install

```bash
go install github.com/LinkStateDev/linkstate-cli@latest
```

Or:

```bash
make install        # builds to $GOPATH/bin/lst
```

## Commands

| Command | Description |
|---------|-------------|
| `lst auth` | Authenticate via browser (OAuth callback) |
| `lst fetch <slug>` | Download a lesson |
| `lst test` | Run compiled test binary locally |
| `lst submit` | Run tests and submit result |
| `lst hint [level]` | Get a hint from the server |
| `lst config` | Show or change settings |
| `lst progress` | Show learning progress |

## Project Structure

```
cmd/              ← Cobra commands (auth, fetch, test, submit, hint, config, progress)
internal/
├── client/       ← HTTP client with JWT auth
├── color/        ← ANSI color helpers
└── config/       ← Configuration store (~/.linkstate/config.json)
```
