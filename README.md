# wp-version-to-slack

[![Go Report Card](https://goreportcard.com/badge/github.com/eripa/wp-version-to-slack)](https://goreportcard.com/report/github.com/eripa/wp-version-to-slack)

![Example](example.png)

Simple Wordpress version checker that sends a message to the given Slack channel if the version differs from last run.

Designed to be run periodically in Cron or Jenkins et al. for informing about new Wordpress version availability.

## Install

```shell
go get github.com/eripa/wp-version-to-slack
```

## Usage

```shel
$ wp-version-to-slack -help
Usage: wp-version-to-slack -slack-token xoxp-1337-12345-67890

No output is sent unless a new version is found.

  -last-file string
        File for storing the previously known version (default "/tmp/wp-version-to-slack.last")
  -slack-channel string
        Slack Channel (without #) to post to (default is set to environment variable SLACK_CHANNEL)
  -slack-emoji string
        Slack message Emoji icon (default ":mailbox:")
  -slack-mention string
        Space separated list of @mentions (default is set to environment variable SLACK_MENTION)
  -slack-token string
        Slack API token (default is set to environment variable SLACK_TOKEN)
  -version
        Show tool version
  -wordpress-api string
        Wordpress API URL (default "https://api.wordpress.org/core/version-check/1.7/")
```

### Example

Simply send a message:

```shell
wp-version-to-slack -slack-token xoxp-1337-12345-67890 -slack-channel operations
2017/01/13 11:48:53 New version: 4.7.1
2017/01/13 11:48:53 Message successfully sent to channel ID CXXXYYYZZZ at 1484304533.000008
```

Mention specific persons:

```shell
wp-version-to-slack -slack-token xoxp-1337-12345-67890 -slack-channel operations -slack-mention "@eric"
2017/01/13 11:48:53 New version: 4.7.1
2017/01/13 11:48:53 Message successfully sent to channel ID CXXXYYYZZZ at 1484304533.000008
```

Using @channel or @here. Bots need to use special syntax for these, note the escaped ! to avoid shell expansion:

```shell
wp-version-to-slack -slack-token xoxp-1337-12345-67890 -slack-channel operations -slack-mention "<\!here>"
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
