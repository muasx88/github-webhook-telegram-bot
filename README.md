# Github Webhook Telegram Bot

Project to receive github notification payload and send to telegram via bot. This project require Github Secret and Telegram Bot Token.

#### Run the Applications following these steps:

1. clone the project
2. copy file .env-example to .env

```bash
$ cp .env-example .env
```

3. Fill the env

4. install modules

```bash
$ go mod tidy

#or
$ go mod download
```

5. run the application into development mode

```bash
$ make dev
```

#### Build process:

```bash
# run build command
$ make build

# execute
$ ./target/github-webhook-telegram-bot
```

#### Build process using Docker:

```bash
# run build command
$ make docker-up
```
