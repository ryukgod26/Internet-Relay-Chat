# Internet-Relay-Chat

An experimental IRC client learning project. It includes a small IRC client implementation in `irc/irc.go` and a Bubble Tea TUI (`tui.go` / `wizard.go`) that asks a short series of questions (domain, port, nick, etc.) and then sends a message.

This README explains how to run the project, how the TUI and IRC pieces interact, and common troubleshooting steps you can use while developing.

## Project layout

- `main.go` — program entry. Runs the TUI to collect connection info, then starts the IRC client and prints incoming server messages.
- `irc/irc.go` — small IRC client (Connect, Auth, Join, Say, SendRaw, GetResponse).
- `tui.go`, `wizard.go` — Bubble Tea-based TUI and question helpers.
- `go.mod` — module file (current module path: `irc`).

## Quick start (PowerShell)

1. From the project root:

```powershell
go mod tidy
go run .\main.go
```

2. The TUI will ask questions in order. Fill them and press Enter. After the final question the program connects and sends the message.

Notes:
- Provide the server `domain` (example: `irc.oftc.net`), `port` (example: `6667`), `username`, `nickname`, `channel` (without `#`), and optionally a password.
- If your module path is not `irc`, update `go.mod` and imports, or import local packages using the `irc/...` prefix as the project currently expects.

## How TUI and IRC are wired

- The program runs the TUI (Bubble Tea) to collect answers. `main.go` then reads answers from the model and initializes the IRC client with `irc.Init(domain, port, password, username, nick)`.
- The IRC connection is read by a single reader goroutine which calls `GetResponse()` and prints or forwards lines. Only one goroutine should read the TCP connection.
- Outgoing commands should be sent through the IRC client's send helpers (`SendRaw`, `Say`) — those helpers add the required CRLF (`\r\n`).

## Important protocol notes and common pitfalls

- Always use CRLF (`\r\n`) for IRC messages. The library's `send_data` appends CRLF — callers must not add another CRLF.
- Handshake sequence the server expects:
  1. (optional) `PASS <password>`
  2. `NICK <nick>`
  3. `USER <username> 0 * :<realname>`
- Respond to `PING` messages with `PONG <payload>`; the client includes a handler for that.
- PRIVMSG syntax is `PRIVMSG <target> :<message>` (note the space before the colon). Incorrect syntax (e.g. `PRIVMSG #chan:msg`) will usually cause a `401 No such nick/channel` error.

## Troubleshooting tips

- If your client is disconnected with `Registration timed out`, check that you send the correct `USER` and `NICK` lines and include CRLFs.
- If you see `No Ident response` or `Could not resolve your hostname` in server NOTICEs, those are informational (server tried an ident/PTR lookup). They don't usually prevent a working connection.
- If PRIVMSG replies aren't showing: verify the raw outgoing line printed by the client (the code prints outgoing lines). Ensure it matches `PRIVMSG #channel :message`.
- Avoid reading from the same `textproto.Reader` in multiple goroutines — it isn't concurrency-safe and can block or panic. Use one reader goroutine and forward received lines into the TUI or other channels.

## Development notes & recommended fixes

The repo is a work in progress. Here are recommended small fixes to stabilize behavior:

- Fix message parsing in `irc.ParseMessage` (prefer `" :"` as the message delimiter and slice using `idx+2` / `idx+1`).
- Ensure `send_data` is the single place CRLF is appended and remove extra `\r\n` strings from `NICK`/`PRIVMSG` calls.
- Use pointer receivers for Bubble Tea model methods (`func (m *model) Update(...)`) so state persists.
- Initialize text input components for each Question so placeholders show and input state is independent.
- When restoring questions (e.g. on ctrl+r), recreate fresh input components instead of copying the old structs.

If you'd like, I can apply any of these patches for you — tell me which one and I'll prepare the changes.

---

If anything here is out of date with your local code, tell me and I will update the README accordingly.