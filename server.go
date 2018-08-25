package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"runtime/pprof"
	"strconv"

	"strings"

	"github.com/go-martini/martini"
)

var user *User = &User{name: "ariel", tenant: "ariel"}

var DebugVerbose bool

func stop() {
	//for breakpoints
}
func stop2(a []int) {
	//for breakpoints

}

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

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {

	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	log.SetOutput(os.Stdout)

	m := martini.Classic()

	m.Get("/prof/stop", func(w http.ResponseWriter, r *http.Request) {
		pprof.StopCPUProfile()

	})

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
		isGenderSensitive := getParamBool(r, "isGenderSensitive")
		isSpreadEvenly := getParamBool(r, "isSpreadEvenly")
		updateSubgroup(user, taskId, groupId, isUnite, isInactive, isGenderSensitive, isSpreadEvenly)
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
			//file := r.URL.Query().Get("file")
			//ec, err = Initialize(file)
			w.Write([]byte("Not implemented"))
			return
		} else {
			ec, err = Initialize2(user, taskId)
			ec.InitialErr = err
		}
		if ec == nil {
			w.Write([]byte("<html><body>Execution aborted. </br>\n" + err + "</br></body></html>"))
			return
		}

		limit := r.URL.Query().Get("limit")
		if len(limit) > 0 {

			ec.timeLimit, _ = strconv.Atoi(limit)
		} else {
			ec.timeLimit = 20
		}
		go RunBackTrack(ec)
		w.Write([]byte(`<html><body>Execution started.to check status, click <a href="/status?e=` + ec.ID() + `" >here</a></br>`))
		if len(err) > 0 {
			w.Write([]byte("Warnings: </br>" + strings.Replace(err, "\n", "<br/>\n", -1)))
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
			return
		}
		w.Write([]byte(`<html dir="rtl"><head>
						<script src="util.js"></script>
						<script src="forge.min.js"></script>
						<script>
							function decryptNames() {
								var cells = document.getElementsByName("encryptedCell");
								var pwd = document.getElementById("pwd").value;
								if (pwd.length > 0) {
									for (var i=0;i<cells.length;i++) {
										var cell = cells[i];
										cell.innerText = decrypt(pwd, cell.innerText)
									}
								}							
					        }
						</script>
						</head>
						<body>
						
						<h1>Execution status</h1>  <a href="/status?e=` + ec.id + `" >Refresh</a><br>
						סיסמת הצפנה<input type="text" id="pwd" value=""><input type="button" onClick= "decryptNames()" id="decrypt" value="פענח" /><br/>`))
		res, t := ec.GetStatusHtml()
		w.Write([]byte("time(sec):" + t))
		w.Write([]byte(res))
		w.Write([]byte("</body></html>"))
	})

	m.Get("/up", func(w http.ResponseWriter, r *http.Request) {

		w.Write([]byte(`<html>
					<script src="util.js"></script>
					<script src="forge.min.js"></script>
					<script>
						function makeKey(){
							var passcodeInput = document.getElementsByName("passcode")[0];
							if (passcodeInput.value != "") {
								passcodeInput.value = btoa(getKey(passcodeInput.value))
							}
						}
					</script>
                    <form action="" method="post" onSubmit="makeKey()" enctype="multipart/form-data">
						<p><input type="file" name="file" value="upload excel file">
						<p><input type="text" name="taskName" value="" placeholder="שם השיבוץ">
						<p><input type="text" name="passcode" value="" placeholder="סיסמת הצפנה - אופציונלי">
                        <p><button type="submit" >שלח</button>
                    </form>
                </html>`))
	})

	m.Get("/pref-graph", func(w http.ResponseWriter, r *http.Request) {
		taskId := getParamInt(r, "task")
		var ec *ExecutionContext
		var err string
		ec, err = Initialize2(user, taskId)
		if err != "" {
			w.Write([]byte("Error: </br>" + err))
			return
		}

		w.Write([]byte(`<html>
  	<head>
    <link href="graph.css" rel="stylesheet" />
    <meta charset=utf-8 />
    <meta name="viewport" content="user-scalable=no, initial-scale=1.0, minimum-scale=1.0, maximum-scale=1.0, minimal-ui">
    <title>גרף העדפות</title>
	<script src="cytoscape.min.js"></script>
	<script>
		var nodes = [`))
		comma := ""
		for i, p := range ec.pupils {
			c := "white"
			if p.IsMale() {
				c = "blue"
			}
			w.Write([]byte(fmt.Sprintf("%s { data: { id: '%d', name: '%s', color: '%s'}}", comma, i, p.name, c)))

			if comma == "" {
				comma = ","
			}
		}

		w.Write([]byte(`];
		var edges = [`))
		comma = ""
		for i, p := range ec.pupils {
			for j := 0; j < len(p.prefs); j++ {
				w.Write([]byte(fmt.Sprintf("%s { data: { source: '%d', target: '%d'}}", comma, i, p.prefs[j])))
				if comma == "" {
					comma = ","
				}
			}

		}
		//{ data: { source: 'j', target: 'g' } },
		w.Write([]byte(`];
	</script>
  </head>
  <body>
    <div id="cy"></div>
    <!-- Load appplication code at the end to ensure DOM is loaded -->
    <script src="graph.js"></script>
  </body>
</html>
`))
	})

	m.Get("/test/encrypt", func(w http.ResponseWriter, r *http.Request) {

		passcode := r.URL.Query().Get("passcode")

		content := "בדיקה"
		key, iv := AESencryptInit(passcode)
		enc := AESencrypt(key, iv, content)

		w.Write([]byte(enc))
	})

	m.Post("/up", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%v\n", "p./up")

		file, header, err := r.FormFile("file")
		taskName := r.FormValue("taskName")
		passcode := r.FormValue("passcode")
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

		uploadExcel(&User{tenant: "ariel"}, path, taskName, passcode)

		os.Remove(path)

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
