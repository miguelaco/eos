package cmd

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

type LoginCmd struct {
	client   *http.Client
	addr     string
	user     string
	password string
	*cobra.Command
}

func newLoginCmd() *cobra.Command {
	lc := LoginCmd{}

	lc.Command = &cobra.Command{
		Use:   "login",
		Short: "Perform login to EOS cluster.",
		Run: func(cmd *cobra.Command, args []string) {
			lc.addr = viper.GetString("addr")
			lc.user = viper.GetString("user")

			cookieJar, _ := cookiejar.New(nil)
			lc.client = &http.Client{
				Jar:       cookieJar,
				Transport: LogRedirects{},
				CheckRedirect: func(req *http.Request, via []*http.Request) error {
					if len(via) >= 10 {
						return errors.New("stopped after 10 redirects")
					}
					req.Header.Set("Referer", lc.addr)
					return nil
				},
			}

			lc.login()
		},
	}

	lc.Command.Flags().StringP("addr", "a", "", "Cluster url")
	lc.Command.Flags().StringP("user", "u", "admin", "Username you want to use")

	//	lc.Command.MarkFlagRequired("addr")

	viper.BindPFlag("addr", lc.Command.Flags().Lookup("addr"))
	viper.BindPFlag("user", lc.Command.Flags().Lookup("user"))

	return lc.Command
}

func (c *LoginCmd) login() {
	log.Printf("Login to %v as %v", c.addr, c.user)

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

	viper.Set("token", token)

	if err = viper.WriteConfig(); err != nil {
		log.Fatal(err)
		return
	}

	log.Println("Login successful")
}

func (c *LoginCmd) getLoginContext() (lc LoginContext, err error) {
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

func (c *LoginCmd) getAction(res *http.Response, formAction string) string {
	actionURL, _ := url.Parse(formAction)

	if !actionURL.IsAbs() {
		actionURL.Scheme = res.Request.URL.Scheme
		actionURL.Host = res.Request.URL.Host
	}

	return actionURL.String()
}

func (c *LoginCmd) getAuthToken(lc LoginContext) (string, error) {
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

func (c *LoginCmd) getCookie(name string) (string, error) {
	rootUrl, _ := url.Parse(c.addr)

	for _, cookie := range c.client.Jar.Cookies(rootUrl) {
		if cookie.Name == name {
			return cookie.Value, nil
		}
	}

	return "", errors.New("Cookie " + name + " not found")
}
