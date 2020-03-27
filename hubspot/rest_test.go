package hubspot

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestQuota(t *testing.T) {
	rest := &RestClient{
		quotatime: time.Duration(time.Millisecond * 100),
	}

	rest.BeginQuota()
	rest.EndQuota()
	start := time.Now().UTC()
	rest.BeginQuota()
	rest.EndQuota()
	end := time.Now().UTC()

	diff := end.Sub(start)
	require.True(t, diff > time.Millisecond*100)
}
