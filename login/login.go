package login

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"

	"github.com/miguelaco/eos/flag"

	"github.com/mitchellh/cli"
)

type LogRedirects struct {
	Transport http.RoundTripper
}

func (l LogRedirects) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	t := l.Transport
	if t == nil {
		t = http.DefaultTransport
	}
	resp, err = t.RoundTrip(req)
	if err != nil {
		return
	}

	log.Println(req.Method, "for", req.URL, "status", resp.StatusCode)

	return
}

type LoginContext struct {
	Action    string
	Lt        string
	Execution string
}

func (lc *LoginContext) Form() (form url.Values) {
	form = url.Values{}
	form.Add("lt", lc.Lt)
	form.Add("_eventId", "submit")
	form.Add("execution", lc.Execution)
	form.Add("submit", "LOGIN")

	return
}

const authCookieName = "dcos-acs-auth-cookie"

type Cmd struct {
	ui       cli.Ui
	flags    *flag.FlagSet
	client   *http.Client
	addr     string
	user     string
	password string
	help     string
}

func New(ui cli.Ui) *Cmd {
	cmd := &Cmd{ui: ui}
	cmd.init()
	return cmd
}

func (c *Cmd) init() {
	c.flags = flag.NewFlagSet("login", c.ui)
	c.flags.StringVar(&c.addr, "addr", "http://mycluster.example.com", "Cluster HTTP address")
	c.flags.StringVar(&c.user, "user", "admin", "Username")
	c.flags.StringVar(&c.password, "password", "", "Password")
	c.help = c.flags.Help(help)

	cookieJar, _ := cookiejar.New(nil)
	c.client = &http.Client{
		Jar:       cookieJar,
		Transport: LogRedirects{},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 10 {
				return errors.New("stopped after 10 redirects")
			}
			req.Header.Set("Referer", c.addr)
			return nil
		},
	}
}

func (c *Cmd) Run(args []string) int {
	c.flags.Parse(args)

	log.Printf("Login to %v as %v", c.addr, c.user)
	c.login()

	return 0
}

func (c *Cmd) login() {
	lc, err := c.getLoginContext()
	if err != nil {
		log.Fatal(err)
		return
	}

	token, err := c.getAuthToken(lc)
	if err != nil {
		log.Fatal(err)
		return
	}

	log.Println(authCookieName, token)
}

func (c *Cmd) getLoginContext() (lc LoginContext, err error) {
	lc = LoginContext{}

	addr := c.addr + "/login"
	res, err := c.client.Get(addr)
	if err != nil {
		log.Fatal(err)
		return
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
		return
	}

	re := regexp.MustCompile(`(?s)action="(.*?)".*name="lt" value="(.*?)".*name="execution" value="(.*?)"`)
	info := re.FindSubmatch(body)
	if info == nil {
		err = errors.New("Login info not found")
		return
	}

	lc.Action = c.getAction(res, string(info[1]))
	lc.Lt = string(info[2])
	lc.Execution = string(info[3])

	return
}

func (c *Cmd) getAction(res *http.Response, formAction string) string {
	actionURL, _ := url.Parse(formAction)

	if !actionURL.IsAbs() {
		actionURL.Scheme = res.Request.URL.Scheme
		actionURL.Host = res.Request.URL.Host
	}

	return actionURL.String()
}

func (c *Cmd) getAuthToken(lc LoginContext) (string, error) {
	form := lc.Form()
	form.Add("username", c.user)
	form.Add("password", c.password)
	form.Add("tenant", "NONE")

	_, err := c.client.PostForm(lc.Action, form)
	if err != nil {
		log.Fatal(err)
	}

	return c.getCookie(authCookieName)
}

func (c *Cmd) getCookie(name string) (string, error) {
	rootUrl, _ := url.Parse(c.addr)

	for _, cookie := range c.client.Jar.Cookies(rootUrl) {
		if cookie.Name == name {
			return cookie.Value, nil
		}
	}

	return "", errors.New("Cookie " + name + " not found")
}

func (c *Cmd) Synopsis() string {
	return synopsis
}

func (c *Cmd) Help() string {
	return c.help
}

const synopsis = "Login to EOS cluster"
const help = `
Usage: eos login [options]
  Login to EOS cluster.
`
