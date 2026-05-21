package cmd

import (
	"fmt"
	"net"
	"net/http"
	"os/exec"
	"runtime"

	"github.com/LinkStateDev/linkstate-cli/internal/config"
	"github.com/spf13/cobra"
)

var authWebURL string

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authenticate via browser",
	RunE: func(cmd *cobra.Command, args []string) error {
		listener, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			return fmt.Errorf("start local server: %w", err)
		}
		port := listener.Addr().(*net.TCPAddr).Port
		callbackURL := fmt.Sprintf("http://localhost:%d", port)

		webURL := authWebURL
		if webURL == "" {
			webURL = "http://localhost:5173"
		}

		authURL := fmt.Sprintf("%s/login?callback=%s", webURL, callbackURL)

		resultCh := make(chan authResult, 1)
		errCh := make(chan error, 1)

		mux := http.NewServeMux()
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			token := r.URL.Query().Get("token")
			email := r.URL.Query().Get("email")
			if token == "" || email == "" {
				w.Header().Set("Content-Type", "text/html")
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("<html><body><h3>Authentication failed</h3><p>Missing token or email.</p><p>You can close this window.</p></body></html>"))
				errCh <- fmt.Errorf("callback missing token or email")
				return
			}
			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte("<html><body><h3>Authenticated!</h3><p>You can close this window.</p></body></html>"))
			resultCh <- authResult{token: token, email: email}
		})

		srv := &http.Server{Handler: mux}
		go func() {
			if err := srv.Serve(listener); err != http.ErrServerClosed {
				errCh <- err
			}
		}()

		fmt.Printf("Opening browser for authentication...\n")
		fmt.Printf("If the browser doesn't open, visit:\n  %s\n", authURL)
		openBrowser(authURL)

		select {
		case res := <-resultCh:
			srv.Close()
			cfg.Token = res.token
			cfg.Email = res.email
			if err := config.Save(cfg); err != nil {
				return fmt.Errorf("save config: %w", err)
			}
			fmt.Printf("Welcome back, %s!\n", res.email)
			return nil
		case err := <-errCh:
			srv.Close()
			return err
		}
	},
}

type authResult struct {
	token string
	email string
}

func openBrowser(url string) {
	var args []string
	switch runtime.GOOS {
	case "darwin":
		args = []string{"open", url}
	case "windows":
		args = []string{"cmd", "/c", "start", url}
	default:
		args = []string{"xdg-open", url}
	}
	exec.Command(args[0], args[1:]...).Start()
}

func init() {
	authCmd.Flags().StringVar(&authWebURL, "web", "", "Web app URL (default http://localhost:5173)")
	rootCmd.AddCommand(authCmd)
}
