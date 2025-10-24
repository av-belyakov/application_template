package confighandler

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"

	"github.com/av-belyakov/enricher_zin/internal/supportingfunctions"
)

func New(rootDir string) (*ConfigApp, error) {
	conf := &ConfigApp{}

	var (
		validate *validator.Validate
		envList  map[string]string = map[string]string{
			"GO_ENRICHERZIN_MAIN": "",

			//Подключение к Zabbix
			"GO_ENRICHERZIN_ZHOST":   "",
			"GO_ENRICHERZIN_ZPORT":   "",
			"GO_ENRICHERZIN_ZUSER":   "",
			"GO_ENRICHERZIN_ZPASSWD": "",

			//Подключение к NetBox
			"GO_ENRICHERZIN_NBHOST":   "",
			"GO_ENRICHERZIN_NBPORT":   "",
			"GO_ENRICHERZIN_NBUSER":   "",
			"GO_ENRICHERZIN_NBPASSWD": "",

			//Настройки доступа к БД в которую будут записыватся логи
			"GO_ENRICHERZIN_DBWLOGHOST":        "",
			"GO_ENRICHERZIN_DBWLOGPORT":        "",
			"GO_ENRICHERZIN_DBWLOGNAME":        "",
			"GO_ENRICHERZIN_DBWLOGUSER":        "",
			"GO_ENRICHERZIN_DBWLOGPASSWD":      "",
			"GO_ENRICHERZIN_DBWLOGSTORAGENAME": "",
		}
	)

	getFileName := func(sf, confPath string, lfs []fs.DirEntry) (string, error) {
		for _, v := range lfs {
			if v.Name() == sf && !v.IsDir() {
				return filepath.Join(confPath, v.Name()), nil
			}
		}

		return "", fmt.Errorf("file '%s' is not found", sf)
	}

	setCommonSettings := func(fn string) error {
		viper.SetConfigFile(fn)
		viper.SetConfigType("yml")
		if err := viper.ReadInConfig(); err != nil {
			return err
		}

		ls := Logs{}
		if ok := viper.IsSet("LOGGING"); ok {
			if err := viper.GetViper().Unmarshal(&ls); err != nil {
				return err
			}

			conf.Common.Logs = ls.Logging
		}

		return nil
	}

	setSpecial := func(fn string) error {
		viper.SetConfigFile(fn)
		viper.SetConfigType("yml")
		if err := viper.ReadInConfig(); err != nil {
			return err
		}

		//Настройки для модуля подключения к Zabbix
		if viper.IsSet("Zabbix.host") {
			conf.Zabbix.Host = viper.GetString("Zabbix.host")
		}
		if viper.IsSet("Zabbix.port") {
			conf.Zabbix.Port = viper.GetInt("Zabbix.port")
		}
		if viper.IsSet("Zabbix.user") {
			conf.Zabbix.User = viper.GetString("Zabbix.user")
		}

		// Настройки доступа к NetBox
		if viper.IsSet("NetBox.host") {
			conf.NetBox.Host = viper.GetString("NetBox.host")
		}
		if viper.IsSet("NetBox.port") {
			conf.NetBox.Port = viper.GetInt("NetBox.port")
		}
		if viper.IsSet("NetBox.user") {
			conf.NetBox.User = viper.GetString("NetBox.user")
		}

		// Настройки доступа к БД в которую будут записыватся логи
		if viper.IsSet("WriteLogDataBase.host") {
			conf.LogDB.Host = viper.GetString("WriteLogDataBase.host")
		}
		if viper.IsSet("WriteLogDataBase.port") {
			conf.LogDB.Port = viper.GetInt("WriteLogDataBase.port")
		}
		if viper.IsSet("WriteLogDataBase.user") {
			conf.LogDB.User = viper.GetString("WriteLogDataBase.user")
		}
		if viper.IsSet("WriteLogDataBase.namedb") {
			conf.LogDB.NameDB = viper.GetString("WriteLogDataBase.namedb")
		}
		if viper.IsSet("WriteLogDataBase.storage_name_db") {
			conf.LogDB.StorageNameDB = viper.GetString("WriteLogDataBase.storage_name_db")
		}

		return nil
	}

	validate = validator.New(validator.WithRequiredStructEnabled())

	for v := range envList {
		if env, ok := os.LookupEnv(v); ok {
			envList[v] = env
		}
	}

	rootPath, err := supportingfunctions.GetRootPath(rootDir)
	if err != nil {
		return conf, err
	}

	confPath := filepath.Join(rootPath, "config")
	list, err := os.ReadDir(confPath)
	if err != nil {
		return conf, err
	}

	fileNameCommon, err := getFileName("config.yml", confPath, list)
	if err != nil {
		return conf, err
	}

	//читаем общий конфигурационный файл
	if err := setCommonSettings(fileNameCommon); err != nil {
		return conf, err
	}

	var fn string
	if envList["GO_ENRICHERZIN_MAIN"] == "development" {
		fn, err = getFileName("config_dev.yml", confPath, list)
		if err != nil {
			return conf, err
		}
	} else if envList["GO_ENRICHERZIN_MAIN"] == "test" {
		fn, err = getFileName("config_test.yml", confPath, list)
		if err != nil {
			return conf, err
		}
	} else {
		fn, err = getFileName("config_prod.yml", confPath, list)
		if err != nil {
			return conf, err
		}
	}

	if err := setSpecial(fn); err != nil {
		return conf, err
	}

	//Настройки для модуля подключения к Zabbix
	if envList["GO_ENRICHERZIN_ZHOST"] != "" {
		conf.Zabbix.Host = envList["GO_ENRICHERZIN_ZHOST"]
	}
	if envList["GO_ENRICHERZIN_ZPORT"] != "" {
		if p, err := strconv.Atoi(envList["GO_ENRICHERZIN_ZPORT"]); err == nil {
			conf.Zabbix.Port = p
		}
	}
	if envList["GO_ENRICHERZIN_ZUSER"] != "" {
		conf.Zabbix.User = envList["GO_ENRICHERZIN_ZUSER"]
	}
	if envList["GO_ENRICHERZIN_ZPASSWD"] != "" {
		conf.Zabbix.Passwd = envList["GO_ENRICHERZIN_ZPASSWD"]
	}

	//Подключение к NetBox
	if envList["GO_ENRICHERZIN_NBHOST"] != "" {
		conf.NetBox.Host = envList["GO_ENRICHERZIN_NBHOST"]
	}
	if envList["GO_ENRICHERZIN_NBPORT"] != "" {
		if p, err := strconv.Atoi(envList["GO_ENRICHERZIN_NBPORT"]); err == nil {
			conf.NetBox.Port = p
		}
	}
	if envList["GO_ENRICHERZIN_NBUSER"] != "" {
		conf.NetBox.User = envList["GO_ENRICHERZIN_NBUSER"]
	}
	if envList["GO_ENRICHERZIN_NBPASSWD"] != "" {
		conf.NetBox.Passwd = envList["GO_ENRICHERZIN_NBPASSWD"]
	}

	//Настройки доступа к БД в которую будут записыватся логи
	if envList["GO_ENRICHERZIN_DBWLOGHOST"] != "" {
		conf.LogDB.Host = envList["GO_ENRICHERZIN_DBWLOGHOST"]
	}
	if envList["GO_ENRICHERZIN_DBWLOGPORT"] != "" {
		if p, err := strconv.Atoi(envList["GO_ENRICHERZIN_DBWLOGPORT"]); err == nil {
			conf.LogDB.Port = p
		}
	}
	if envList["GO_ENRICHERZIN_DBWLOGNAME"] != "" {
		conf.LogDB.NameDB = envList["GO_ENRICHERZIN_DBWLOGNAME"]
	}
	if envList["GO_ENRICHERZIN_DBWLOGUSER"] != "" {
		conf.LogDB.User = envList["GO_ENRICHERZIN_DBWLOGUSER"]
	}
	if envList["GO_ENRICHERZIN_DBWLOGPASSWD"] != "" {
		conf.LogDB.Passwd = envList["GO_ENRICHERZIN_DBWLOGPASSWD"]
	}
	if envList["GO_ENRICHERZIN_DBWLOGSTORAGENAME"] != "" {
		conf.LogDB.StorageNameDB = envList["GO_ENRICHERZIN_DBWLOGSTORAGENAME"]
	}

	//выполняем проверку заполненой структуры
	if err = validate.Struct(conf); err != nil {
		return conf, err
	}

	return conf, nil
}
