//nolint:errcheck // ignore error
package metrics

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewHTTPMetricsServer(t *testing.T) {
	r := require.New(t)

	server, err := NewDefaultHTTPServer(NewDefaultHTTPServerConfig())
	r.NoError(err)

	defer server.Stop()
	go func() { err = server.Start(); r.NoError(err) }()

	resp, err := http.Get("http://localhost:8019/metrics")
	r.NoError(err)
	r.Equal(http.StatusOK, resp.StatusCode)
}
