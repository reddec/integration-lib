update_godoc:
	find ./ -type d -not -path "./.*" | sed 's|./||' | xargs -n 1 -i curl "https://godoc.org/github.com/reddec/integration-lib/{}"