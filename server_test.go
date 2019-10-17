package srv

import (
	"testing"

	"github.com/karlpokus/routest"
	"github.com/karlpokus/srv/testdata/routes"
)

func TestRoutes(t *testing.T) {
	routest.Test(t, []routest.Data{
		{
			"hi",
			nil,
			routes.Hello("dude"),
			200,
			[]byte("dude"),
		},
	})
}
