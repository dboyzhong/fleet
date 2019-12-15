package kolide

import(
	"time"
	"encoding/json"
	_"fmt"
	"errors"
)

type HostnameType struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type OsMatch struct {
	Name      string `json:"name"`
	Accuracy  uint32 `json:"accuracy"`
	Line      uint32 `json:"line"`
}

type PService struct {
	NServer     string `json:"nServer"`
	State       string `json:"state"`
	Product     string `json:"product"`
	Version     string `json:"version"`
	HRef        string `json:"href"`       
	CPes        []string `json:"cpes"`
}

type ServicePort struct {
	Id         int `json:"id"` 
	Services   PService `json:"services"`
}

type PropertyParsed struct {
	Hostnames []*HostnameType `json:"hostnames"`
	OsMatches []*OsMatch      `json:"os_matches"`
	Time      time.Time       `json:"time"`
	Address   string          `json:"address"`
	//Ports     []*ServicePort  `json:"ports"`
	Ports     json.RawMessage  `json:"ports"`
}

/////////////////////////////

type PropertyAddr struct {
	Addr string `json:"addr"`
}

type PropertyOs struct {
	OsMatches  []*OsMatch `json:"os_matches"`
}

type PropertyHost struct {
	Os     PropertyOs  `json:"os"`
	Address    []*PropertyAddr `json:"addresses"` 
	Hostnames  []*HostnameType `json:"hostnames"`
	Ports      json.RawMessage `json:"ports"`
	OsFp   *string     `json:"os_fingerprints"`
}

type PropertyOrigin struct {
	Hosts []*PropertyHost `json:"hosts"`
}

func ParseProperty(origin, hostname string) (string, error) {

	var err error
	var propertyOrigin PropertyOrigin
	if err = json.Unmarshal([]byte(origin), &propertyOrigin); err != nil {
		return "", err
	}

	if (propertyOrigin.Hosts == nil || len(propertyOrigin.Hosts) == 0) {
		return "", errors.New("parse property no property found")
	}

	propertyParsed := make([]*PropertyParsed, 0, len(propertyOrigin.Hosts))

	for _, v := range(propertyOrigin.Hosts) {

		hostnames := make([]*HostnameType, 0)
		hn := &HostnameType{
			Name: hostname,
		}
		hostnames = append(hostnames, hn)

		osMatches := v.Os.OsMatches
		if (v.Os.OsMatches == nil) {
			osMatches = make([]*OsMatch, 0)
		}

		var addr string
		if (v.Address != nil && len(v.Address) > 0) {
			addr = (v.Address)[0].Addr
		}

		ports := v.Ports
		if (v.Ports == nil) {
			ports = make(json.RawMessage, 0)
		}

		pParsed := &PropertyParsed{
			Hostnames : hostnames,
			OsMatches : osMatches,
			Time      : time.Now(),
			Address   : addr,
			Ports     : ports,
		}

		propertyParsed = append(propertyParsed, pParsed)
	}

	var res []byte 
	if res, err = json.Marshal(propertyParsed); err != nil {
		return "", err
	}
	return string(res), nil
}