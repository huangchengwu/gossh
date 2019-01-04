package server

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

var (
	err          error
	client       *ssh.Client
	session      *ssh.Session
	oldState     *terminal.State
	data         map[string]string
	stdout       io.Reader
	stdin        io.WriteCloser
	p            []byte
	c            *websocket.Conn
	Httpserver   string
	Socketserver string
	Logdir       string
	Wlan         string
	Lan          string
)

// func Index(w http.ResponseWriter, r *http.Request) {
// 	http.Redirect(w, r, "http://"+host+"/home/www", http.StatusFound)
// }

func Get_log() *log.Logger {
	logFile, err := os.OpenFile(Logdir+time.Now().Format("20060102")+".log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {

		fmt.Println("错误无法创建日志")
	}
	loger := log.New(logFile, "", log.Ltime|log.Lshortfile|log.Ldate)
	return loger
}

func SSHLogin(user, pwd, host, height, width string, ch chan bool) {
	l := Get_log()

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(pwd),
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}
	//ssh tcp握手
	client, err = ssh.Dial("tcp", host, config)

	if err != nil {
		l.SetPrefix("[TCP]")
		l.Printf("\nuser:%s\npwd:%s\nhost:%s\n握手失败 %s", user, pwd, host, err)

		ch <- false
		return
	}
	//开启远程会话
	session, err = client.NewSession()
	if err != nil {
		l.SetPrefix("[SSH]")
		l.Println("远程开启回话失败", err)

		ch <- false
		return
	} else {
		ch <- true
	}
	defer session.Close()

	fd := int(os.Stdin.Fd())
	oldState, _ = terminal.MakeRaw(fd)
	if err != nil {
		log.SetPrefix("[SSH]")
		log.Println(err)
		ch <- false
		return
	}
	defer terminal.Restore(fd, oldState)

	stdout, _ = session.StdoutPipe()

	stdin, _ = session.StdinPipe()

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}
	h, _ := strconv.Atoi(height)
	w, _ := strconv.Atoi(width)
	session.RequestPty("xterm-256color", h, w, modes)
	session.Run("/bin/bash")

}

func HTTP(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	w.Header().Set("content-type", "application/json")

	fmt.Println(req.Method)
	body, _ := ioutil.ReadAll(req.Body)

	json.Unmarshal(body, &data)
	ch := make(chan bool, 1)
	l := Get_log()

	go SSHLogin(data["user"], data["pwd"], data["host"], data["height"], data["width"], ch)

	switch <-ch {
	case true:

		l.SetPrefix("[Login]")
		l.Println(data)

		data["status"] = "true"
		data["sock"] = "ws://" + Wlan + ":3002/socket"
		d, _ := json.Marshal(data)

		fmt.Fprint(w, string(d))
		return

	case false:
		data["status"] = "false"

		l.SetPrefix("[Login]")
		l.Println("登录失败", data)
		d, _ := json.Marshal(data)
		fmt.Fprint(w, string(d))
		return
	}

}

var upgrader = websocket.Upgrader{
	// 解决跨域问题
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Socket(w http.ResponseWriter, req *http.Request) {
	l := Get_log()

	c, err = upgrader.Upgrade(w, req, nil)
	defer c.Close()

	if err != nil {
		l.SetPrefix("[SSH]")
		l.Println("握手失败", data, err)
		return
	} else {
		l.SetPrefix("[SSH]")
		l.Println("握手成功", data, err)
	}
	go func() {
		for {
			//读取客户端数据
			_, p, err = c.ReadMessage()
			if err != nil {

				l.SetPrefix("[READ]")
				l.Println("读取失败退出登录", err)
				stdin.Write([]byte("exit\n"))
				return
			}
			go io.Copy(stdin, strings.NewReader(string(p)))

		}
	}()
	for {
		buf := make([]byte, 4096)
		n1, err := stdout.Read(buf)
		if err != nil {
			return
		}

		data["result"] = string(buf[:n1])

		//发送客户端数据
		err = c.WriteJSON(data)
		if err != nil {
			l.SetPrefix("[READ]")
			l.Println("发送失败", err)

		}
	}

}

//启动http
func Starthttp() {
	fmt.Printf("-------启动ssh http %s\n", Lan+":3001")
	http.HandleFunc("/login", HTTP)
	http.ListenAndServe(Lan+":3001", nil)
}

func Startsocket() {
	fmt.Printf("-------启动ssh socket %s\n", Lan+":3002")
	http.HandleFunc("/socket", Socket)
	http.ListenAndServe(Lan+":3002", nil)
}
