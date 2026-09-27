# go-mailer

A small Go program that sends a personalised email to every recipient in a CSV file, using a producer/consumer pipeline with a pool of concurrent workers.

## How it works

```
emails.csv ──► producer (loadRecipient) ──► channel ──► 5 workers (emailWorker) ──► SMTP
```

- **Producer** ([producer.go](producer.go)) reads `emails.csv`, skips the header row, and sends each row to a channel as a `Recipient`. It closes the channel when finished.
- **Workers** ([consumer.go](consumer.go)) read recipients from the channel, render [email.tmpl](email.tmpl) for each one, and send it with `net/smtp`.
- **main** ([main.go](main.go)) creates the channel, starts the producer and 5 workers, and waits for all workers to finish.

## Project layout

| File | Purpose |
|------|---------|
| `main.go` | Entry point, `Recipient` type, `executeTemplate` |
| `producer.go` | `loadRecipient`: reads the CSV into the channel |
| `consumer.go` | `emailWorker`: renders and sends each email |
| `email.tmpl` | Go `text/template` for the message (subject + body) |
| `emails.csv` | Campaign recipients |

## Requirements

- Go 1.25+
- Docker (to run a local SMTP server)

## Setup

### 1. Start a local SMTP server (Mailpit)

Mailpit catches all outgoing mail so nothing is actually delivered.

```sh
docker run -d \
  --restart unless-stopped \
  --name=mailpit \
  -p 8025:8025 \
  -p 1025:1025 \
  axllent/mailpit
```

- SMTP: `localhost:1025` (what the app sends to)
- Web UI: http://localhost:8025 (view the received emails)

### 2. Prepare the recipient CSV

`emails.csv` must have a header row followed by one recipient per line:

```csv
Name,Email
User1,user1@example.com
User2,user2@example.com
```

### 3. Edit the template (optional)

`email.tmpl` is rendered once per recipient. Fields from `Recipient` are available as `{{.Name}}` and `{{.Email}}`. The first line is the `Subject:` header, followed by a blank line and then the body.

## Run

```sh
go run .
```

Then open http://localhost:8025 to see the sent emails.

## Configuration

These values are currently hardcoded:

| Setting | Location | Default |
|---------|----------|---------|
| CSV path | `main.go` | `./emails.csv` |
| Worker count | `main.go` | `5` |
| SMTP host / port | `consumer.go` | `localhost:1025` |
| Sender address | `consumer.go` | `<youremail>@gmail.com` |
| Delay per email | `consumer.go` | `50ms` |
