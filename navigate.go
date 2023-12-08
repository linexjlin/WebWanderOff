package main

import (
	"log"
	"net/http"
	"text/template"
)

var htmlTemplate = `
<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<title>Site List</title>
	<style>
		body {
			font-family: Arial, sans-serif;
			margin: 0;
			padding: 0;
			background-color: #f4f4f4;
		}
		.container {
			max-width: 800px;
			margin: 20px auto;
			background-color: #fff;
			padding: 20px;
			border-radius: 8px;
			box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
		}
		.site-list {
			list-style: none;
			padding: 0;
		}
		.site-item {
			margin-bottom: 15px;
			padding-bottom: 15px;
			border-bottom: 1px solid #eeeeee;
		}
		.site-item:last-child {
			border-bottom: none;
		}
		.site-icon {
			max-width: 50px;
			max-height: 50px;
			vertical-align: middle;
		}
		.site-name a {
			font-size: 20px;
			text-decoration: none;
			color: #333;
		}
		.site-description {
			margin-top: 5px;
			font-size: 14px;
		}
	</style>
</head>
<body>
	<div class="container">
		<h1>Available Sites</h1>
		<ul class="site-list">
		{{ range . }}
			<li class="site-item">
				
				<div class="site-name">
					<img class="site-icon" src="data:image/*;base64,{{ .Icon }}">
					<a href="http://{{.ListenAddr}}">{{.Name}}</a> - {{.Description}}
				</div>
			</li>
		{{ end }}
		</ul>
	</div>
</body>
</html>
`

func renderSiteList(w http.ResponseWriter, r *http.Request, configs []Config) {
	tmpl, err := template.New("siteList").Parse(htmlTemplate)
	if err != nil {
		// Handle error
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, configs)
	if err != nil {
		// Handle error
		http.Error(w, "Error executing template", http.StatusInternalServerError)
		return
	}
}

func navigateServe(addr string, configs []Config) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		renderSiteList(w, r, configs)
	})
	log.Println("Navigate site listening on:", addr)
	log.Println(http.ListenAndServe(addr, mux))
}
