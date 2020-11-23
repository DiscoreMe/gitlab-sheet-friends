package config

import (
	"fmt"
	"github.com/brianvoe/gofakeit"
	"github.com/spf13/viper"
	"os"
	"sort"
	"time"
)

const (
	errRequiredField = "error %s: type is the required field"
)

type Config struct {
	DB              string    `mapstructure:"db"`
	Gits            []Git     `mapstructure:"-"`
	SpreadSheetID   string    `mapstructure:"spreadsheet_id"`
	Members         []Member  `mapstructure:"-"`
	StartingTimeInt int64     `mapstructure:"starting_time"`
	StartingTime    time.Time `mapstructure:"-"`
	TmplSheetID     int64     `mapstructure:"tmpl_sheet_id"`
}

type Member struct {
	ID       string
	Name     string
	Services map[string]string
}

type Git struct {
	Name      string
	Type      string
	URL       string
	Available string
	Token     string
}

func newGit() Git {
	return Git{
		Name:      gofakeit.Company(),
		Available: "internal",
	}
}

func newMember(id, name string) Member {
	return Member{
		ID:       id,
		Name:     name,
		Services: make(map[string]string),
	}
}

func setDefaultValues() {
	viper.SetDefault("db", "db.db")
	viper.SetDefault("sheet_id", "")
	viper.SetDefault("starting_time", 0)
	viper.SetDefault("tmpl_sheet_id", -1)
}

// LoadConfig loads config from a file
func LoadConfig(filename string) (*Config, error) {
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
	setDefaultValues()

	fileConf, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer fileConf.Close()

	if err := viper.ReadConfig(fileConf); err != nil {
		return nil, err
	}

	conf := &Config{}
	if err := viper.Unmarshal(conf); err != nil {
		return nil, err
	}

	services := viper.GetStringMap("services")
	for key, value := range services {
		git := newGit()
		if key != "" {
			git.Name = key
		}
		serv := value.(map[string]interface{})
		if t, ok := serv["type"]; ok {
			git.Type = t.(string)
		} else {
			return nil, fmt.Errorf(errRequiredField, key)
		}
		if url, ok := serv["url"]; ok {
			git.URL = url.(string)
		} else {
			return nil, fmt.Errorf(errRequiredField, key)
		}
		if token, ok := serv["token"]; ok {
			git.Token = token.(string)
		} else {
			return nil, fmt.Errorf(errRequiredField, key)
		}
		if av, ok := serv["available"]; ok {
			git.Available = av.(string)
		}

		conf.Gits = append(conf.Gits, git)
	}

	members := viper.GetStringMap("members")
	for key, value := range members {
		memberMap := value.(map[string]interface{})
		member := newMember(key, memberMap["name"].(string))
		if services, ok := memberMap["services"].([]interface{}); ok {
			for _, vs := range services {
				for kvs, vvs := range vs.(map[interface{}]interface{}) {
					member.Services[kvs.(string)] = vvs.(string)
				}
			}
		}
		conf.Members = append(conf.Members, member)
		sort.Slice(conf.Members, func(i, j int) bool {
			return conf.Members[i].Name < conf.Members[j].Name
		})
	}

	conf.StartingTime = time.Unix(conf.StartingTimeInt, 0)

	return conf, nil
}
