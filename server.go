package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"strings"

	"github.com/go-martini/martini"
)

var user *User = &User{name: "ariel", tenant: "ariel"}

func getParamInt(r *http.Request, param string) int {
	val := r.URL.Query().Get(param)
	res, err := strconv.Atoi(val)
	if err == nil {
		return res
	}
	return -1
}
func getParamBool(r *http.Request, param string) bool {
	val := r.URL.Query().Get(param)
	if strings.ToLower(val) == "true" {
		return true
	}
	return false
}

func main() {
	log.SetOutput(os.Stdout)

	m := martini.Classic()

	m.Get("/api/pupils", func(w http.ResponseWriter, r *http.Request) {
		task := getParamInt(r, "task")
		json, err := getPupilList2(&User{tenant: "ariel"}, task)
		if err == nil {
			w.Write(json)
		} else {
			//todo
		}

		//w.Write([]byte(`[{"id":"1", "name":"אריאל"}, {"id":"2", "name":"מעין"}, {"id":"2", "name":"יובל"}]`))
	})

	m.Post("/api/pupils", func(w http.ResponseWriter, r *http.Request) (int, string) {
		return 201, "Successufuly added"
	})

	m.Get("/api/tasks", func(w http.ResponseWriter, r *http.Request) {
		json, err := getTaskList(user)
		if err == nil {
			w.Write(json)
		}
	})

	m.Delete("/api/tasks", func(w http.ResponseWriter, r *http.Request) (int, string) {
		taskId := getParamInt(r, "task")
		deleteTask(user, taskId)
		return 201, "Successufuly deleted task"
	})

	m.Get("/api/subgroups", func(w http.ResponseWriter, r *http.Request) {
		taskId := getParamInt(r, "task")
		json, err := getSubgroupList2(user, taskId)
		if err == nil {
			w.Write(json)
		}
	})
	m.Post("/api/subgroup", func(w http.ResponseWriter, r *http.Request) (int, string) {
		taskId := getParamInt(r, "task")
		name := r.URL.Query().Get("name")
		createNewSubgroup(user, taskId, name)
		return 201, "Successufuly added"
	})

	m.Put("/api/subgroup", func(w http.ResponseWriter, r *http.Request) (int, string) {
		taskId := getParamInt(r, "task")
		groupId := getParamInt(r, "groupId")
		isUnite := getParamBool(r, "isUnite")
		isInactive := getParamBool(r, "isInactive")
		updateSubgroup(user, taskId, groupId, isUnite, isInactive)
		return 201, "Successfully updated"
	})

	m.Get("/api/pupil/prefs", func(w http.ResponseWriter, r *http.Request) {
		taskId := getParamInt(r, "task")
		pupilId := getParamInt(r, "pupilId")
		json, err := getPupilPrefs(user, taskId, pupilId)
		if err == nil {
			w.Write(json)
		}
	})
	m.Post("/api/pupil/prefs", func(w http.ResponseWriter, r *http.Request) (int, string) {
		taskId := getParamInt(r, "task")
		pupilId := getParamInt(r, "pupilId")
		decoder := json.NewDecoder(r.Body)

		defer r.Body.Close()
		err := setPupilPrefs(user, taskId, pupilId, decoder)
		if err == nil {
			return 201, "Updated Successfully"
		}
		return 500, err.Error()
	})

	m.Get("/api/subgroup/pupils", func(w http.ResponseWriter, r *http.Request) {
		taskId := getParamInt(r, "task")
		groupId := getParamInt(r, "groupId")
		json, err := getSubGroupPupils2(user, taskId, groupId)
		if err == nil {
			w.Write(json)
		}
	})
	m.Post("/api/subgroup/pupils", func(w http.ResponseWriter, r *http.Request) (int, string) {
		taskId := getParamInt(r, "task")
		groupId := getParamInt(r, "groupId")
		decoder := json.NewDecoder(r.Body)

		defer r.Body.Close()
		err := setSubGroupPupils2(user, taskId, groupId, decoder)
		if err == nil {
			return 201, "Updated Successfully"
		}
		return 500, err.Error()
	})

	m.Get("/run", func(w http.ResponseWriter, r *http.Request) {
		taskId := getParamInt(r, "task")
		var ec *ExecutionContext
		var err string
		if taskId == -1 {
			file := r.URL.Query().Get("file")
			ec, err = Initialize(file)
		} else {
			ec, err = Initialize2(user, taskId)
		}
		if ec == nil {
			w.Write([]byte("<html><body>Execution aborted. </br>\n" + err + "</br></body></html>"))
			return
		}

		limit := r.URL.Query().Get("limit")
		if len(limit) > 0 {

			ec.timeLimit, _ = strconv.Atoi(limit)
		} else {
			ec.timeLimit = 5
		}
		go RunBackTrack(ec)
		w.Write([]byte(`<html><body>Execution started.to check status, click <a href="/status?e=` + ec.ID() + `" >here</a></br>`))
		if len(err) > 0 {
			w.Write([]byte("Warnings: </br>" + err))
		}
		w.Write([]byte("</body></html>"))
	})

	m.Get("/cancel", func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("e")

		ec := FindExecutionContext(id)
		ec.Cancel = true
	})

	m.Get("/status", func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("e")

		ec := FindExecutionContext(id)
		if ec == nil {
			w.Write([]byte("context not found"))
		}
		w.Write([]byte(`<html dir="rtl"><body><h1>Execution status</h1>  <a href="/status?e=` + ec.id + `" >Refresh</a><br>`))
		res, t := ec.GetStatusHtml()
		w.Write([]byte("time(sec):" + t))
		w.Write([]byte(res))
		w.Write([]byte("</body></html>"))
	})

	m.Get("/up", func(w http.ResponseWriter, r *http.Request) {

		w.Write([]byte(`<html>
                    <form action="" method="post" enctype="multipart/form-data">
						<p><input type="file" name="file" value="upload excel file">
						<p><input type="text" name="taskName" value="שם עבודה">
                        <p><button type="submit">Submit</button>
                    </form>
                </html>`))
	})
	m.Post("/up", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%v\n", "p./up")

		file, header, err := r.FormFile("file")
		taskName := r.FormValue("taskName")
		defer file.Close()

		if err != nil {
			fmt.Fprintln(w, err)
			return
		}
		path := "/tmp/" + header.Filename
		out, err := os.Create(path)
		if err != nil {
			fmt.Fprintf(w, "Failed to open the file for writing")
			return
		}
		defer out.Close()
		_, err = io.Copy(out, file)
		if err != nil {
			fmt.Fprintln(w, err)
		}

		uploadExcel(&User{tenant: "ariel"}, path, taskName)

		// the header contains useful info, like the original file name
		fmt.Fprintf(w, "File %s uploaded successfully.", header.Filename)
	})
	/*
		m.Use(cors.Allow(&cors.Options{
			AllowOrigins:     []string{"http://localhost:3001"},
			AllowMethods:     []string{"PUT", "PATCH", "GET", "POST"},
			AllowHeaders:     []string{"Origin"},
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: true,
		}))
	*/
	m.Use(martini.Static("ui/"))

	m.Run()

}
