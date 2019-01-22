# workstep
步骤器

## 插件机制
通过插件扩展功能
### 已经实现插件
* ftp
* copy
* delete
* zip
* unzip
* bash

## 实现一个插件
```golang
package xxx

import (
	"github.com/Fengxq2014/workstep"
)

func Register(session *workstep.Session) {
	// 将插件注册到session
	// 第一个参数为具体方法
	// 第二个参数为插件名称
	session.HandlerRegister.Add(workstep.Handler(doxxx), "xxx")
}

func doxxx(s *workstep.Session) error {
	// do somthing
	s.Args 可用拿到配置文件中对args
	return nil
}
```

## 插件使用
### ftp
```json
{
    "type": "ftp",
    "args": "addr=10.39.12.56:21;user=ftpuser;password=passw0rd!;path=/wxpt_imgs/service/bkjdcx.png;des=/Users/defned/1-{datetime}.png"
  }
```
参数 | 说明 | 其他
---- | --- | -----
addr | ftp地址  | ip:端口
user | ftp用户  | 
password | ftp密码 |
path | ftp文件地址 |
des  | 下载地址 | 路径中可使用{date},{datetime}
参数之间用分号分割

### copy
```json
{
    "type": "copy",
    "args": "/Users/defned/Desktop/WechatIMG6.png /Users/defned/Desktop/WechatIMG688888-{datetime}.png"
  }

```
源文件地址，目标文件地址，中间用空格分割

### delete
```json
{
    "type": "delete",
    "args": "/Users/defned/git/workstep/test*"
  }
```

### zip
```json
{
    "type": "zip",
    "args": "files=/Users/defned/Desktop/WechatIMG6.png,/Users/defned/Desktop/WechatIMG7.png,/Users/defned/git/authz;des=/Users/defned/git/workstep/test.zip"
  }
```
参数 | 说明 | 其他
---- | --- | -----
files | 要压缩的文件  | 支持文件和文件夹，多个目标之间用逗号分割
des | 压缩文件路径  | 生成的压缩文件存放路径，根据文件后缀生成压缩文件,支持：zip、tar、tar.gz、rar、tar.xz等
参数之间用分号分割

### unzip
```json
{
    "type": "unzip",
    "args": "file=/Users/defned/git/workstep/test.zip;des=/Users/defned/git/workstep/testunzip"
  }
```
参数 | 说明 | 其他
---- | --- | -----
file | 要解缩的文件  | 支持多种压缩文件格式
des | 解压路径  | 
参数之间用分号分割

### bash
```json
{
    "type": "bash",
    "args": "jps"
  }
```
args为要执行等命令，支持windows cmd中可执行的命令

对nohup进行了特殊处理，命令不需要&

## 步骤配置文件
见项目step.json

## 使用
```bash
// 帮助
workstep -h
```
* -c 配置文件路径
* -ec bool，step出错是否继续
```bash
workstep -c /User/user/step.json -ec true
```

## 问题
* 为什么不直接用shell脚本
> shell脚本不能跨平台，windows/linux需要各自实现
* windows路径报错
> windos路径分隔符为\\, 请使用\\\\,既双反斜杠