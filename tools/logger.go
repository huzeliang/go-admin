package tools

import (
	"errors"
	log "github.com/sirupsen/logrus"
	config2 "go-admin/tools/config"
	"os"
	"time"
)

func InitLogger() {
	log.SetFormatter(&log.TextFormatter{FieldMap: log.FieldMap{
		log.FieldKeyTime:  "@timestamp",
		log.FieldKeyLevel: "@level",
		log.FieldKeyMsg:   "@message"}, TimestampFormat: "2006-01-02 15:04:05"})

	switch Mode(config2.ApplicationConfig.Mode) {
	case ModeDev, ModeTest:
		log.SetOutput(os.Stdout)
		log.SetLevel(log.TraceLevel)
	case ModeProd:
		// 创建日志目录
		path := config2.LogConfig.Dir
		_, err := os.Stat(path)
		if err != nil && os.IsNotExist(err) {
			_ = os.MkdirAll(path, os.ModePerm) //创建多级目录
		}

		file, err := os.OpenFile(path+"/api-"+time.Now().Format("2006-01-02")+".log", os.O_WRONLY|os.O_APPEND|os.O_CREATE|os.O_SYNC, 0600)
		if err != nil {
			log.Fatalln("log init failed")
		}

		var info os.FileInfo
		info, err = file.Stat()
		if err != nil {
			log.Fatal(err)
		}
		fileWriter := logFileWriter{file, info.Size()}
		log.SetOutput(&fileWriter)
		log.SetLevel(log.WarnLevel)
	}

	log.SetReportCaller(true)
}

type logFileWriter struct {
	file *os.File
	size int64
}

func (p *logFileWriter) Write(data []byte) (n int, err error) {
	if p == nil {
		return 0, errors.New("logFileWriter is nil")
	}
	if p.file == nil {
		return 0, errors.New("file not opened")
	}
	n, e := p.file.Write(data)
	p.size += int64(n)
	//每天一个文件
	if p.file.Name() != config2.LogConfig.Dir+"/api-"+time.Now().Format("2006-01-02")+".log" {
		p.file.Close()
		p.file, _ = os.OpenFile(config2.LogConfig.Dir+"/api-"+time.Now().Format("2006-01-02")+".log", os.O_WRONLY|os.O_APPEND|os.O_CREATE|os.O_SYNC, 0600)
		p.size = 0
	}
	return n, e
}
