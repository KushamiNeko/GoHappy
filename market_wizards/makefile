css:
	~/programming_projects/tools/go/web_templates/bin/web_templates -input ./templates -operation scss
	~/programming_projects/tools/go/web_templates/bin/web_templates -input ./templates -operation templates
	rm templates/**/**/*.css

js:
	~/programming_projects/tools/go/web_templates/bin/web_templates -input ./templates -operation dart 
	~/programming_projects/tools/go/web_templates/bin/web_templates -input ./templates -operation templates
	rm templates/**/**/*.js*

web:
	make css
	make js
