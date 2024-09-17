package engine

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"net"
	"os"
	"strings"
	"time"

	"github.com/nzin/taxsi2/internal/com"
	"github.com/nzin/taxsi2/internal/db"
	"github.com/sirupsen/logrus"
)

type WafEngine interface {
	/*
	 register scanning plugin
	*/
	RegisterPlugin(plugin WafEnginePlugin)
	/*
	 main WAF scanning function
	 Returns true if we dont block
	*/
	Scan(payload *com.TaxsiCom) bool
}

type WafEnginePlugin interface {
	// WAF Plugin name
	Name() string
	/*
	 WAF Plugin scanning function
	 Returns true if we dont block
	*/
	Scan(payload *com.TaxsiCom) bool
}

type WafEngineImpl struct {
	/*
	  analysisOutput is a comma separated list of
	  output (ALL:stdout, BLOCKED:/filepath ... ) to log scan events
	*/
	analysisOutput []WafOuput
	/*
	  scan events format string:
	  - {{.Date}} (human readable)
	  - {{.Timestamp}} (UTC timestamp)
	  - {{.Url}}
	  - {{.UrlHostname}}
	  - {{.UrlPath}}
	  - {{.Method}}
	  - {{.Remoteaddr}}
	  - {{.Scanresult}} (blocked, dryrun, pass)
	*/
	analysisOutputTemplate *template.Template
	config                 WafConfig
	plugins                map[string]WafEnginePlugin
}

type WafOuput struct {
	OutputType string
	Writer     io.Writer
}

func NewWafEngineImpl(ds db.DbService, analysisOutput string, analysisOutputFormat string) (WafEngine, error) {
	config, err := NewWafConfig(ds)
	if err != nil {
		return nil, err
	}

	tmpl, err := template.New("outputformat").Parse(analysisOutputFormat)
	if err != nil {
		return nil, fmt.Errorf("error analysis output format: %v", err)
	}

	wafoutputs := []WafOuput{}
	outputs := strings.Split(analysisOutput, ",")
	for _, output := range outputs {
		outputtype := "all"
		path := "logrus"
		if strings.HasPrefix(output, "all:") {
			outputtype = "all"
			path = output[4:]
		}
		if strings.HasPrefix(output, "blocked:") {
			outputtype = "blocked"
			path = output[8:]
		}

		if path == "logrus" {
			wafoutputs = append(wafoutputs, WafOuput{
				OutputType: outputtype,
				Writer:     nil,
			})
		} else if path == "stdout" {
			wafoutputs = append(wafoutputs, WafOuput{
				OutputType: outputtype,
				Writer:     os.Stdout,
			})
		} else {
			f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND, 0644)
			if err != nil {
				return nil, fmt.Errorf("not able to open %s for output: %v", path, err)
			}
			wafoutputs = append(wafoutputs, WafOuput{
				OutputType: outputtype,
				Writer:     f,
			})
		}
	}

	return &WafEngineImpl{
		analysisOutput:         wafoutputs,
		analysisOutputTemplate: tmpl,
		config:                 *config,
		plugins:                make(map[string]WafEnginePlugin),
	}, nil
}

func (we *WafEngineImpl) RegisterPlugin(plugin WafEnginePlugin) {
	we.plugins[plugin.Name()] = plugin
}

/*
main scanning function
*/
func (we *WafEngineImpl) Scan(payload *com.TaxsiCom) bool {
	res := true

	if we.config.Mode == "disabled" {
		we.output(payload, "pass")
		return true
	}

	// deny/allow
	remoteAddr := net.ParseIP(payload.RemoteAddr)
	if remoteAddr != nil {
		if we.config.IsIpAllowListed(remoteAddr) {
			we.output(payload, "pass")
			return true
		}
		if we.config.IsIpDenyListed(remoteAddr) {
			we.output(payload, "blocked")
			return false
		}
	}

	// Engine scan
	for name, plugin := range we.plugins {
		if we.config.EnabledPlugin[name] {
			res = plugin.Scan(payload)
			// if the plugin attempt to block
			if !res {
				if we.config.Mode == "enabled" {
					we.output(payload, "blocked")
					return false
				} else {
					we.output(payload, "dryrun")
					return true
				}
			}
		}
	}

	we.output(payload, "pass")

	return res
}

type OutputVariables struct {
	Date        string
	Timestamp   string
	Url         string
	UrlHostname string
	UrlPath     string
	Method      string
	Remoteaddr  string
	Scanresult  string
}

func (we *WafEngineImpl) output(payload *com.TaxsiCom, scanresult string) {
	now := time.Now()
	v := OutputVariables{
		Date:        now.Format(time.RFC3339),
		Timestamp:   fmt.Sprintf("%d", now.Unix()),
		Url:         payload.Url.String(),
		UrlHostname: payload.Url.Host,
		UrlPath:     payload.Url.RawPath,
		Method:      payload.Method,
		Remoteaddr:  payload.RemoteAddr,
		Scanresult:  scanresult,
	}

	for _, o := range we.analysisOutput {
		// if we want to output only blocked queries
		if o.OutputType == "blocked" && scanresult != "blocked" {
			continue
		}

		// if the output is logrus (then o.Writer is nil)
		if o.Writer == nil {
			var output bytes.Buffer
			_ = we.analysisOutputTemplate.Execute(&output, v)

			// Log the output as a string
			logrus.Info(output.String())

		}

		_ = we.analysisOutputTemplate.Execute(o.Writer, v)
	}
}
