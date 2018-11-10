package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh/terminal"
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

	fmt.Println(req.Method, "for", req.URL, "status", resp.StatusCode)

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
		RunE: func(cmd *cobra.Command, args []string) error {
			lc.addr = viper.GetString("addr")
			lc.user = viper.GetString("user")

			if err := lc.validate(); err != nil {
				return err
			}

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
			return nil
		},
	}

	lc.Command.Flags().StringP("addr", "a", "", "Cluster url")
	lc.Command.Flags().StringP("user", "u", "admin", "Username you want to use")

	viper.BindPFlag("addr", lc.Command.Flags().Lookup("addr"))
	viper.BindPFlag("user", lc.Command.Flags().Lookup("user"))

	return lc.Command
}

func (c *LoginCmd) login() {
	fmt.Println("Login to", c.addr, "as", c.user)

	lc, err := c.getLoginContext()
	if err != nil {
		fmt.Println(err)
		return
	}

	c.password = c.promptPassword("Password: ")

	token, err := c.getAuthToken(lc)
	if err != nil {
		fmt.Println(err)
		return
	}

	viper.Set("token", token)

	if err = viper.WriteConfig(); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Login successful")
}

func (c *LoginCmd) validate() error {
	missing := []string{}

	if c.addr == "" {
		missing = append(missing, "addr")
	}

	if c.user == "" {
		missing = append(missing, "user")
	}

	if len(missing) > 0 {
		return fmt.Errorf(`"%s" not set as flags or config`, strings.Join(missing, `", "`))
	}

	return nil
}

func (c *LoginCmd) getLoginContext() (lc LoginContext, err error) {
	lc = LoginContext{}

	addr := c.addr + "/login"
	res, err := c.client.Get(addr)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
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
		fmt.Println(err)
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

func (c *LoginCmd) promptPassword(msg string) string {
	fmt.Print(msg)
	defer fmt.Print("\n")

	pass, _ := terminal.ReadPassword(int(os.Stdin.Fd()))
	return string(pass)
}
