package config

import (
	"encoding/json"

	"github.com/mylxsw/glacier/infra"
)

type Config struct {
	Listen           string `json:"listen"`
	Debug            bool   `json:"debug"`
	LogPath          string `json:"log_path"`
	Version          string `json:"version"`
	GitCommit        string `json:"git_commit"`
	DBConnStr        string `json:"-"`
	StoragePath      string `json:"storage_path"`
	SessionKey       string `json:"-"`
	LDAP             LDAP   `json:"ldap"`
	WeakPasswordMode bool   `json:"weak_password_mode"`
}

type LDAP struct {
	URL         string `json:"url"`
	BaseDN      string `json:"base_dn"`
	Username    string `json:"username"`
	Password    string `json:"-"`
	DisplayName string `json:"display_name"`
	UID         string `json:"uid"`
}

func (conf *Config) Serialize() string {
	rs, _ := json.Marshal(conf)
	return string(rs)
}

// Get return config object from container
func Get(cc infra.Resolver) *Config {
	return cc.MustGet(&Config{}).(*Config)
}
