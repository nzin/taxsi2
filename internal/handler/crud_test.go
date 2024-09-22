package handler

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/nzin/taxsi2/internal/com"
	"github.com/nzin/taxsi2/internal/engine"
	"github.com/nzin/taxsi2/swagger_gen/restapi/operations/health"
	"github.com/nzin/taxsi2/swagger_gen/restapi/operations/waf"
	"github.com/stretchr/testify/assert"
)

func TestHealth(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		c := crud{
			ds:        nil,
			wafEngine: nil,
		}
		res := c.GetHealthcheck(health.GetHealthParams{})
		_, ok := res.(*health.GetHealthOK)
		assert.Equal(t, true, ok)
	})
}

type WafEngineMock struct {
	result bool
}

func (we *WafEngineMock) RegisterPlugin(plugin engine.WafEnginePlugin) {

}
func (we *WafEngineMock) Scan(payload *com.TaxsiCom) bool {
	return we.result
}

func TestHGetHealth(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		c := crud{
			ds:        nil,
			wafEngine: nil,
		}
		res := c.GetHealthcheck(health.GetHealthParams{})
		_, ok := res.(*health.GetHealthOK)
		assert.Equal(t, true, ok)
	})
}

func TestHPostSubmit(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		we := WafEngineMock{
			result: true,
		}
		c := crud{
			ds:        nil,
			wafEngine: &we,
		}

		r, err := http.NewRequest("GET", "https://foo/bar", nil)
		assert.Nil(t, err)

		payload, err := com.NewTaxsiCom(r)
		assert.Nil(t, err)

		var buf bytes.Buffer
		err = payload.Marshall(&buf)
		assert.Nil(t, err)

		call, err := http.NewRequest("POST", "https://taxsi2", &buf)
		assert.Nil(t, err)
		res := c.PostSubmit(waf.PostSubmitParams{
			HTTPRequest: call,
			Request:     io.NopCloser(&buf),
		})

		// returning 200?
		_, ok := res.(*waf.PostSubmitOK)
		assert.Equal(t, true, ok)
	})
}
