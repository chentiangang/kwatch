package xlog

import (
	"fmt"
	"os"
	"sync"
	"time"
)

// 创建一个把日志写入文件的对象
type XFile struct {
	filename string
	file     *os.File
	*XLogBase

	logChan chan *LogData
	wg      *sync.WaitGroup
	curDay int
}

// 给文件这个对象进行赋值
func NewXFile(level int, filename, module string) XLog {
	logger := &XFile{
		filename: filename,
	}
	logger.XLogBase = &XLogBase{
		level:  level,
		module: module,
	}

	logger.curDay = time.Now().Day()
	logger.wg = &sync.WaitGroup{}
	logger.logChan = make(chan *LogData, 10000)
	logger.wg.Add(1)
	go logger.syncLog()
	return logger
}

// 这个对像有了值以后，添加以下的方法来实现XLog这个接口。

func (c *XFile) Init() (err error) {
	c.file, err = os.OpenFile(c.filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	return
}

func (c *XFile) syncLog() {
	for data := range c.logChan {
			c.splitLog()
		c.writeLog(c.file, data)
	}
	c.wg.Done()
}

func (c *XFile) splitLog() {
	now := time.Now()

	fmt.Println(now.Day(),c.curDay)
	// 如果当前的小时等于创建文件对象时的小时时间，则返回，不进行日志切割
	if now.Day() == c.curDay {
		return
	}

	c.curDay = now.Day()
	// 把文件从缓冲区同步到磁盘
	c.file.Sync()

	// 关闭文件
	c.file.Close()

	// 新生成的文件名
	newFilename := fmt.Sprintf("%s-%04d-%02d-%02d-%02d", c.filename,
		now.Year(), now.Month(), now.Day(), now.Hour())

	// os.Rename 修改文件名
	os.Rename(c.filename, newFilename)

	// 初始化
	c.Init()
}

func (c *XFile) writeToChan(level int, module string, format string, args ...interface{}) {
	if c.level > level {
		return
	}
	LogData := c.formatLogger(level, module, format, args...)
	select {
	case c.logChan <- LogData:
	default:
	}
}

func (c *XFile) LogDebug(format string, args ...interface{}) {
	c.writeToChan(XLogLevelDebug, c.module, format, args...)
}

func (c *XFile) LogTrace(format string, args ...interface{}) {
	c.writeToChan(XLogLevelTrace, c.module, format, args...)
}

func (c *XFile) LogInfo(format string, args ...interface{}) {
	c.writeToChan(XLogLevelInfo, c.module, format, args...)
}

func (c *XFile) LogWarn(format string, args ...interface{}) {
	c.writeToChan(XLogLevelWarn, c.module, format, args...)
}

func (c *XFile) LogError(format string, args ...interface{}) {
	c.writeToChan(XLogLevelError, c.module, format, args...)
}

func (c *XFile) LogFatal(format string, args ...interface{}) {
	c.writeToChan(XLogLevelFatal, c.module, format, args...)
}

func (c *XFile) SetLevel(level int) {
	c.level = level
}

func (c *XFile) Close() {
	if c.logChan != nil {
		close(c.logChan)
	}
	c.wg.Wait()
	if c.file != nil {
		c.file.Sync()
		c.file.Close()
	}
}
