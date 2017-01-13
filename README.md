# wp-version-to-slack

Simple Wordpress version checker that sends a message to the given Slack channel if the version differs from last run.

Designed to be run periodically in Cron or Jenkins et al. for informing about new Wordpress version availability.

## Install

```shell
go get github.com/eripa/wp-version-to-slack
```

## Usage

```shel
$ wp-version-to-slack -help
Usage: wordpress-to-slack -slack-token xoxp-1337-12345-67890

No output is sent unless a new version is found.

  -slack-channel string
        Slack Channel (without #) to post to (alternatively use env var SLACK_CHANNEL)
  -slack-token string
        Slack API token (u(alternatively use env var SLACK_TOKEN)
  -version
        Show tool version
  -wordpress-api string
        Wordpress API URL (default "https://api.wordpress.org/core/version-check/1.7/")
```

### Example

```shell
wp-version-to-slack -slack-token xoxp-1337-12345-67890 -slack-channel operations
2017/01/13 11:48:53 New version: 4.7.1
2017/01/13 11:48:53 Message successfully sent to channel ID CXXXYYYZZZ at 1484304533.000008
```

## Build

```shell
git clone http://github.com/eripa/wp-version-to-slack.git "$GOPATH/src/github.com/eripa/wp-version-to-slack"
cd "$GOPATH/src/github.com/eripa/wp-version-to-slack"
go get
go install
```

## License

The MIT License (MIT), See `LICENSE`
