# Gator

Gator is a command-line application that manages user interactions and RSS feed aggregation. It includes features like user login, registration, feed management, and periodic RSS feed aggregation.

## Features

- **User Management**: Register, log in, reset users, and fetch user details.
- **RSS Feed Management**: Add feeds, follow/unfollow feeds, and periodically scrape RSS feeds.
- **Postgres Integration**: Connects to a PostgreSQL database to store user, feed, and post data.
- **Command-line Interface**: Run various commands to manage users, feeds, and posts.

## Technologies Used

- **Golang**: Core language for the application.
- **PostgreSQL**: Database used for storing user, feed, and post data.
- **SQLC**: Generates type-safe Go code to query the database.
- **Goose**: A tool for managing database migrations.
- **Go Modules**: For dependency management.

## Project Structure

.
├─ internal 
│  ├── commands   # Command handlers for CLI interactions 
│  ├── config     # Application configuration and settings 
│  ├── database   # Database queries and interaction 
│  ├── rss        # RSS feed fetching and parsing 
│  └── state      # Application state management 
└─ main.go        # Main entry point for the application

## Installation

1. Clone the repository:
   ```
   git clone https://github.com/acehotel33/bootdev-gator.git
   ```
2. Install the dependencies:
```
   git mod tidy
```
3. Set up your PostgreSQL database and update the ```.gatorconfig.json``` file with your database connection string:
```
  {
    "db_url": "your_postgres_connection_string",
    "current_user_name": ""
  }
```
4. Run the database migrations using goose:
```
   goose up
```
## Usage

### Available Commands
- User Commands:

* register [username] - Register a new user.
* login [username] - Log in with a specified username.
* reset - Reset all users.
* users - List all users.

- Feed Management:

* addfeed [feed_name] [feed_url] - Add a new RSS feed.
* feeds - List all RSS feeds.
* follow [feed_url] - Follow an RSS feed.
* following - List all feeds the user is following.
* unfollow [feed_url] - Unfollow a feed.
* browse [limit] - Browse posts from followed feeds.

- Aggregator:

* agg [interval] - Periodically collect RSS feeds with a given time interval (e.g., "1m" for 1 minute).

## Example Usage

### To register a new user:
```
go run main.go register johndoe
```
### To follow an RSS feed:
```
go run main.go follow https://example.com/rss
```
### To aggregate feeds every 1 minute:
```
go run main.go agg 1m
```
To roll back the last migration:
```
goose down
```
## Configuration
The application reads configuration from the .gatorconfig.json file located in your home directory. This file should include the PostgreSQL connection string and other necessary settings.

### Example .gatorconfig.json:
```
{
  "db_url": "postgres://user:password@localhost:5432/dbname",
  "current_user_name": ""
}
```
