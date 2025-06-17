package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/repo-scm/playground/sandbox"
)

var (
	serveAddress string
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run the server",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		if err := runServe(ctx, serveAddress); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
	},
}

// nolint:gochecknoinits
func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.PersistentFlags().StringVarP(&serveAddress, "address", "a", ":9090", "serve address")
}

func runServe(_ context.Context, address string) error {
	box, err := sandbox.NewSandbox()
	if err != nil {
		return errors.Wrap(err, "failed to initialize sandbox\n")
	}

	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./static")

	api := r.Group("/api")
	{
		api.POST("/containers", box.CreateContainer)
		api.GET("/containers", box.ListContainers)
		api.POST("/containers/:id/start", box.StartContainer)
		api.POST("/containers/:id/stop", box.StopContainer)
		api.DELETE("/containers/:id", box.RemoveContainer)
		api.GET("/containers/:id/logs", box.GetLogs)
	}

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	fmt.Printf("Server starting on %s\n", address)

	return r.Run(address)
}
