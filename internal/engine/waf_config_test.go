package engine

import (
	"net"
	"testing"

	"github.com/nzin/taxsi2/internal/db"
	"github.com/stretchr/testify/assert"
)

type DbServiceConfigMock struct {
	config       map[string]string
	lastSetKey   string
	lastSetValue string
}

func (c *DbServiceConfigMock) SubscribeChanges(table int, listener db.DbChangeListener) {

}
func (c *DbServiceConfigMock) GetConfigs() (map[string]string, error) {
	return c.config, nil
}
func (c *DbServiceConfigMock) GetConfigValueForKey(key string) (string, error) {
	return c.config[key], nil
}
func (c *DbServiceConfigMock) SetConfigValueForKey(key string, value string) error {
	c.lastSetKey = key
	c.lastSetValue = value
	return nil
}

func TestWafConfig(t *testing.T) {
	t.Run("happy path: creating an empty waf config", func(t *testing.T) {
		c := DbServiceConfigMock{
			config: make(map[string]string),
		}

		_, err := NewWafConfig(&c)
		assert.Nil(t, err)
	})

	t.Run("happy path: creating waf config with parameters", func(t *testing.T) {
		c := DbServiceConfigMock{
			config: make(map[string]string),
		}
		c.config["mode"] = "enabled"
		c.config["plugin_foo"] = "enabled"
		c.config["allowlist"] = "8.8.8.8/32,1.1.1.0/24"
		c.config["denylist"] = "2.2.2.2/32"

		wc, err := NewWafConfig(&c)
		assert.Nil(t, err)
		assert.Equal(t, "enabled", wc.Mode)
		assert.Equal(t, 1, len(wc.EnabledPlugin))
		assert.Equal(t, true, wc.EnabledPlugin["foo"])
		assert.Equal(t, 2, len(wc.AllowList))
		assert.Equal(t, 1, len(wc.DenyList))
		assert.Equal(t, "2.2.2.2", wc.DenyList[0].IP.String())
	})

	t.Run("happy path: testing notification", func(t *testing.T) {
		c := DbServiceConfigMock{
			config: make(map[string]string),
		}
		c.config["mode"] = "enabled"
		c.config["plugin_foo"] = "enabled"
		c.config["allowlist"] = "8.8.8.8/32,1.1.1.0/24"
		c.config["denylist"] = "2.2.2.2/32"

		wc, err := NewWafConfig(&c)
		assert.Nil(t, err)

		c.config["allowlist"] = "4.4.4.4/32"
		wc.NotifyDbChange("allowlist")

		assert.Equal(t, "enabled", wc.Mode)
		assert.Equal(t, 1, len(wc.EnabledPlugin))
		assert.Equal(t, true, wc.EnabledPlugin["foo"])
		assert.Equal(t, 1, len(wc.AllowList))
		assert.Equal(t, 1, len(wc.DenyList))
	})

	t.Run("happy path: testing ips", func(t *testing.T) {
		c := DbServiceConfigMock{
			config: make(map[string]string),
		}
		c.config["mode"] = "enabled"
		c.config["plugin_foo"] = "enabled"
		c.config["allowlist"] = "8.8.8.8/32,1.1.1.0/24"
		c.config["denylist"] = "2.2.2.2/32"

		wc, err := NewWafConfig(&c)
		assert.Nil(t, err)

		remoteAddr := net.ParseIP("1.1.1.2")
		assert.True(t, wc.IsIpAllowListed(remoteAddr))
		remoteAddr = net.ParseIP("4.5.6.7")
		assert.False(t, wc.IsIpAllowListed(remoteAddr))

		remoteAddr = net.ParseIP("2.2.2.2")
		assert.True(t, wc.IsIpDenyListed(remoteAddr))
		remoteAddr = net.ParseIP("4.5.6.7")
		assert.False(t, wc.IsIpDenyListed(remoteAddr))
	})

	t.Run("happy path: testing set allowlist ips", func(t *testing.T) {
		c := DbServiceConfigMock{
			config: make(map[string]string),
		}
		c.config["mode"] = "enabled"
		c.config["plugin_foo"] = "enabled"
		c.config["allowlist"] = "8.8.8.8/32,1.1.1.0/24"
		c.config["denylist"] = "2.2.2.2/32"

		wc, err := NewWafConfig(&c)
		assert.Nil(t, err)

		remoteAddr := net.ParseIP("1.1.1.2")
		assert.True(t, wc.IsIpAllowListed(remoteAddr))
		remoteAddr = net.ParseIP("4.5.6.7")
		assert.False(t, wc.IsIpAllowListed(remoteAddr))

		err = wc.SetAllowList("4.4.4.4/")
		assert.NotNil(t, err)

		err = wc.SetAllowList("4.4.4.4/32,4.5.6.7/32")
		assert.Nil(t, err)
		assert.Equal(t, "allowlist", c.lastSetKey)
		assert.Equal(t, "4.4.4.4/32,4.5.6.7/32", c.lastSetValue)

		remoteAddr = net.ParseIP("1.1.1.2")
		assert.False(t, wc.IsIpAllowListed(remoteAddr))
		remoteAddr = net.ParseIP("4.5.6.7")
		assert.True(t, wc.IsIpAllowListed(remoteAddr))
	})

	t.Run("happy path: testing set deny ips", func(t *testing.T) {
		c := DbServiceConfigMock{
			config: make(map[string]string),
		}
		c.config["mode"] = "enabled"
		c.config["plugin_foo"] = "enabled"
		c.config["allowlist"] = "8.8.8.8/32,1.1.1.0/24"
		c.config["denylist"] = "2.2.2.2/32"

		wc, err := NewWafConfig(&c)
		assert.Nil(t, err)

		remoteAddr := net.ParseIP("2.2.2.2")
		assert.True(t, wc.IsIpDenyListed(remoteAddr))
		remoteAddr = net.ParseIP("4.5.6.7")
		assert.False(t, wc.IsIpDenyListed(remoteAddr))

		err = wc.SetDenyList("4.4.4.4/")
		assert.NotNil(t, err)

		err = wc.SetDenyList("4.5.0.0/16")
		assert.Nil(t, err)
		assert.Equal(t, "denylist", c.lastSetKey)
		assert.Equal(t, "4.5.0.0/16", c.lastSetValue)

		remoteAddr = net.ParseIP("2.2.2.2")
		assert.False(t, wc.IsIpDenyListed(remoteAddr))
		remoteAddr = net.ParseIP("4.5.6.7")
		assert.True(t, wc.IsIpDenyListed(remoteAddr))
	})

	t.Run("happy path: testing set mode", func(t *testing.T) {
		c := DbServiceConfigMock{
			config: make(map[string]string),
		}
		c.config["mode"] = "enabled"
		c.config["plugin_foo"] = "enabled"
		c.config["allowlist"] = "8.8.8.8/32,1.1.1.0/24"
		c.config["denylist"] = "2.2.2.2/32"

		wc, err := NewWafConfig(&c)
		assert.Nil(t, err)

		err = wc.SetMode("foobar")
		assert.NotNil(t, err)

		err = wc.SetMode("disabled")
		assert.Nil(t, err)
		assert.Equal(t, "mode", c.lastSetKey)
		assert.Equal(t, "disabled", c.lastSetValue)
	})

	t.Run("happy path: testing set plugin", func(t *testing.T) {
		c := DbServiceConfigMock{
			config: make(map[string]string),
		}
		c.config["mode"] = "enabled"
		c.config["plugin_foo"] = "enabled"
		c.config["allowlist"] = "8.8.8.8/32,1.1.1.0/24"
		c.config["denylist"] = "2.2.2.2/32"

		wc, err := NewWafConfig(&c)
		assert.Nil(t, err)

		err = wc.EnablePlugin("bar", true)
		assert.Nil(t, err)
		assert.Equal(t, "plugin_bar", c.lastSetKey)
		assert.Equal(t, "enabled", c.lastSetValue)
	})
}
