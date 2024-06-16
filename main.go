package main

import (
	"embed"
	"html/template"
	"math/rand"
	"net/http"
	"strconv"
	"strings"

	"github.com/syumai/workers"
)

//go:embed pages
var fs embed.FS

func main() {
	for k, v := range handlers() {
		http.HandleFunc(k, v)
	}
	workers.Serve(nil)
}

func handlers() map[string]http.HandlerFunc {
	return map[string]http.HandlerFunc{
		"/":     index(),
		"/spin": spin(),
	}
}

func index() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p, err := fs.ReadFile("pages/index.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl := template.Must(template.New("index").Funcs(template.FuncMap{"add": add}).Parse(string(p)))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		tmpl.Execute(w, nil)
	}
}

func spin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		member := r.FormValue("member")
		ms := strings.Split(member, "\n")

		members := make([]string, 0)
		for _, m := range ms {
			s := strings.TrimSpace(m)
			members = append(members, s)
		}

		group := r.FormValue("group")
		ig, err := strconv.Atoi(group)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		max := len(members) / ig

		rand.Shuffle(len(members), func(i, j int) {
			members[i], members[j] = members[j], members[i]
		})

		groups := make([][]string, ig)

		var i int
		for k, v := range members {
			if k != 0 && k%max == 0 {
				i++
			}
			if i >= ig {
				i = 0
			}
			groups[i] = append(groups[i], v)
		}

		p, err := fs.ReadFile("pages/index.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl := template.Must(template.New("index").Funcs(template.FuncMap{"add": add}).Parse(string(p)))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		tmpl.Execute(w, groups)
	}
}

func add(a, b int) int {
	return a + b
}
