package main

import (
	"fmt"
	"github.com/codegangsta/martini"
	"github.com/codegangsta/martini-contrib/cors"
	"github.com/codegangsta/martini-contrib/render"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

func main() {
	log.SetOutput(os.Stdout)

	m := martini.Classic()

	m.Use(render.Renderer(render.Options{
		Directory: "templates", // Specify what path to load the templates from.
		Layout:    "layout",    // Specify a layout template. Layouts can call {{ yield }} to render the current template.
		Charset:   "UTF-8",     // Sets encoding for json and html content-types.
	}))

	m.Get("/api/files", func(w http.ResponseWriter, r *http.Request) {
		sb := NewStringBuffer()
		sb.Append(`{"files":[`)
		files, err := ioutil.ReadDir("c:/temp/groupallocation")
		if err != nil {
			log.Fatal(err)
		}

		for i, file := range files {
			if i > 0 {
				sb.Append(",")
			}
			sb.AppendFormat(`{"name":"%s"}`, file.Name())
		}

		sb.Append("]}")
		w.Write([]byte(sb.ToString()))
	})

	m.Get("/run", func(w http.ResponseWriter, r *http.Request) {
		file := r.URL.Query().Get("file")
		ec, err := Initialize(file)
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
                        <p><button type="submit">Submit</button>
                    </form>
                </html>`))
	})
	m.Post("/up", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%v\n", "p./up")

		file, header, err := r.FormFile("file")
		defer file.Close()

		if err != nil {
			fmt.Fprintln(w, err)
			return
		}

		out, err := os.Create("c:/temp/groupallocation/" + header.Filename)
		if err != nil {
			fmt.Fprintf(w, "Failed to open the file for writing")
			return
		}
		defer out.Close()
		_, err = io.Copy(out, file)
		if err != nil {
			fmt.Fprintln(w, err)
		}

		// the header contains useful info, like the original file name
		fmt.Fprintf(w, "File %s uploaded successfully.", header.Filename)
	})

	m.Use(cors.Allow(&cors.Options{
		AllowOrigins:     []string{"http://localhost:3001"},
		AllowMethods:     []string{"PUT", "PATCH", "GET", "POST"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	//m.Use(martini.Static("C:/Users/i022021/AppData/Roaming/npm/group-allocation/build"))

	m.Run()

}
