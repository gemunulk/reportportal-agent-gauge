package listener

import (
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"reflect"

	"strconv"
	"time"

	"reportportal-agent-gauge/event"

	"github.com/gorilla/websocket"
)

const flashServerPort = "FLASH_SERVER_PORT"

var port = fmt.Sprintf(":%d", getPort())

var connections []*websocket.Conn

var events []event.Event

var homeTemplate *template.Template

var timestamp = time.Now().Format("2006-01-02 15:04:05")

func home(w http.ResponseWriter, r *http.Request) {
	data := struct {
		URL       string
		Project   string
		Timestamp string
	}{
		URL:       fmt.Sprintf("127.0.0.1%s/progress", port),
		Project:   GetProjectRoot(),
		Timestamp: timestamp,
	}
	homeTemplate.Execute(w, data)
}

func GetProjectRoot() string {
	projectRoot := os.Getenv("GAUGE_PROJECT_ROOT")
	if projectRoot == "" {
		return "Sample"
	}
	return filepath.Base(projectRoot)
}

func progress(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Upgrade(w, r, w.Header(), 1024, 1024)
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
	}
	connections = append(connections, conn)
	for _, e := range events {
		conn.WriteJSON(e)
	}
}

func handleEvents(e chan event.Event) {
	for {
		et := <-e
		events = append(events, et)
		endEvent := event.NewEndEvent(false)
		if reflect.TypeOf(et) == reflect.TypeOf(endEvent) {
			for _, c := range connections {
				c.WriteJSON(et)
				c.Close()
			}
			e <- endEvent
			return
		}
		for _, c := range connections {
			c.WriteJSON(et)
		}
	}
}

type httpListener struct {
	event chan event.Event
}

func NewHttpListener(e chan event.Event) Listener {
	return &httpListener{event: e}
}

func (l *httpListener) Start() {
	http.HandleFunc("/", home)
	http.HandleFunc("/progress", progress)
	go handleEvents(l.event)
	fmt.Println("[ReportPortal] Starting ...")
	var err error
	homeTemplate, err = template.New("home").Parse(html)
	if err != nil {
		log.Fatalf(err.Error())
	}
	log.Fatal(http.ListenAndServe(port, nil))
}

func getPort() int {
	port := os.Getenv(flashServerPort)
	if port != "" {
		p, err := strconv.Atoi(port)
		if err == nil {
			return p
		}
		fmt.Printf("[ReportPortal] Cannot use %s='%s' value. Error: %s\n", flashServerPort, port, err)
	}
	l, err := net.ListenTCP("tcp", &net.TCPAddr{Port: 0})
	if err != nil {
		log.Fatalf(err.Error())
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port
}
