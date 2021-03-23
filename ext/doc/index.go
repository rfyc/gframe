package doc

var html_index = `
	{{range $_, $v := .keys}}
		<div>
			{{$v}}  
 		</div> 
  	{{end}}
`
