package main

import (
	"encoding/json"
	"fmt"
	"html"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
	uuid "github.com/satori/go.uuid"
	melody "gopkg.in/olahol/melody.v1"
)

var socketMap = map[string]*connection{}
var m = melody.New()

// HTTP サーバー起動
func main() {

	router := httprouter.New()
	router.GET("/", viewHandler)
	router.ServeFiles("/static/*filepath", http.Dir("dist/static"))

	m.Config.MaxMessageSize = 1024 * 1024 * 10

	router.GET("/ws", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		m.HandleRequest(w, r)
	})

	m.HandleMessage(echo)
	m.HandleConnect(connect)
	m.HandleDisconnect(disConnect)

	fmt.Println(http.ListenAndServe(":8080", router))
}

func viewHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	data, err := ioutil.ReadFile("dist/index.html")
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(w, string(data))
}

func connect(s *melody.Session) {
	u1 := uuid.NewV4().String()
	s.Set("id", u1)
	ret := &connection{Connection: s, ID: u1}
	socketMap[u1] = ret
}

func disConnect(s *melody.Session) {
	id := getID(s)
	delete(socketMap, id)
}

func getID(s *melody.Session) string {
	var ret string
	id, exist := s.Get("id")
	if exist {
		if xi, ok := id.(string); ok {
			ret = xi
		}
	}
	return ret
}

func echo(s *melody.Session, msg []byte) {
	text := string(msg)
	split := strings.Index(text, " ")
	code := text[0:split]
	body := text[split:len(text)]
	jsonBytes := ([]byte)(body)

	switch code {
	case "join":
		data := new(joinRequest)
		if err := json.Unmarshal(jsonBytes, data); err != nil {
			fmt.Println("JSON Unmarshal error:", err)
			return
		}
		join(s, data)
		break
	case "line":
		data := new(oneLine)
		if err := json.Unmarshal(jsonBytes, data); err != nil {
			fmt.Println("JSON Unmarshal error:", err)
			return
		}
		updateLine(data)
		break
	case "updates":
		data := new(updates)
		if err := json.Unmarshal(jsonBytes, data); err != nil {
			fmt.Println("JSON Unmarshal error:", err)
			return
		}
		updateLines(data)
		break
	case "follow":
		data := new(followLatestRequest)
		if err := json.Unmarshal(jsonBytes, data); err != nil {
			fmt.Println("JSON Unmarshal error:", err)
			return
		}
		followLatest(data)
		break
	default:
		break
	}
}

// トップページハンドラ
func index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}

func getBody(r *http.Request) []byte {
	length, err := strconv.Atoi(r.Header.Get("Content-Length"))
	if err != nil {
		return nil
	}

	body := make([]byte, length)
	length, err = r.Body.Read(body)
	if err != nil && err != io.EOF {
		return nil
	}

	return body
}

type connection struct {
	Connection *melody.Session
	ID         string
}

func join(s *melody.Session, data *joinRequest) {
	jsonBytes, err := json.Marshal(emitJoin{
		Type: "join",
		ID:   getID(s),
	})
	if err != nil {
		fmt.Println("JSON Marshal error:", err)
		return
	}
	m.BroadcastOthers(jsonBytes, s)
}

// ファイル全情報取得
func followLatest(data *followLatestRequest) {
	jsonBytes, err := json.Marshal(allLine{
		Type:     "text",
		FileName: data.FileName,
		Text:     data.Text,
	})
	if err != nil {
		fmt.Println("JSON Marshal error:", err)
		return
	}

	// 個人宛か全体かを分岐
	if data.ID == "" {
		m.Broadcast(jsonBytes)
	} else {
		socket, ok := socketMap[data.ID]
		if ok {
			socket.Connection.Write(jsonBytes)
		}
	}
}

// 1行更新
func updateLine(data *oneLine) {
	jsonBytes, err := json.Marshal(oneLine{
		Type: "line",
		R:    data.R,
		C:    data.C,
		L:    data.L,
	})
	if err != nil {
		fmt.Println("JSON Marshal error:", err)
		return
	}
	m.Broadcast(jsonBytes)
}

// 複数行更新
func updateLines(data *updates) {
	jsonBytes, err := json.Marshal(updates{
		Type:    "updates",
		Updates: data.Updates,
	})
	if err != nil {
		fmt.Println("JSON Marshal error:", err)
		return
	}
	m.Broadcast(jsonBytes)
}

type joinRequest struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

type emitJoin struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

type emitLeave struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

type followLatestRequest struct {
	Type     string `json:"type"`
	ID       string `json:"id"`
	FileName string `json:"fileName"`
	Text     string `json:"text"`
}

type allLine struct {
	Type     string `json:"type"`
	FileName string `json:"fileName"`
	Text     string `json:"text"`
}

type oneLine struct {
	Type string `json:"type"`
	R    int    `json:"r"`
	C    int    `json:"c"`
	L    string `json:"l"`
}

type unitUpdate struct {
	Sr int    `json:"sr"`
	Sc int    `json:"sc"`
	Er int    `json:"er"`
	Ec int    `json:"ec"`
	T  string `json:"t"`
}

type updates struct {
	Type    string       `json:"type"`
	Updates []unitUpdate `json:"updates"`
}
