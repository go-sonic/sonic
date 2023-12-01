package config

type Config struct {
	Server     Server      `mapstructure:"server"`
	Log        Log         `mapstructure:"logging"`
	PostgreSQL *PostgreSQL `mapstructure:"postgre"`
	MySQL      *MySQL      `mapstructure:"mysql"`
	SQLite3    *SQLite3    `mapstructure:"sqlite3"`
	Sonic      Sonic       `mapstructure:"sonic"`
}

type PostgreSQL struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	DB       string `mapstructure:"db"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

type MySQL struct {
	Dsn string `mapstructure:"dsn"`
}

type SQLite3 struct {
	Enable bool `mapstructure:"enable"`
	File   string
}
type Server struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}
type Log struct {
	FileName string `mapstructure:"filename"`
	Levels   Levels `mapstructure:"level"`
	MaxSize  int    `mapstructure:"maxsize"`
	MaxAge   int    `mapstructure:"maxage"`
	Compress bool   `mapstructure:"compress"`
}
type Levels struct {
	App  string `mapstructure:"app"`
	Gorm string `mapstructure:"gorm"`
}

type LogMode string

const (
	Console LogMode = "console"
	File    LogMode = "file"
)

type Sonic struct {
	Mode              string  `mapstructure:"mode"`
	LogMode           LogMode `mapstructure:"log_mode"`
	WorkDir           string  `mapstructure:"work_dir"`
	UploadDir         string
	LogDir            string `mapstructure:"log_dir"`
	TemplateDir       string `mapstructure:"template_dir"`
	ThemeDir          string
	AdminResourcesDir string
	AdminURLPath      string `mapstructure:"admin_url_path"`
}
