package cmd

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os/exec"
	"runtime"
	"time"

	"github.com/LinkStateDev/linkstate-cli/internal/config"
	"github.com/LinkStateDev/linkstate-cli/internal/ui"
	"github.com/spf13/cobra"
)

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

		webURL := cfg.Server
		if webURL == "" {
			webURL = "http://localhost"
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
				if f, ok := w.(http.Flusher); ok {
					f.Flush()
				}
				// Hold the connection in "active" state briefly so the
				// browser actually receives the body before the main
				// goroutine tears the server down.
				time.Sleep(500 * time.Millisecond)
				errCh <- fmt.Errorf("callback missing token or email")
				return
			}
			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte("<html><body><h3>Authenticated!</h3><p>You can close this window.</p></body></html>"))
			if f, ok := w.(http.Flusher); ok {
				f.Flush()
			}
			time.Sleep(500 * time.Millisecond)
			resultCh <- authResult{token: token, email: email}
		})

		srv := &http.Server{Handler: mux}
		go func() {
			if err := srv.Serve(listener); err != http.ErrServerClosed {
				errCh <- err
			}
		}()

		fmt.Printf("%s %s\n", ui.Bold.Render("Opening browser:"), ui.Hint.Render(authURL))
		fmt.Println(ui.Muted.Render("If the browser does not open automatically, paste the link above."))
		openBrowser(authURL)

		var res authResult
		err = withSpinner("Waiting for browser callback…", func() error {
			select {
			case r := <-resultCh:
				res = r
				return nil
			case err := <-errCh:
				return err
			}
		})

		// Graceful shutdown waits for in-flight responses to finish so the
		// browser actually receives the "Authenticated!" page.
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		_ = srv.Shutdown(shutdownCtx)

		if err != nil {
			return err
		}

		cfg.Token = res.token
		cfg.Email = res.email
		if err := config.Save(cfg); err != nil {
			return fmt.Errorf("save config: %w", err)
		}
		fmt.Printf("%s %s\n", ui.Success.Render(ui.GlyphPass), ui.Bold.Render("Welcome back, "+res.email))
		return nil
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
	rootCmd.AddCommand(authCmd)
}
