package site

import (
    "fmt"
	"sync"
)

type SiteSettings struct {
    Domain         string
    AutobuyFlat    string
    FreeUser       bool
    FreeUserRoles  string
    FreeUserCons   int
    FreeUserTime   int
    FakeUsers      int
    FakeOngoing    int
    CNC            bool
    SSH            string
    Telnet         string
	API			   bool
	Subnet1       string   
	Subnet2       string  
	Subnet3       string   
	Subnet4       string   
	Subnet5       string   
	Subnet6       string   
	Subnet7       string   
	Subnet8       string  
	Subnet9       string   
}

var AllowedFields = map[string]bool{
    "domain":         true,
    "autobuy_flat":   true,
    "freeuser":       true,
    "freeuser_roles": true,
    "freeuser_cons":  true,
    "freeuser_time":  true,
    "fake_users":     true,
    "fake_ongoing":   true,
    "cnc":            true,
    "ssh":            true,
    "telnet":         true,
    "api":            true,      // ← new
    // subnets:
    "subnet1":        true,
    "subnet2":        true,
    "subnet3":        true,
    "subnet4":        true,
    "subnet5":        true,
    "subnet6":        true,
    "subnet7":        true,
    "subnet8":        true,
    "subnet9":        true,
}

var (
    fakeOngoingCurrent int
    fakeOngoingMu      sync.Mutex
	Site 			   *SiteSettings
)

func (conn *Instance) LoadSiteSettings() error {
	query := `
	  SELECT
		domain, autobuy_flat, freeuser, freeuser_roles,
		freeuser_cons, freeuser_time, fake_users, fake_ongoing,
		cnc, ssh, telnet, api,
		subnet1, subnet2, subnet3, subnet4, subnet5,
		subnet6, subnet7, subnet8, subnet9
	  FROM settings_conf
	  LIMIT 1
	`
	row := conn.conn.QueryRow(query)
  
	var s SiteSettings
	err := row.Scan(
	  &s.Domain, &s.AutobuyFlat, &s.FreeUser, &s.FreeUserRoles,
	  &s.FreeUserCons, &s.FreeUserTime, &s.FakeUsers, &s.FakeOngoing,
	  &s.CNC, &s.SSH, &s.Telnet, &s.API,
	  &s.Subnet1, &s.Subnet2, &s.Subnet3, &s.Subnet4, &s.Subnet5,
	  &s.Subnet6, &s.Subnet7, &s.Subnet8, &s.Subnet9,
	)
	if err != nil {
	  return err
	}
	Site = &s
    fakeOngoingMu.Lock()
    fakeOngoingCurrent = s.FakeOngoing
    fakeOngoingMu.Unlock()

	return nil
  }


func (conn *Instance) UpdateSiteSettingField(field string, value any) error {
    if !AllowedFields[field] {
		return fmt.Errorf("invalid field name: %s", field)
	}
	query := fmt.Sprintf("UPDATE settings_conf SET %s = ?", field)

	_, err := conn.conn.Exec(query, value)
	if err != nil {
		return err
	}

	return conn.LoadSiteSettings()
}