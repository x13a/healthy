package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/textproto"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	Version = "0.1.1"

	ExitSuccess = 0
	ExitFailure = 1
	ExitUsage   = 2

	FlagTimeout = "t"
	FlagHeader  = "H"
	FlagFail    = "f"
	FlagSkip    = "s"

	Hostname   = "127.0.0.1"
	DefaultURL = "http://" + Hostname + ":8000/ping"
)

type Header http.Header

func (h Header) Set(s string) error {
	header := strings.SplitN(s, ":", 2)
	if len(header) != 2 {
		return fmt.Errorf("Invalid header: %q", s)
	}
	for index, value := range header {
		header[index] = strings.TrimSpace(value)
	}
	textproto.MIMEHeader(h).Add(header[0], header[1])
	return nil
}

func (h Header) String() string {
	return ""
}

type Opts struct {
	url     string
	timeout time.Duration
	header  Header
	fail    bool
	skip    bool
}

func newOpts() *Opts {
	return &Opts{header: make(Header)}
}

func usage() {
	var name string
	if len(os.Args) == 0 {
		name = "PROG_NAME"
	} else {
		name = filepath.Base(os.Args[0])
	}
	fmt.Fprintf(
		flag.CommandLine.Output(),
		"%s [URL (default: %s)]\n",
		name,
		DefaultURL,
	)
	flag.PrintDefaults()
}

func parseURL(opts *Opts) error {
	opts.url = flag.Arg(0)
	if opts.url == "" {
		opts.url = DefaultURL
	}
	uri, err := url.ParseRequestURI(opts.url)
	if err != nil {
		return err
	}
	if uri.Scheme != "http" && uri.Scheme != "https" {
		return fmt.Errorf("Invalid scheme: %q", uri.Scheme)
	}
	hostname := uri.Hostname()
	if hostname != Hostname {
		return fmt.Errorf("Invalid hostname: %q", hostname)
	}
	return nil
}

func getOpts() *Opts {
	flag.Usage = usage
	opts := newOpts()
	isVersion := flag.Bool("V", false, "Print version and exit")
	flag.DurationVar(&opts.timeout, FlagTimeout, 0, "Timeout")
	flag.Var(&opts.header, FlagHeader, "Header")
	flag.BoolVar(&opts.fail, FlagFail, true, "Fail silently")
	flag.BoolVar(&opts.skip, FlagSkip, true, "InsecureSkipVerify")
	flag.Parse()
	if *isVersion {
		fmt.Println(Version)
		os.Exit(ExitSuccess)
	}
	if err := parseURL(opts); err != nil {
		fmt.Fprintln(flag.CommandLine.Output(), err)
		os.Exit(ExitUsage)
	}
	return opts
}

func request(opts *Opts) error {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: opts.skip}
	client := &http.Client{Transport: transport, Timeout: opts.timeout}
	req, err := http.NewRequest(http.MethodGet, opts.url, nil)
	if err != nil {
		return err
	}
	for k, vs := range opts.header {
		for _, v := range vs {
			req.Header.Add(k, v)
		}
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < http.StatusOK || resp.StatusCode > 299 {
		return fmt.Errorf("HTTP status code: %d", resp.StatusCode)
	}
	return nil
}

func main() {
	opts := getOpts()
	if err := request(opts); err != nil {
		if !opts.fail {
			log.Println(err)
		}
		os.Exit(ExitFailure)
	}
}
