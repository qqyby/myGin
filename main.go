package main

/*
	配置文件
	数据库
	日志
*/


import (
	"context"
	"flag"
	"fmt"
	"go.uber.org/zap"
	"log"
	"math/rand"
	"myGin/inits"
	"myGin/internal/routers"
	"myGin/settings"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	cfgPath string
	isVersion    bool
	onlyVersion  bool
	buildTime    string
	buildVersion string
	gitCommitID  string
)

func init() {
	rand.Seed(time.Now().UnixNano())

	flag.StringVar(&cfgPath, "c", "", "配置文件路径:/xxx/xxx/config.yaml")
	flag.BoolVar(&isVersion, "v", false, "编译信息")
	flag.BoolVar(&onlyVersion, "version", false, "版本号")
	flag.Parse()
}

func main()  {
	flag.Parse()

	if isVersion {
		fmt.Printf("build_time: %s\n", buildTime)
		fmt.Printf("build_version: %s\n", buildVersion)
		fmt.Printf("git_commit_id: %s\n", gitCommitID)
		return
	}

	if onlyVersion {
		fmt.Println(buildVersion)
		return
	}

	// 初始化配置和日志
	inits.InitCfgAndLog(cfgPath)

	// 初始化数据库
	inits.InitDao()
	defer inits.Close()

	// todo 初始化系统服务

	// 注册路由
	r := routers.Setup()

	// 启动服务
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", settings.AppCfg.Port),
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// 等待中断信号来优雅地关闭服务器，为关闭服务器操作设置一个5秒的超时
	quit := make(chan os.Signal, 1) // 创建一个接收信号的通道
	// kill 默认会发送 syscall.SIGTERM 信号
	// kill -2 发送 syscall.SIGINT 信号，我们常用的Ctrl+C就是触发系统SIGINT信号
	// kill -9 发送 syscall.SIGKILL 信号，但是不能被捕获，所以不需要添加它
	// signal.Notify把收到的 syscall.SIGINT或syscall.SIGTERM 信号转发给quit
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // 此处不会阻塞
	<-quit                                               // 阻塞在此，当接收到上述两种信号时才会往下执行
	zap.L().Warn("Shutdown Server ...")
	// 创建一个5秒超时的context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// 5秒内优雅关闭服务（将未处理完的请求处理完再关闭服务），超过5秒就超时退出
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Fatal("Server Shutdown: ", zap.Error(err))
	}
	zap.L().Warn("Server exiting")


	fmt.Println(">>>> print")
	return
}
