# Gator CLI

Gator is a command-line tool to manage and interact with rss feeds. This README will guide you through the steps needed to install and run the tool.

## Requirements

Before using Gator, make sure you have the following installed:

1. **PostgreSQL** - [Install PostgreSQL](https://www.postgresql.org/download/) if you haven’t already.
2. **Go** - Gator is written in Go, so you'll need Go installed to compile and run it. [Install Go](https://go.dev/doc/install) if it’s not already on your system.

## Commands

Below is a list of commands available in the `gator` CLI tool. Each command has a brief description, usage syntax, and example.

### Command Overview

| Command     | Description                                                   |
| ----------- | ------------------------------------------------------------- |
| `login`     | Log in to the application.                                    |
| `register`  | Register a new user account.                                  |
| `users`     | Retrieve a list of all registered users.                      |
| `agg`       | Perform aggregation operations on rss collection and storage. |
| `addfeed`   | Add a new feed to your account.                               |
| `feeds`     | List all available feeds.                                     |
| `follow`    | Follow another feed.                                          |
| `following` | View the list of feeds you are following.                     |
| `unfollow`  | Unfollow a feed you are currently following.                  |
| `browse`    | Browse saved by agg feeds.                                    |
