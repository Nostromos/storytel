package story

// const tpl = `<!DOCTYPE html>
// <html>
// <head>
//   <meta charset="utf-8">
//   <title>{{.Title}}</title>
//   <style>
//     .option-button {
//       background-color: #4CAF50;
//       color: white;
//       border: none;
//       padding: 10px 20px;
//       font-size: 16px;
//       cursor: pointer;
//       margin: 5px;
//     }
//   </style>
// </head>
// <body>
//   <h1>{{.Title}}</h1>
//   <div>
//     {{range .Story}}
//       <p>{{.}}</p>
//     {{end}}
//   </div>
//   <div>
//     {{range .Options}}
//       <button class="option-button" onclick="window.location.href='{{printf "/%s" .Arc}}'">
//         {{.Text}}
//       </button>
//     {{end}}
//   </div>
// </body>
// </html>`

// const (
// 	defaultStoryPath = "./gopher.json"
// )

// func story () {
// 	// load our story
// 	story, err := parser.LoadStory(defaultStoryPath)
// 	if err != nil {
// 		panic(err)
// 	}

// 	fmt.Println(story)

// 	// parse our template
// 	templ, err := template.New("story").Parse(tpl)
// 	if err != nil {
// 		panic(err)
// 	}

// 	fmt.Println(templ)

// 	handler.HandlePathing()
// }
