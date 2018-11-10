package cmd

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
)

type ServerCmd struct {
	*cobra.Command
	bindAddr string
}

func newServerCmd() *cobra.Command {
	sc := ServerCmd{}

	sc.Command = &cobra.Command{
		Use:   "server",
		Short: "Start test server",
		Run: func(cmd *cobra.Command, args []string) {
			router := mux.NewRouter()
			router.HandleFunc("/v1/sys/status", sc.healthHandler).Methods(http.MethodGet)
			router.HandleFunc("/login", sc.loginGetHandler).Methods(http.MethodGet)
			router.HandleFunc("/login", sc.loginPostHandler).Methods(http.MethodPost)

			log.Printf("Listening on: %s", sc.bindAddr)
			log.Fatal(http.ListenAndServe(sc.bindAddr, router))
		},
	}

	sc.Command.Flags().StringVarP(&sc.bindAddr, "bindAddr", "b", "", "Bind address to listen to")

	return sc.Command
}

func (sc *ServerCmd) healthHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Health Request: %s", r.URL.Path)

	type Status struct {
		Status string `json:"status"`
	}

	status := Status{"ok"}
	json.NewEncoder(w).Encode(status)

	s, _ := json.Marshal(status)
	log.Printf("Response: %s", s)
}

func (sc *ServerCmd) loginGetHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("GET Login Request: %s", r.URL.Path)

	s := `action="/login" name="lt" value="1234" name="execution" value="5678"`
	io.WriteString(w, s)

	log.Printf("Response: %s", s)
}

func (sc *ServerCmd) loginPostHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("POST Login Request: %s", r.URL.Path)

	r.ParseForm()
	log.Println("Form:", r.Form)

	user := r.Form.Get("username")
	password := r.Form.Get("password")

	if user == "admin" && password == "1234" {
		c := http.Cookie{Name: "dcos-acs-auth-cookie", Value: "12345678"}
		http.SetCookie(w, &c)
		log.Printf("Cookie: %s", c)
	}
}
