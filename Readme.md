# Internet-Relay-Chat

Small test IRC client + TUI (Bubble Tea) project that demonstrates:
- an `irc` package for connecting to IRC servers (connect / auth / join / send / read)
- a Bubble Tea TUI for asking questions and sending messages
- wiring between TUI and IRC via channels

--- 

## Repo layout (important files)
- `main.go` — program entry, TUI + IRC orchestration (starts TUI and IRC goroutines).
- `irc/irc.go` — IRC client library (Init, Connect, Auth, Join, Say, SendRaw, GetResponse).
- `tui.go`, `wizard.go`, `main.go` (TUI parts) — Bubble Tea model, question/input components.
- `go.mod` — module file (currently `module irc` in this repo).

---

## Quick start (Windows PowerShell)

1. From project root:
   - Initialize or fix module path (recommended to use a canonical path). Example:
     ```powershell
     # choose a module path (replace with your github path)
     go mod edit -module=github.com/you/Internet-Relay-Chat
     go mod tidy
     ```
     If you keep `module irc`, imports must use that prefix (e.g. `irc/irc`, `irc/tui`).

2. Build and run:
   ```powershell
   go mod tidy
   go run .\main.go
   # or build
   go build -o ircclient .
   .\ircclient.exe
   ```

---

## Configuration

Edit `main.go` and set these constants (or replace with your preferred config handling):

- `domain` — IRC server domain (e.g. `irc.libera.chat`, `irc.oftc.net`)
- `port` — `6667` for plain TCP or `6697` for TLS (TLS requires using `tls.Dial`)
- `user` — username (single token, not the whole USER line)
- `nick` — desired nickname
- `channel` — channel name (without `#`) used by the TUI

Examples:
```go
const (
    domain = "irc.oftc.net"
    port   = "6667"
    user   = "building101"
    nick   = "building101"
)
```

---

## How it works (integration notes)

- TUI runs with Bubble Tea. The program starts the TUI in a goroutine (or with `p.Start()`), then starts IRC goroutines.
- Communication between TUI and IRC uses channels:
  - `ircIn`  (chan string) — incoming raw server lines forwarded into the TUI.
  - `ircOut` (chan string) — outgoing raw IRC commands produced by the TUI and consumed by an IRC writer goroutine.

Design rules:
- Only one goroutine should read from the TCP connection. Use a single reader goroutine that calls `GetResponse()` and forwards lines to TUI (`p.Send(...)`) and/or `ircIn`.
- Use one writer goroutine that reads from `ircOut` and calls `SendRaw` / `Send_data`.
- `Send_data` should append CRLF (`\r\n`) once; callers must NOT include CRLF.

---

## Important IRC protocol notes / common fixes

- Use CRLF for IRC lines: `"COMMAND args\r\n"`. Your send helper should append `\r\n`.
- Registration handshake:
  - PASS (optional) — `PASS <password>`
  - NICK — `NICK <nick>`
  - USER — `USER <username> 0 * :<realname>`
- PING/PONG: respond to server PING lines with the same payload: `PONG <payload>` or `PONG :<payload>`.
- PRIVMSG format: `PRIVMSG <target> :<message>`
  - Example to channel: `PRIVMSG #testchannel :hello world`
  - Wrong: `PRIVMSG #testchannel:hello` (colon in wrong place) — server will interpret as a nickname target.
- If server says `No Ident response` or `Could not resolve your hostname` — informational; ident/ptr lookups are optional.
- If server disconnects with `Registration timed out` — likely malformed handshake or CRLF missing.

---

## TUI / Bubble Tea notes & common issues

- Placeholder not shown / truncated:
  - Ensure you construct and configure the textinput for each question (set `Placeholder`, `Focus()`, `Width`, `CharLimit`).
  - Example: `ti := textinput.New(); ti.Placeholder = "Enter your answer"; ti.Focus(); ti.Width = 60`
- Use pointer receivers on your Bubble Tea model methods (`*model`) so state (like input model) is preserved.
- Forward key messages to the child text input component:
  ```go
  m.answerField, cmd = m.answerField.Update(msg)
  ```
- To change input style when empty (e.g., red border), modify the style in `View()` depending on `strings.TrimSpace(m.answerField.Value()) == ""`.
- Do not call `GetResponse()` or `reader.ReadLine()` from multiple goroutines (textproto.Reader is not concurrency-safe).

---

## Debugging tips

- Log outgoing raw IRC commands: print the raw line before sending (e.g. `fmt.Println(">>>", line)`).
- Observe server replies in the TUI — forward every server line into the TUI with `p.Send(IrcMsg(line))`.
- If PRIVMSG is ignored, check:
  - Did Join succeed? Wait for `JOIN` confirmation before sending.
  - Channel modes (e.g., +m) or ChanServ access lists may prevent you from sending; server returns numeric error codes (401/403/404).
- If you see `No such nick/channel` it usually means the server thought the target was a nick (bad PRIVMSG syntax).

---

## Tests & development workflow

- Use `go vet` / `golangci-lint` for static checks.
- Use `go test ./...` if you add tests.
- Use `go mod tidy` after changing imports or module path.

---

## Known TODO & suggestions

- Move to TLS for networks requiring SSL (use `crypto/tls` and `tls.Dial`).
- Implement SASL auth if the network requires it.
- Improve message parsing: expose parsed msg struct with exported fields.
- Add proper shutdown handling: close channels and goroutines gracefully on quit.
- Consider splitting TUI and IRC into separate modules/packages if you want independent reuse.

---

If you want, I can:
- Create/replace this README in the repo (`d:\Internet_relay_chat\Internet-Relay-Chat\README.md`).
- Patch `go.mod` to a canonical module path and update imports in `main.go`.
- Produce a small example `config.go` to hold domain/port/user/nick constants.