# Paste Bay
### This is a simple server for creating pastes and sharing its shortened url with friends.
### BY [Islombay](https://t.me/islombay)

## Installation

First of all, clone the repository to your local computer, using git command:

```git clone https://github.com/islombay/PasteBay.git```

Next, create `.env` file that will contain environmental variables

```
DB_USER=postgresUsername
DB_PASSWORD=postgresPassword

AUTH_TOKENTTL=daysAmountForTokenLifeTime
AUTH_SECRETKEY=secretKeyForHashingToken
```

The code itself should create `data/` directory, but it would be better if
you create it before running the server.

- You can also specify the type of environment in which you are going to run the server. For this you need to change file `configs/config.json` and set one of the following values
  - `development`
  - `production`

#### Now, let's run!
Use command `go run cmd/main.go` to start the server.

## API

There is swagger available for this.
`/swagger/index.html`. Use this url path.