package sftp

import (
	"errors"
	"fmt"
	"github.com/Fengxq2014/workstep"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io"
	"os"
	"strings"
	"time"
)

// Register 将插件注册到session
func Register(session *workstep.Session) {
	session.HandlerRegister.Add(workstep.Handler(dosftp), "sftp")
}

func dosftp(s *workstep.Session) error {
	if s.Args == "" {
		return errors.New("args is empty")
	}
	maps := make(map[string]string)
	split := strings.Split(s.Args, ";")
	for _, spl := range split {
		temp := strings.Split(spl, "=")
		if _, ok := maps[temp[0]]; !ok {
			maps[temp[0]] = temp[1]
		}
	}
	params := []string{"addr", "user", "password", "path", "des", "methods"}
	for _, param := range params {
		if _, ok := maps[param]; !ok {
			return errors.New("cant found param " + param)
		}
	}

	sftpClient, err := connect(maps["user"], maps["password"], maps["addr"])
	if err != nil {
		return err
	}
	defer sftpClient.Close()

	if strings.EqualFold(maps["methods"], "get") {
		srcFile, err := sftpClient.Open(maps["path"])
		if err != nil {
			return err
		}
		defer srcFile.Close()
		//stat, err := os.Stat(maps["des"])
		//if err != nil {
		//	return err
		//}
		//if stat.IsDir() {
		//	srcFile
		//}
		dstFile, err := os.Create(workstep.FormatStr(maps["des"]))
		if err != nil {
			return err
		}
		defer dstFile.Close()

		_, err = srcFile.WriteTo(dstFile)
		return err
	}
	if strings.EqualFold(maps["methods"], "put") {
		srcFile, err := os.Open(maps["path"])
		if err != nil {
			return err
		}
		defer srcFile.Close()
		dstFile, err := sftpClient.Create(workstep.FormatStr(maps["des"]))
		if err != nil {
			return err
		}
		defer dstFile.Close()
		_, err = io.Copy(dstFile, srcFile)
		return err
	}
	return fmt.Errorf("not support methods:%s", maps["methods"])
}

func connect(user, password, host string) (*sftp.Client, error) {
	var (
		auth         []ssh.AuthMethod
		clientConfig *ssh.ClientConfig
		sshClient    *ssh.Client
		sftpClient   *sftp.Client
		err          error
	)
	// get auth method
	auth = make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(password))
	auth = append(auth, ssh.KeyboardInteractive(func(user, instruction string, questions []string, echos []bool) (answers []string, err error) {
		answers = make([]string, len(questions))
		// The second parameter is unused
		for n, _ := range questions {
			answers[n] = password
		}

		return answers, nil
	}))

	clientConfig = &ssh.ClientConfig{
		User:            user,
		Auth:            auth,
		Timeout:         10 * time.Second,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// connet to ssh
	if sshClient, err = ssh.Dial("tcp", host, clientConfig); err != nil {
		return nil, err
	}

	// create sftp client
	if sftpClient, err = sftp.NewClient(sshClient); err != nil {
		return nil, err
	}

	return sftpClient, nil
}
