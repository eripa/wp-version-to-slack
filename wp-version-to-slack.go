package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	"os"

	"net/http"
	"net/url"

	"github.com/nlopes/slack"
)

const (
	toolVersion = "0.1.0"
)

type wordpressResponse struct {
	Offers []struct {
		Response string `json:"response"`
		Download string `json:"download"`
		Locale   string `json:"locale"`
		Packages struct {
			Full       string `json:"full"`
			NoContent  string `json:"no_content"`
			NewBundled string `json:"new_bundled"`
			Partial    bool   `json:"partial"`
			Rollback   bool   `json:"rollback"`
		} `json:"packages"`
		Current        string `json:"current"`
		Version        string `json:"version"`
		PhpVersion     string `json:"php_version"`
		MysqlVersion   string `json:"mysql_version"`
		NewBundled     string `json:"new_bundled"`
		PartialVersion bool   `json:"partial_version"`
		NewFiles       bool   `json:"new_files,omitempty"`
	} `json:"offers"`
	Translations []interface{} `json:"translations"`
}

func notifySlack(api *slack.Client, channelID, message string) (err error) {
	params := slack.PostMessageParameters{
		Username: "Wordpress Version Checker",
		// IconEmoji: ":wizard:",
	}
	channel, timestamp, err := api.PostMessage(channelID, message, params)
	if err != nil {
		return err
	}
	log.Printf("Message successfully sent to channel ID %s at %s", channel, timestamp)
	return
}

func channelNametoID(api *slack.Client, name string) (ID string, err error) {
	channels, err := api.GetChannels(false)
	if err != nil {
		return "", err
	}
	for _, channel := range channels {
		if channel.Name == name {
			return channel.ID, nil
		}
	}
	return "", fmt.Errorf("no such channel: %s", name)
}

func getLatestWordpressVersion(apiURL string) (version string, err error) {
	u, err := url.Parse(apiURL)
	if err != nil {
		return version, err
	}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return version, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return version, err
	}
	defer resp.Body.Close()

	var body wordpressResponse
	err = json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		return version, err
	}

	for _, offer := range body.Offers {
		if offer.Response == "upgrade" {
			version = offer.Version
		}
	}
	if version == "" {
		return version, fmt.Errorf("could not find version, got offers: %+v", body.Offers)
	}
	return version, nil
}

func isNew(version string) (isnew bool, err error) {
	persistenceFile := "/tmp/wp-version-to-slack.last"
	noSuchFileError := fmt.Sprintf("stat %s: no such file or directory", persistenceFile)
	_, err = os.Stat(persistenceFile)
	if err != nil {
		switch err.Error() {
		case noSuchFileError:
			// assume no file means first run, return true as new version
			isnew = true
			vb := []byte(version)
			err = ioutil.WriteFile(persistenceFile, vb, 0644)
			if err != nil {
				return false, err
			}
		default:
			return false, err
		}
	}
	// Get stored version
	dat, err := ioutil.ReadFile(persistenceFile)
	if err != nil {
		return false, err
	}
	storedVersion := string(dat)
	if storedVersion != version {
		isnew = true
		// Update stored version
		vb := []byte(version)
		err = ioutil.WriteFile(persistenceFile, vb, 0644)
		if err != nil {
			return false, err
		}
	}
	return isnew, nil
}

func main() {
	versionCheck := flag.Bool("version", false, "Show tool version")
	slackToken := flag.String("slack-token", os.Getenv("SLACK_TOKEN"), "Slack API token (u(alternatively use env var SLACK_TOKEN)")
	slackChannel := flag.String("slack-channel", os.Getenv("SLACK_CHANNEL"), "Slack Channel (without #) to post to (alternatively use env var SLACK_CHANNEL)")
	wordpressAPI := flag.String("wordpress-api", "https://api.wordpress.org/core/version-check/1.7/", "Wordpress API URL")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s -slack-token xoxp-1337-12345-67890\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nNo output is sent unless a new version is found.\n\n")
		flag.PrintDefaults()
	}
	flag.Parse()
	if *versionCheck {
		fmt.Printf("%s v%s \\ʕ◔ϖ◔ʔ/\n", os.Args[0], toolVersion)
		os.Exit(0)
	}
	if *slackToken == "" || *slackChannel == "" {
		flag.Usage()
		os.Exit(1)
	}
	api := slack.New(*slackToken)
	channelID, err := channelNametoID(api, *slackChannel)
	if err != nil {
		log.Fatal(err)
	}
	version, err := getLatestWordpressVersion(*wordpressAPI)
	if err != nil {
		log.Fatal(err)
	}
	isNewVersion, err := isNew(version)
	if err != nil {
		log.Fatal(err)
	}
	if isNewVersion {
		log.Printf("New version: %s", version)
		notifySlack(api, channelID, fmt.Sprintf("New version available: %s", version))
	}
}