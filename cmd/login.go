package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/miguelaco/eos/common"
	"github.com/miguelaco/eos/config"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

const authCookieName = "dcos-acs-auth-cookie"

type LoginCmd struct {
	*cobra.Command
	client    *common.HttpClient
	cluster   *config.Cluster
	verbose   bool
	action    string
	lt        string
	execution string
	password  string
}

func newLoginCmd() *cobra.Command {
	lc := LoginCmd{}

	lc.Command = &cobra.Command{
		Use:   "login",
		Short: "Perform login to EOS cluster.",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			if err := lc.validate(); err != nil {
				return err
			}

			lc.client = common.NewHttpClient()
			lc.client.ResetReferer(true)
			lc.client.Verbose(lc.verbose)

			lc.login()

			return nil
		},
	}

	lc.cluster = config.GetAttachedCluster()
	lc.Command.Flags().StringVarP(&lc.cluster.User, "user", "u", "admin", "Username you want to use")
	lc.Command.Flags().BoolVarP(&lc.verbose, "verbose", "v", false, "Trace http requests")

	return lc.Command
}

func (c *LoginCmd) login() {
	fmt.Println("Login to", c.cluster.Addr, "as", c.cluster.User)

	if err := c.getForm(); err != nil {
		fmt.Println("Login error:", err)
		os.Exit(2)
	}

	c.promptPassword("Password: ")

	if err := c.postForm(); err != nil {
		fmt.Println("Login error:", err)
		os.Exit(2)
	}

	if err := config.Save(); err != nil {
		fmt.Println("Cannot write config:", err)
		os.Exit(3)
	}

	fmt.Println("Login successful")
}

func (c *LoginCmd) validate() error {
	missing := []string{}

	if !c.cluster.Active {
		fmt.Println("No attached cluster")
		os.Exit(2)
	}

	if c.cluster.Addr == "" {
		missing = append(missing, "addr")
	}

	if c.cluster.User == "" {
		missing = append(missing, "user")
	}

	if len(missing) > 0 {
		return fmt.Errorf(`"%s" not set as flags or config`, strings.Join(missing, `", "`))
	}

	return nil
}

func (c *LoginCmd) getForm() (err error) {
	addr := c.cluster.Addr + "/login"
	res, err := c.client.Get(addr)
	if err != nil {
		return
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}

	re := regexp.MustCompile(`(?s)action="(.*?)".*name="lt" value="(.*?)".*name="execution" value="(.*?)"`)
	info := re.FindSubmatch(body)
	if info == nil {
		err = errors.New("Login info not found")
		return
	}

	c.action = c.getAction(res, string(info[1]))
	c.lt = string(info[2])
	c.execution = string(info[3])

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

func (c *LoginCmd) postForm() (err error) {
	form := c.form()

	_, err = c.client.PostForm(c.action, form)
	if err != nil {
		return
	}

	c.cluster.Token, err = c.client.GetCookie(c.cluster.Addr, authCookieName)

	return
}

func (c *LoginCmd) form() (form url.Values) {
	form = url.Values{}
	form.Add("username", c.cluster.User)
	form.Add("password", c.password)
	form.Add("tenant", "NONE")
	form.Add("lt", c.lt)
	form.Add("_eventId", "submit")
	form.Add("execution", c.execution)
	form.Add("submit", "LOGIN")

	return
}

func (c *LoginCmd) promptPassword(msg string) {
	fmt.Print(msg)
	defer fmt.Print("\n")

	pass, _ := terminal.ReadPassword(int(os.Stdin.Fd()))
	c.password = string(pass)
}
