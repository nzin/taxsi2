package engine

import (
	"fmt"
	"net"
	"strings"

	"github.com/nzin/taxsi2/internal/db"
	"github.com/sirupsen/logrus"
)

/*
Taxsi2 main (global) configuration
*/
type WafConfig struct {
	ds            db.DbService
	Mode          string // enabled, dryrun, disabled. (learning in the future?)
	EnabledPlugin map[string]bool
	AllowList     []*net.IPNet
	DenyList      []*net.IPNet
}

func NewWafConfig(ds db.DbService) (*WafConfig, error) {

	wc := WafConfig{
		ds:            ds,
		Mode:          "enabled",
		EnabledPlugin: make(map[string]bool),
		AllowList:     []*net.IPNet{},
		DenyList:      []*net.IPNet{},
	}

	if err := wc.loadConfigs(); err != nil {
		return nil, err
	}
	ds.SubscribeChanges(db.CHANGELOG_TABLE_CONFIG, &wc)
	return &wc, nil
}

func (wc *WafConfig) loadConfigs() error {
	configs, err := wc.ds.GetConfigs()
	if err != nil {
		return err
	}
	for k, v := range configs {
		wc.parseKeyValue(k, v)
	}
	return nil
}

func (wc *WafConfig) NotifyDbChange(key string) {
	value, err := wc.ds.GetConfigValueForKey(key)
	if err != nil {
		logrus.Errorf("Error reading config %s: %v", key, err)
		return
	}

	wc.parseKeyValue(key, value)
}

func (wc *WafConfig) parseKeyValue(k string, v string) {
	// mode
	if k == "mode" {
		if v == "enabled" || v == "dryrun" || v == "disabled" {
			wc.Mode = v
		}
	}
	// plugin enable
	if strings.HasPrefix(k, "plugin_") && (v == "enabled" || v == "disabled") {
		wc.EnabledPlugin[k[len("plugin_"):]] = v == "enabled"
	}

	// allow list
	if k == "allowlist" {
		wc.AllowList = []*net.IPNet{}
		nets := strings.Split(v, ",")
		for _, n := range nets {
			_, network, err := net.ParseCIDR(n)
			if err != nil {
				logrus.Errorf("not able to parse allow net %s: %v", n, err)
			} else {
				wc.AllowList = append(wc.AllowList, network)
			}
		}
	}

	// deny list
	if k == "denylist" {
		wc.DenyList = []*net.IPNet{}
		nets := strings.Split(v, ",")
		for _, n := range nets {
			_, network, err := net.ParseCIDR(n)
			if err != nil {
				logrus.Errorf("not able to parse deny net %s: %v", n, err)
			} else {
				wc.DenyList = append(wc.DenyList, network)
			}
		}
	}
}

/*
IsIpAllowListed to know if an IP is allow/white listed
*/
func (wc *WafConfig) IsIpAllowListed(remoteAddr net.IP) bool {
	for _, n := range wc.AllowList {
		if n.Contains(remoteAddr) {
			return true
		}
	}
	return false
}

/*
IsIpAllowListed to know if an IP is deny/black listed
*/
func (wc *WafConfig) IsIpDenyListed(remoteAddr net.IP) bool {
	for _, n := range wc.DenyList {
		if n.Contains(remoteAddr) {
			return true
		}
	}
	return false
}

func (wc *WafConfig) SetMode(mode string) error {
	if !(mode == "enabled" || mode == "dryrun" || mode == "disabled") {
		return fmt.Errorf("bad mode (must be enabled,dryrun or disabled)")
	}

	wc.Mode = mode
	wc.ds.SetConfigValueForKey("mode", mode)
	return nil
}

func (wc *WafConfig) SetAllowList(allowlist string) error {
	// check values
	nets := strings.Split(allowlist, ",")
	al := []*net.IPNet{}
	for _, n := range nets {
		_, network, err := net.ParseCIDR(n)
		if err != nil {
			return fmt.Errorf("not able to parse allow net %s: %v", n, err)
		} else {
			al = append(al, network)
		}
	}

	wc.AllowList = al
	return wc.ds.SetConfigValueForKey("allowlist", allowlist)
}

func (wc *WafConfig) SetDenyList(denylist string) error {
	// check values
	nets := strings.Split(denylist, ",")
	dl := []*net.IPNet{}
	for _, n := range nets {
		_, network, err := net.ParseCIDR(n)
		if err != nil {
			return fmt.Errorf("not able to parse allow net %s: %v", n, err)
		} else {
			dl = append(dl, network)
		}
	}

	wc.DenyList = dl
	return wc.ds.SetConfigValueForKey("denylist", denylist)
}

func (wc *WafConfig) EnablePlugin(pluginname string, enable bool) error {
	wc.EnabledPlugin[pluginname] = enable

	e := "enabled"
	if !enable {
		e = "disabled"
	}
	wc.ds.SetConfigValueForKey("plugin_"+pluginname, e)
	return nil
}
