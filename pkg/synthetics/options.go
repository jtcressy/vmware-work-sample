package synthetics

import (
	"log"
	"net/url"
	"time"
)

type MultipleUrls []*url.URL

func (u *MultipleUrls) String() string {
    out := ""
    for _,url := range *u {
        out = out + "," + url.String()
    }
    return out
}

func (u *MultipleUrls) Set(value string) error {
    url, err := url.Parse(value)
    if err != nil {
        return err
    }
    *u = append(*u, url)
    return nil
}

type Options struct {
	FetchInterval time.Duration
	TestUrls MultipleUrls
    BindAddress string
}

func (o *Options) defaults() {
    if len(o.TestUrls) < 1 {
        if err := o.TestUrls.Set("https://httpstat.us/200"); err != nil {
            log.Fatalln(err)
        }
        if err := o.TestUrls.Set("https://httpstat.us/503"); err != nil {
            log.Fatalln(err)
        }
    }
}