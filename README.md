# An URL shortener written in Golang
Inspired by Mathias Bynens' [PHP URL Shortener](https://github.com/mathiasbynens/php-url-shortener), and triggered by a wish to learn Go, I wanted to try and see if I could build an URL shortener in Go.

## Features

* Redirect to your main website when no slug, or incorrect slug, is entered, e.g. `http://wiere.ma/` → `http://samwierema.nl/`.
* Generates short URLs using only `[a-z0-9]` characters.
* Doesn’t create multiple short URLs when you try to shorten the same URL. In this case, the script will simply return the existing short URL for that long URL.

## Installation
1. Download the source code and install it using the `go install` command.
2. Use `database.sql` to create the `redirect` table in a database of choice.
3. Create a config file in `/path/to/.go-url-shortener/` named `config.(json|yaml|toml)`. Use `config-example.json` as a example.
4. Run the program as a daemon using one of the many methods: write a script for [upstart](https://launchpad.net/upstart), init, use [daemonize](http://software.clapper.org/daemonize/), [Supervisord](http://supervisord.org/), [Circus](http://circus.readthedocs.org/) or just plain old `nohup`. You can even start (and manage) it in a `screen` session.
5. Adding the following configuration to Apache (make sure you've got [mod_proxy](http://httpd.apache.org/docs/2.2/mod/mod_proxy.html) enabled):
```
<VirtualHost *:80>
	ServerName your-short-domain.ext

	ProxyPreserveHost on
	ProxyPass / http://localhost:8080/
	ProxyPassReverse / http://localhost:8080/
</VirtualHost>
```

### Using the example init script
You will find an example init script in the `scripts` folder. To use, you **must** at least change the GOPATH line to point to your Go root path.

## To-do
* Add tests
* Add checks for duplicate slugs (i.e. make creation of slugs better)

## Author
* [Sam Wierema](http://wiere.ma)
