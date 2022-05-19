package settings

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"myGin/pkg/utils"
	"os"
	"path"
	"path/filepath"
)

var (
	FFmpeg, FFprobe string
	LocalIpPort     string // 本机的ip:port
	VerifyLicense string // 验证license的工具
)

type appConfig struct {
	RootDir             string
	Name                string `mapstructure:"name"`
	RunMode             string `mapstructure:"run_mode"`
	Ip                  string `mapstructure:"Ip"`
	Port                int    `mapstructure:"port"`
	SnowflakeStartTime  string `mapstructure:"snowflake_start_time"`
	SnowflakeMachineID  int64  `mapstructure:"snowflake_machine_id"`
	OutputDir           string `mapstructure:"output_dir"`
	RequestTokenTimeout int64  `mapstructure:"request_token_timeout"`
	ApiSecretKey        string `mapstructure:"api_secret_key"`
	OssDownloadUrl      string `mapstructure:"oss_download_url"`
	OssVideoUploadUrl   string `mapstructure:"oss_video_upload_url"`
	LicenseFileDir      string `mapstructure:"license_file_dir"`
	NodeConcurrentJob   int64  `mapstructure:"node_concurrent_job"`
	SlaveOf             string `mapstructure:"slave_of"`
}

type logConfig struct {
	Level     string `mapstructure:"level"`
	Filename  string `mapstructure:"filename"`
	MaxSize   int    `mapstructure:"max_size"`
	MaxAge    int    `mapstructure:"max_age"`
	MaxBackup int    `mapstructure:"max_backup"`
}

type mysqlConfig struct {
	Host        string `mapstructure:"host"`
	Port        int    `mapstructure:"port"`
	User        string `mapstructure:"user"`
	Password    string `mapstructure:"password"`
	DbName      string `mapstructure:"dbname"`
	MaxOpenConn int    `mapstructure:"max_open_conn"`
	MaxIdleConn int    `mapstructure:"max_idle_conn"`
}

type redisConfig struct {
	MasterName string   `mapstructure:"master_name"`
	Host       []string `mapstructure:"host"`
	Password   string   `mapstructure:"password"`
	Db         int      `mapstructure:"db"`
	PoolSize   int      `mapstructure:"pool_size"`
}

var (
	AppCfg   = new(appConfig)
	LogCfg   = new(logConfig)
	MysqlCfg = new(mysqlConfig)
	RedisCfg = new(redisConfig)
	//LicenseCfg = NewTesterLicense()
)

func Init(cfgPath string) error {
	AppCfg.RootDir = inferRootDir()
	FFmpeg = path.Join(AppCfg.RootDir, "objs", "ffmpeg")
	FFprobe = path.Join(AppCfg.RootDir, "objs", "ffprobe")
	VerifyLicense = path.Join(AppCfg.RootDir, "objs", "verify_license")

	if cfgPath != "" {
		viper.SetConfigFile(cfgPath)
	} else {
		cfgPath = path.Join(AppCfg.RootDir, "configs", "config.yaml")
		viper.SetConfigFile(cfgPath)
	}

	if err := viper.ReadInConfig(); err != nil {
		return errors.Wrapf(err, "viper read config failed")
	}

	if err := readSection("app", AppCfg); err != nil {
		return errors.Wrapf(err, "read app config failed")
	}

	if err := readSection("log", LogCfg); err != nil {
		return errors.Wrapf(err, "read log config failed")
	}

	if err := readSection("mysql", MysqlCfg); err != nil {
		return errors.Wrapf(err, "read mysql config failed")
	}

	if err := readSection("redis", RedisCfg); err != nil {
		return errors.Wrapf(err, "read redis config failed")
	}

	LocalIpPort = fmt.Sprintf("%v:%v", AppCfg.Ip, AppCfg.Port)

	// 加载license
	if err := LoadLicense(); err != nil {
		fmt.Printf("load license failed:%v; use default", err)
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		if err := reloadAllSection(); err != nil {
			fmt.Printf("config file:%v changed, read section failed: %v\n ", in.Name, err)
		} else {
			fmt.Printf("config file:%v changed, read section success", in.Name)
		}
	})

	return nil

}

// inferRootDir 递归推导项目根目录
func inferRootDir() string {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	var infer func(d string) string
	infer = func(d string) string {
		if d == "/" {
			panic("请确保在项目根目录或子目录下运行程序，当前在：" + cwd)
		}

		if utils.Exist(d + "/configs") {
			return d
		}

		return infer(filepath.Dir(d))
	}

	return infer(cwd)
}

var sections = make(map[string]interface{})

func readSection(k string, v interface{}) error {
	err := viper.UnmarshalKey(k, v)
	if err != nil {
		return err
	}
	if _, ok := sections[k]; !ok {
		sections[k] = v
	}
	return nil
}

func reloadAllSection() error {
	for k, v := range sections {
		err := readSection(k, v)
		if err != nil {
			return err
		}
	}
	return nil
}
