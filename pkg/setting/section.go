package setting

import "time"

type ServerSettingS struct {
	Host         string
	Domain       string
	RunMode      string
	RootPath     string
	HttpPort     string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type OauthSettingS struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
	AuthURL      string
	TokenURL     string
	UserInfoURL  string
	LogoutURL    string
}

type DatabaseSettingS struct {
	DBType          string
	UserName        string
	Password        string
	Host            string
	DBName          string
	TablePrefix     string
	Charset         string
	ParseTime       bool
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
}
type MeetSettingS struct {
	Host     string
	RoBotKey string
}

type ThumbSettingS struct {
	MaxWidth  int
	MaxHeight int
	Quality   int
}

type DifySettingS struct {
	APIURL string
	APIKey string
}

func (s *Setting) ReadSection(k string, v interface{}) error {
	err := s.vp.UnmarshalKey(k, v)
	if err != nil {
		return err
	}
	return nil
}
