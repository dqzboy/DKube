package service

import (
	"dkube/config"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/wonderivan/logger"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
	"log"
	"net/http"
	"time"
)

var Terminal terminal

type terminal struct{}

func (t *terminal) WsHandler(w http.ResponseWriter, r *http.Request) {
	conf, err := clientcmd.BuildConfigFromFlags("", config.Kubeconfig)
	if err != nil {
		logger.Error("加载k8s配置失败, " + err.Error())
		return
	}
	if err := r.ParseForm(); err != nil {
		logger.Error("解析参数失败, " + err.Error())
		return
	}
	namespace := r.Form.Get("namespace")
	podName := r.Form.Get("pod_name")
	containerName := r.Form.Get("container_name")
	logger.Info("exec pod: %s, container: %s, namespace: %s\n", podName, containerName, namespace)

	pty, err := NewTerminalSession(w, r, nil)
	if err != nil {
		logger.Error("实例化TerminalSession失败, " + err.Error())
		return
	}
	defer func() {
		logger.Info("关闭TerminalSession")
		pty.Close()
	}()
	req := K8s.ClientSet.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(podName).
		Namespace(namespace).
		SubResource("exec").
		VersionedParams(&v1.PodExecOptions{
			Stdin:     true,
			Stdout:    true,
			Stderr:    true,
			TTY:       true,
			Container: containerName,
			Command:   []string{"/bin/bash"},
		}, scheme.ParameterCodec)
	logger.Info("exec post request url: ", req)

	executor, err := remotecommand.NewSPDYExecutor(conf, "POST", req.URL())
	if err != nil {
		logger.Error("建立SPDY连接失败, " + err.Error())
		return
	}
	err = executor.Stream(remotecommand.StreamOptions{
		Stdin:             pty,
		Stdout:            pty,
		Stderr:            pty,
		Tty:               true,
		TerminalSizeQueue: pty,
	})

	if err != nil {
		logger.Error("执行 pod 命令失败, " + err.Error())
		pty.Write([]byte("执行 pod 命令失败, " + err.Error()))
		pty.Done()
	}
}

type terminalMessage struct {
	Operation string `json:"operation"`
	Data      string `json:"data"`
	Rows      uint16 `json:"rows"`
	Cols      uint16 `json:"cols"`
}

type TerminalSession struct {
	wsConn   *websocket.Conn
	sizeChan chan remotecommand.TerminalSize
	doneChan chan struct{}
}

var upgrader = func() websocket.Upgrader {
	upgrader := websocket.Upgrader{}
	upgrader.HandshakeTimeout = time.Second * 2
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	return upgrader
}()

func NewTerminalSession(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (*TerminalSession, error) {
	conn, err := upgrader.Upgrade(w, r, responseHeader)
	if err != nil {
		return nil, errors.New("升级websocket失败," + err.Error())
	}
	session := &TerminalSession{
		wsConn:   conn,
		sizeChan: make(chan remotecommand.TerminalSize),
		doneChan: make(chan struct{}),
	}

	return session, nil
}

func (t *TerminalSession) Read(p []byte) (int, error) {
	_, message, err := t.wsConn.ReadMessage()
	if err != nil {
		log.Printf("read message err: %v", err)
		return 0, err
	}
	var msg terminalMessage
	if err := json.Unmarshal(message, &msg); err != nil {
		log.Printf("read parse message err: %v", err)
		return 0, err
	}
	switch msg.Operation {
	case "stdin":
		return copy(p, msg.Data), nil
	case "resize":
		t.sizeChan <- remotecommand.TerminalSize{Width: msg.Cols, Height: msg.Rows}
		return 0, nil
	case "ping":
		return 0, nil
	default:
		log.Printf("unknown message type '%s'", msg.Operation)
		return 0, fmt.Errorf("unknown message type '%s'", msg.Operation)
	}
}

func (t *TerminalSession) Write(p []byte) (int, error) {
	msg, err := json.Marshal(terminalMessage{
		Operation: "stdout",
		Data:      string(p),
	})
	if err != nil {
		log.Printf("write parse message err: %v", err)
		return 0, err
	}
	if err := t.wsConn.WriteMessage(websocket.TextMessage, msg); err != nil {
		log.Printf("write message err: %v", err)
		return 0, err
	}

	return len(p), nil
}

func (t *TerminalSession) Done() {
	close(t.doneChan)
}

func (t *TerminalSession) Close() {
	t.wsConn.Close()
}

func (t *TerminalSession) Next() *remotecommand.TerminalSize {
	select {
	case size := <-t.sizeChan:
		return &size
	case <-t.doneChan:
		return nil
	}
}
