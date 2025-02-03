package Plugins

import (
	common2 "RLscan/pkg/common"
	"fmt"
	"github.com/jlaffaye/ftp"
	"strings"
	"time"
)

func FtpScan(info *common2.HostInfo) (tmperr error) {
	if common2.IsBrute {
		return
	}
	starttime := time.Now().Unix()
	flag, err := FtpConn(info, "anonymous", "")
	if flag && err == nil {
		return err
	} else {
		errlog := fmt.Sprintf("[-] ftp %v:%v %v %v", info.Host, info.Ports, "anonymous", err)
		common2.LogError(errlog)
		tmperr = err
		if common2.CheckErrs(err) {
			return err
		}
	}

	for _, user := range common2.Userdict["ftp"] {
		for _, pass := range common2.Passwords {
			pass = strings.Replace(pass, "{user}", user, -1)
			flag, err := FtpConn(info, user, pass)
			if flag && err == nil {
				return err
			} else {
				errlog := fmt.Sprintf("[-] ftp %v:%v %v %v %v", info.Host, info.Ports, user, pass, err)
				common2.LogError(errlog)
				tmperr = err
				if common2.CheckErrs(err) {
					return err
				}
				if time.Now().Unix()-starttime > (int64(len(common2.Userdict["ftp"])*len(common2.Passwords)) * common2.Timeout) {
					return err
				}
			}
		}
	}
	return tmperr
}

func FtpConn(info *common2.HostInfo, user string, pass string) (flag bool, err error) {
	flag = false
	Host, Port, Username, Password := info.Host, info.Ports, user, pass
	conn, err := ftp.DialTimeout(fmt.Sprintf("%v:%v", Host, Port), time.Duration(common2.Timeout)*time.Second)
	if err == nil {
		err = conn.Login(Username, Password)
		if err == nil {
			flag = true
			result := fmt.Sprintf("[+] ftp %v:%v:%v %v", Host, Port, Username, Password)
			dirs, err := conn.List("")
			//defer conn.Logout()
			if err == nil {
				if len(dirs) > 0 {
					for i := 0; i < len(dirs); i++ {
						if len(dirs[i].Name) > 50 {
							result += "\n   [->]" + dirs[i].Name[:50]
						} else {
							result += "\n   [->]" + dirs[i].Name
						}
						if i == 5 {
							break
						}
					}
				}
			}
			common2.LogSuccess(result)
		}
	}
	return flag, err
}
