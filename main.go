package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	flags "github.com/jessevdk/go-flags"
)

type options struct {
	GitHubAPIURI string `long:"github-api-uri" default:"https://api.github.com/licenses" default-mask:"-"  env:"GITHUB_API_URI" description:"URI for GitHub licenses API, without trailing slash"`
	ListFlag     bool   `short:"l" long:"list" description:"show license name list"`
	Username     string `short:"u" long:"username" env:"USERNAME" description:"Username to embed to LICENSE file"`
	Year         string `short:"y" long:"year" default-mask:"current year" description:"License year"`
	PosArgs      struct {
		License string `positional-arg-name:"LICENSE"`
	} `positional-args:"true"`
}

type license struct {
	Key     string `json:"key"`
	Name    string `json:"name"`
	SPDXID  string `json:"spdx_id"`
	URL     string `json:"url"`
	HTMLURL string `json:"html_url"`
}

type licenseDetail struct {
	Key            string   `json:"key"`
	Name           string   `json:"name"`
	SPDXID         string   `json:"spdx_id"`
	URL            string   `json:"url"`
	HTMLURL        string   `json:"html_url"`
	Description    string   `json:"description"`
	Implementation string   `json:"implementation"`
	Permissions    []string `json:"permissions"`
	Conditions     []string `json:"conditions"`
	Limitations    []string `json:"limitations"`
	Body           string   `json:"body"`
}

func main() { os.Exit(exec()) }

func exec() int {
	var opts options
	if _, err := flags.Parse(&opts); err != nil {
		if fe, ok := err.(*flags.Error); ok && fe.Type == flags.ErrHelp {
			return 0
		}

		printErr(err)
		return 1
	}
	if opts.ListFlag {
		return showList(opts)
	}
	return showLicense(opts)
}

func showList(opts options) int {
	res, err := http.Get(opts.GitHubAPIURI)
	if err != nil {
		printErr(err)
		return 1
	}
	if res.StatusCode != http.StatusOK {
		apiErr(res.StatusCode, res.Body)
		return 1
	}
	var licenses []license
	if err := json.NewDecoder(res.Body).Decode(&licenses); err != nil {
		printErr(err)
		return 1
	}
	var names []string
	for _, l := range licenses {
		names = append(names, l.Key)
	}
	fmt.Print(strings.Join(names, "\n"))
	return 0
}

func showLicense(opts options) int {
	if opts.PosArgs.License == "" || opts.Username == "" {
		printErr(fmt.Errorf("license and username is required"))
		return 1
	}
	uri := opts.GitHubAPIURI + "/" + strings.ToLower(opts.PosArgs.License)
	res, err := http.Get(uri)
	if err != nil {
		printErr(err)
		return 1
	}
	if res.StatusCode != http.StatusOK {
		apiErr(res.StatusCode, res.Body)
		return 1
	}
	var license licenseDetail
	if err := json.NewDecoder(res.Body).Decode(&license); err != nil {
		printErr(err)
		return 1
	}
	year := opts.Year
	if year == "" {
		year = strconv.Itoa(time.Now().Year())
	}
	text := license.Body
	text = strings.Replace(text, "[year]", year, -1)
	text = strings.Replace(text, "[fullname]", opts.Username, -1)
	fmt.Print(text)
	return 0
}

func apiErr(code int, body io.Reader) {
	b, err := ioutil.ReadAll(body)
	if err != nil {
		printErr(err)
	}
	printErr(fmt.Errorf("api error [%d] %s", code, string(b)))
}

func printErr(err error) {
	fmt.Fprintf(os.Stderr, "error occured: %s\n", err)
}
