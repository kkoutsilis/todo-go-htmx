.PHONY: templ run dev

templ:
	go tool templ generate

run: templ
	go run .

dev: 
	go tool templ generate --watch &
	go run . 
