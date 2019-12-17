package kolide

import(
	"time"
	"encoding/json"
	_"fmt"
	"errors"
	"strings"
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
	Fp          string   `json:"service_fp"`
}

type StateData struct {
	State string `json:"state"`
}

type ServicesPort struct {
	Id         int `json:"id"` 
	Services   PService  `json:"services"`
}

type ServicePort struct {
	Id         int `json:"id"` 
	Services   PService  `json:"service"`
	State      StateData `json:"state,omitempty"`
	Status     StateData `json:"status,omitempty"`
}

type PropertyParsed struct {
	Hostnames []*HostnameType `json:"hostnames"`
	OsMatches []*OsMatch      `json:"os_matches"`
	Time      time.Time       `json:"time"`
	Address   string          `json:"address"`
	Ports     []*ServicesPort  `json:"ports"`
	//Ports     json.RawMessage  `json:"ports"`
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
	Ports      []*ServicePort  `json:"ports"`
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
			ports = make([]*ServicePort, 0)
		}

		parsedPorts := make([]*ServicesPort, 0)

		for _, p := range ports {

			pP := &ServicesPort{}
			pP.Id = p.Id
			if p.Services.CPes == nil {
				pP.Services.CPes = make([]string, 0)
			} else {
				pP.Services.CPes = p.Services.CPes
			}
			pP.Services.Fp = p.Services.Fp
			pP.Services.Product = p.Services.Product
			pP.Services.Version = p.Services.Version

			if p.State.State == "open" || p.Status.State == "up" {
				pP.Services.State = "online"
			} else {
				pP.Services.State = "offline"
			}
			var ns string
			idx := strings.Index(p.Services.Fp, "nServer")
			if idx != -1 {
				nst := p.Services.Fp[idx + len("nServer"):]
				idx = strings.Index(nst, "\\r")
				if idx != -1 {
					ns = nst[:idx]
				}
			}
			pP.Services.NServer = ns	

			ns = ""
			idx = strings.Index(p.Services.Fp, "href")
			if idx != -1 {
				nst := p.Services.Fp[idx + len("href"):]
				idx = strings.Index(nst, "\\r")
				if idx != -1 {
					ns = nst[:idx]
				}
			}
			pP.Services.HRef = ns
			parsedPorts = append(parsedPorts, pP)
		}

		pParsed := &PropertyParsed{
			Hostnames : hostnames,
			OsMatches : osMatches,
			Time      : time.Now(),
			Address   : addr,
			Ports     : parsedPorts,
		}

		propertyParsed = append(propertyParsed, pParsed)
	}

	var res []byte 
	if res, err = json.Marshal(propertyParsed); err != nil {
		return "", err
	}
	return string(res), nil
}