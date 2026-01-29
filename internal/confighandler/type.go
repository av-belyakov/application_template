package confighandler

// ConfigApp конфигурационные настройки приложения
type ConfigApp struct {
	Common             CfgCommon
	LogDB              CfgWriteLogDB
	Service            CfgService
	WebServer          CfgWebServer
	AuthenticationData CfgAuthenticationData
}

// CfgCommon общие настройки
type CfgCommon struct {
	Logs []*LogSet
}

// Logs настройки логирования
type Logs struct {
	Logging []*LogSet
}

type LogSet struct {
	MsgTypeName   string `validate:"oneof=error info warning" yaml:"msgTypeName"`
	PathDirectory string `validate:"required" yaml:"pathDirectory"`
	MaxFileSize   int    `validate:"min=1000" yaml:"maxFileSize"`
	WritingStdout bool   `validate:"required" yaml:"writingStdout"`
	WritingFile   bool   `validate:"required" yaml:"writingFile"`
	WritingDB     bool   `validate:"required" yaml:"writingDB"`
}

// CfgWriteLogDB настройки записи данных в БД
type CfgWriteLogDB struct {
	Host          string `yaml:"host"`
	User          string `yaml:"user"`
	NameDB        string `yaml:"namedb"`
	StorageNameDB string `yaml:"storage_name_db"`
	Port          int    `validate:"gt=0,lte=65535" yaml:"port"`
}

// CfgService настройки доступа к некоторому сервису
type CfgService struct {
	Host string `validate:"required" yaml:"host"`
	User string `validate:"required" yaml:"user"`
	Port int    `validate:"gt=0,lte=65535" yaml:"port"`
}

// CfgWebServer настройки веб сервера
type CfgWebServer struct {
	Host     string `validate:"required" yaml:"host"`
	Port     int    `validate:"gt=0,lte=65535" yaml:"port"`
	IsActive bool   `yaml:"isActive"`
}

type CfgAuthenticationData struct {
	SomeToken        string `validate:"required"`
	ServicePasswd    string `validate:"required"`
	WriteLogBDPasswd string `yaml:"passwd"`
}
