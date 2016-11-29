package mackerel

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

/*
{
  "monitors": [
    {
      "id": "2cSZzK3XfmG",
      "type": "connectivity",
      "isMute": false,
      "scopes": [],
      "excludeScopes": []
    },
    {
      "id"  : "2cSZzK3XfmG",
      "type": "host",
      "isMute": false,
      "name": "disk.aa-00.writes.delta",
      "duration": 3,
      "metric": "disk.aa-00.writes.delta",
      "operator": ">",
      "warning": 20000.0,
      "critical": 400000.0,
      "scopes": [
        "SomeService"
      ],
      "excludeScopes": [
        "SomeService: db-slave-backup"
      ],
      "notificationInterval": 60
    },
    {
      "id"  : "2cSZzK3XfmG",
      "type": "service",
      "isMute": false,
      "name": "SomeService - custom.access_num.4xx_count",
      "service": "SomeService",
      "duration": 1,
      "metric": "custom.access_num.4xx_count",
      "operator": ">",
      "warning": 50.0,
      "critical": 100.0
    },
    {
      "id"  : "2cSZzK3XfmG",
      "type": "external",
      "isMute": false,
      "name": "example.com",
      "url": "http://www.example.com",
      "service": "SomeService",
      "maxCheckAttempts": 1,
      "responseTimeCritical": 10000,
      "responseTimeWarning": 5000,
      "responseTimeDuration": 5,
      "certificationExpirationCritical": 15,
      "certificationExpirationWarning": 30,
      "containsString": "Example",
      "skipCertificateVerification": true
    }
  ]
}
*/

// Monitor represents interface to which each monitor type must confirm to.
type Monitor interface {
	MonitorType() string
	MonitorID() string
	MonitorName() string

	isMonitor()
}

const (
	monitorTypeConnectivity  = "connectivity"
	monitorTypeHostMeric     = "host"
	monitorTypeServiceMetric = "service"
	monitorTypeExternalHTTP  = "external"
	monitorTypeExpression    = "expression"
)

// Ensure each monitor type conforms to the Monitor interface.
var (
	_ Monitor = (*MonitorConnectivity)(nil)
	_ Monitor = (*MonitorHostMetric)(nil)
	_ Monitor = (*MonitorServiceMetric)(nil)
	_ Monitor = (*MonitorExternalHTTP)(nil)
	_ Monitor = (*MonitorExpression)(nil)
)

// Ensure only monitor types defined in this package can be assigned to the
// Monitor interface.
//
func (m *MonitorConnectivity) isMonitor()  {}
func (m *MonitorHostMetric) isMonitor()    {}
func (m *MonitorServiceMetric) isMonitor() {}
func (m *MonitorExternalHTTP) isMonitor()  {}
func (m *MonitorExpression) isMonitor()    {}

// MonitorConnectivity represents connectivity monitor.
type MonitorConnectivity struct {
	ID                   string `json:"id,omitempty"`
	Name                 string `json:"name,omitempty"`
	Type                 string `json:"type,omitempty"`
	IsMute               bool   `json:"isMute,omitempty"`
	NotificationInterval uint64 `json:"notificationInterval,omitempty"`

	Scopes        []string `json:"scopes,omitempty"`
	ExcludeScopes []string `json:"excludeScopes,omitempty"`
}

// MonitorType returns monitor type.
func (m *MonitorConnectivity) MonitorType() string { return monitorTypeConnectivity }

// MonitorName returns monitor name.
func (m *MonitorConnectivity) MonitorName() string { return m.Name }

// MonitorID returns monitor id.
func (m *MonitorConnectivity) MonitorID() string { return m.ID }

// MonitorHostMetric represents host metric monitor.
type MonitorHostMetric struct {
	ID                   string `json:"id,omitempty"`
	Name                 string `json:"name,omitempty"`
	Type                 string `json:"type,omitempty"`
	IsMute               bool   `json:"isMute,omitempty"`
	NotificationInterval uint64 `json:"notificationInterval,omitempty"`

	Metric   string  `json:"metric,omitempty"`
	Operator string  `json:"operator,omitempty"`
	Warning  float64 `json:"warning,omitempty"`
	Critical float64 `json:"critical,omitempty"`
	Duration uint64  `json:"duration,omitempty"`

	Scopes        []string `json:"scopes,omitempty"`
	ExcludeScopes []string `json:"excludeScopes,omitempty"`
}

// MonitorType returns monitor type.
func (m *MonitorHostMetric) MonitorType() string { return monitorTypeHostMeric }

// MonitorName returns monitor name.
func (m *MonitorHostMetric) MonitorName() string { return m.Name }

// MonitorID returns monitor id.
func (m *MonitorHostMetric) MonitorID() string { return m.ID }

// MonitorServiceMetric represents service metric monitor.
type MonitorServiceMetric struct {
	ID                   string `json:"id,omitempty"`
	Name                 string `json:"name,omitempty"`
	Type                 string `json:"type,omitempty"`
	IsMute               bool   `json:"isMute,omitempty"`
	NotificationInterval uint64 `json:"notificationInterval,omitempty"`

	Service  string  `json:"service,omitempty"`
	Metric   string  `json:"metric,omitempty"`
	Operator string  `json:"operator,omitempty"`
	Warning  float64 `json:"warning,omitempty"`
	Critical float64 `json:"critical,omitempty"`
	Duration uint64  `json:"duration,omitempty"`
}

// MonitorType returns monitor type.
func (m *MonitorServiceMetric) MonitorType() string { return monitorTypeServiceMetric }

// MonitorName returns monitor name.
func (m *MonitorServiceMetric) MonitorName() string { return m.Name }

// MonitorID returns monitor id.
func (m *MonitorServiceMetric) MonitorID() string { return m.ID }

// MonitorExternalHTTP represents external HTTP monitor.
type MonitorExternalHTTP struct {
	ID                   string `json:"id,omitempty"`
	Name                 string `json:"name,omitempty"`
	Type                 string `json:"type,omitempty"`
	IsMute               bool   `json:"isMute,omitempty"`
	NotificationInterval uint64 `json:"notificationInterval,omitempty"`

	URL                             string  `json:"url,omitempty"`
	MaxCheckAttempts                float64 `json:"maxCheckAttempts,omitempty"`
	Service                         string  `json:"service,omitempty"`
	ResponseTimeCritical            float64 `json:"responseTimeCritical,omitempty"`
	ResponseTimeWarning             float64 `json:"responseTimeWarning,omitempty"`
	ResponseTimeDuration            float64 `json:"responseTimeDuration,omitempty"`
	ContainsString                  string  `json:"containsString,omitempty"`
	CertificationExpirationCritical uint64  `json:"certificationExpirationCritical,omitempty"`
	CertificationExpirationWarning  uint64  `json:"certificationExpirationWarning,omitempty"`
	SkipCertificateVerification     bool    `json:"skipCertificateVerification,omitempty"`
}

// MonitorType returns monitor type.
func (m *MonitorExternalHTTP) MonitorType() string { return monitorTypeExternalHTTP }

// MonitorName returns monitor name.
func (m *MonitorExternalHTTP) MonitorName() string { return m.Name }

// MonitorID returns monitor id.
func (m *MonitorExternalHTTP) MonitorID() string { return m.ID }

// MonitorExpression represents expression monitor.
type MonitorExpression struct {
	ID                   string `json:"id,omitempty"`
	Name                 string `json:"name,omitempty"`
	Type                 string `json:"type,omitempty"`
	IsMute               bool   `json:"isMute,omitempty"`
	NotificationInterval uint64 `json:"notificationInterval,omitempty"`

	Expression string  `json:"expression,omitempty"`
	Operator   string  `json:"operator,omitempty"`
	Warning    float64 `json:"warning,omitempty"`
	Critical   float64 `json:"critical,omitempty"`
}

// MonitorType returns monitor type.
func (m *MonitorExpression) MonitorType() string { return monitorTypeExpression }

// MonitorName returns monitor name.
func (m *MonitorExpression) MonitorName() string { return m.Name }

// MonitorID returns monitor id.
func (m *MonitorExpression) MonitorID() string { return m.ID }

// FindMonitors find monitors
func (c *Client) FindMonitors() ([]Monitor, error) {
	req, err := http.NewRequest("GET", c.urlFor("/api/v0/monitors").String(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.Request(req)
	defer closeResponse(resp)
	if err != nil {
		return nil, err
	}

	var data struct {
		Monitors []json.RawMessage `json:"monitors"`
	}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}
	ms := make([]Monitor, 0, len(data.Monitors))
	for _, rawmes := range data.Monitors {
		m, err := decodeMonitor(rawmes)
		if err != nil {
			return nil, err
		}
		ms = append(ms, m)
	}
	return ms, err
}

// CreateMonitor creating monitor
func (c *Client) CreateMonitor(param Monitor) (Monitor, error) {
	resp, err := c.PostJSON("/api/v0/monitors", param)
	defer closeResponse(resp)
	if err != nil {
		return nil, err
	}
	return decodeMonitorReader(resp.Body)
}

// UpdateMonitor update monitor
func (c *Client) UpdateMonitor(monitorID string, param Monitor) (Monitor, error) {
	resp, err := c.PutJSON(fmt.Sprintf("/api/v0/monitors/%s", monitorID), param)
	defer closeResponse(resp)
	if err != nil {
		return nil, err
	}
	return decodeMonitorReader(resp.Body)
}

// DeleteMonitor update monitor
func (c *Client) DeleteMonitor(monitorID string) (Monitor, error) {
	req, err := http.NewRequest(
		"DELETE",
		c.urlFor(fmt.Sprintf("/api/v0/monitors/%s", monitorID)).String(),
		nil,
	)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := c.Request(req)
	defer closeResponse(resp)
	if err != nil {
		return nil, err
	}
	return decodeMonitorReader(resp.Body)
}

// decodeMonitor decodes json.RawMessage and returns monitor.
func decodeMonitor(mes json.RawMessage) (Monitor, error) {
	var typeData struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(mes, &typeData); err != nil {
		return nil, err
	}
	var m Monitor
	switch typeData.Type {
	case monitorTypeConnectivity:
		m = &MonitorConnectivity{}
	case monitorTypeHostMeric:
		m = &MonitorHostMetric{}
	case monitorTypeServiceMetric:
		m = &MonitorServiceMetric{}
	case monitorTypeExternalHTTP:
		m = &MonitorExternalHTTP{}
	case monitorTypeExpression:
		m = &MonitorExpression{}
	}
	if err := json.Unmarshal(mes, m); err != nil {
		return nil, err
	}
	return m, nil
}

func decodeMonitorReader(r io.Reader) (Monitor, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return decodeMonitor(b)
}
