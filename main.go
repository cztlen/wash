package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"time"
	"wash/model"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/spf13/viper"
)

var (
	// Db        *gorm.DB
	ch        chan string
	wg        sync.WaitGroup
	goCount   = 5
	conf_path = "D:\\go\\bb\\wash\\conf\\field.yaml"
	tmp       = "2006-01-02 15:04:05"
)

type Account struct {
	Id          uint   `gorm:"primary_key;AUTO_INCREMENT"  json:"-"`
	Player_id   string `json:"#account_id"`
	Type        string `gorm:"-" json:"#user_set"`
	Time        string `gorm:"-" json:"#time"`
	Channel_id  string `gorm:"type:varchar(10)" json:"-"`
	Uid         string `gorm:"index:idx_uid_time;NOT NULL;type:varchar(100)"  json:"-"`
	Prop        *Prop  `json:"properties"`
	Active_time int    `gorm:"index:idx_uid_time;type:int(10);NOT NULL"  json:"-"`
}

type Prop struct {
	Active_time string `json:"active_time"`
	Channel_id  string `json:"channel_id"`
	Uid         string `json:"uid"`
}

func main() {

	ch = make(chan string, 1024)
	// //获取文件路径
	file_path := getPath(conf_path)

	files := getFiles(file_path)
	aLog := createLogFile()
	for i := 0; i < goCount; i++ {
		wg.Add(1)
		go aLog.readChanData()
	}
	readFilesData(files)

	fmt.Println("执行结束...")
	// wg.Wait()
}

//遍历文件，读取数据
func readFilesData(files []string) error {
	defer wg.Done()
	for _, file := range files {
		fmt.Println(file)
		getFileData(file)

		// return nil

	}
	return nil
}

//获取文件夹下所有文件
func getFiles(path string) []string {
	var files []string
	fileInfo, _ := ioutil.ReadDir(path)
	for _, f := range fileInfo {
		files = append(files, path+f.Name())
		// fmt.Println(path + f.Name())
	}
	return files
}

//获取存放文件路径 dir
func getPath(conf_path string) string {
	viper.SetConfigFile(conf_path)

	if err := viper.ReadInConfig(); err != nil {
		return err.Error()
	}
	path := viper.GetString("account_file_path")

	return path
}

//获取文件内容并放入channel
//path :E:\gamelog\sgzj2\android2_1\actionlog\action.20220710.log
func getFileData(name string) {

	file, err := os.Open(name)
	if err != nil {
		panic("文件读取异常")
	}
	defer file.Close()
	read := bufio.NewReader(file)
	for {
		lineBytes, _, err := read.ReadLine()
		if err == io.EOF {
			close(ch)
			// model.Db.Close()
			fmt.Println("通道已关闭")
			break
		}
		lineStr := string(lineBytes)
		// fmt.Println(lineStr)

		// break
		ch <- lineStr
	}

}

//将channel中数据写入到表中
func (alog *accountLog) readChanData() {
	defer wg.Done()
	for lineStr := range ch {
		a := Account{}
		lineSlice := strings.Split(lineStr, ",")
		a.Player_id = lineSlice[0]
		a.Channel_id = lineSlice[1]
		a.Uid = lineSlice[2]
		res, _ := time.ParseInLocation(tmp, lineSlice[3], time.Local)
		a.Active_time = int(res.Unix())
		a.Type = "user_set"
		a.Time = time.Now().Format(tmp)

		a.Prop = &Prop{}
		a.Prop.Active_time = res.Format(tmp)
		a.Prop.Channel_id = a.Channel_id
		a.Prop.Uid = a.Uid

		// fmt.Println(a)
		// a.Prop = &Prop{
		// 	Active_time: a.Active_time,
		// 	Channel_id:  a.Channel_id,
		// 	Uid:         a.Uid,
		// }
		// fmt.Println(a)
		data, err := json.Marshal(&a)
		if err != nil {
			fmt.Printf("序列号错误 err=%v\n", err)
		}
		// fmt.Println(lineStr)
		alog.writeDataToLog(string(data) + "\n")
		// fmt.Printf("%v\n", string(data))A
		// break

		model.Db.Create(&a)

	}

}

//生成日志数据文件
type accountLog struct {
	File     *os.File
	BasePath string
}

//创建日志文件
func createLogFile() *accountLog {
	aLog := &accountLog{}
	aLog.BasePath = "./log.user_set.log"
	aLog.File, _ = os.OpenFile(aLog.BasePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	return aLog
}

//生成一定的日志格式写入文件中 json格式
func (aLog *accountLog) writeDataToLog(data string) {
	aLog.File.WriteString(data)
}

//docker镜像 定义应用程序及其运行所需要的一切
//解析依赖 工作区里的每个Go Module 在解析依赖时都被当做根Module
//在1.18以前，module A 新增feature,module b 再删除module B的
//go.mod 文件里的replace指令

//有了go工作区模式之后，针对上述场景，我们有了更为简单的方案：可以在
//工作区维护一个go.work文件来管理你的所有依赖。go.work里的use和replace指令
//会覆盖工作区目录下的每个Go Module的go.mod文件，因此没有必要去修改Go Module的
//go.mod文件了。

//go work init [moddirs]
//moddirs 是go module所在的本地目录。如果有多个go module,就用空格分开。如果go work
//init 后面没有参数，会创建一个空的workspace
//go work use 新增go module
//go work use [-r] moddir
//go.work go/use（添加一个本地磁盘上的go module到workspace的主module）/replace
//几种使用场景机器最佳实践
//
//
