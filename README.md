# gator
GATOR CLI : Go App to gather and aggregate RSS feeds

## Prerequisites :
go 1.23.1
postgres 

## Installation

### 1. Install Go 1.23 or later
Gator CLI requires Golang installation, and only works on Linux. If you are on Windows, you will need to use WSL. Make sure you install go within your linux / wsl terminal and not your windows terminal.

**Option 1**: [The webi installer](https://webinstall.dev/golang/) is the simplest way for most people. Just run this in your terminal:
```bash
curl -sS https://webi.sh/golang | sh
```

**Option 2**: Use the [official installation instructions](https://go.dev/doc/install).

Run `go version` on your command line to ensure the installation worked. If it did _move on to step 2_.

### 2. Install Postgres 

Linux / WSL :
```bash
sudo apt update
sudo apt install postgresql postgresql-contrib
```

Ensure the installation worked. Make sure you are on version 16+ of Postgres:

```bash
psql --version
```

Setup a postgres password:
```bash
sudo passwd postgres
```
Enter a password, and be sure you won`t forget it.

Now, start your postgres server in the background.
```bash
sudo service postgresql start
```

Connect to the server and create a new database
```bash
sudo -u postgres psql
```

```sql
CREATE DATABASE gator;
```

Connect to your database

```sql
\c gator
```
Set the user password
```sql
ALTER USER postgres PASSWORD 'postgres';
```

Once everything is working, you can type `exit` to leave `psql` shell.

### 3. Setup your config file.

Create a `.gatorconfig.json` in your home directory.
It should look like so :
```json
{"db_url":"postgres://postgres:password@localhost:5432/gator?sslmode=disable","current_user_name":""}
```

### 4. Install gator

```bash
go install github.com/ablanchetmd/gator@latest
```
And you should be ready to use the CLI!

## Commands

```bash
gator register username
```
To register usernames, and set that username as the current user.

```bash
gator login username
```
login as current username to add feeds, follow, unfollow, see followed streams and browse posts.

```bash
gator reset
```
To reset the entire gator database from users, posts, feeds,

```bash
gator users
```
List users currently registered to the CLI app.

```bash
gator agg 10s
```
This is a heartbeat aggregator function where time is the time between each tick (1s, 1m, 1h).
This function will lock your current terminal in a loop and you'll need to use Ctrl+c to exit that terminal screen.

```bash
gator addfeed name url
```
Where name is the name of the rss feed and url is the link to it.

```bash
gator feeds
```
List all the feeds currently saved in the gator CLI.

```bash
gator follow url
```
Adds the feed to the current user`s feeds, where url is an url of a feed currently available in the feeds of the CLI

```bash
gator unfollow url
```
Removes the feed from the current user`s feeds, where url is the url of a feed.

```bash
gator following
```
List the feeds from the current user`s feeds.

```bash
gator browse 10
```
Display the 10 latest posts from your followed Rss feeds, where 10 can be changed to the number of posts you want to display at a time.

