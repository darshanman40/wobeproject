# Wobe Project

## Project Requirement:
The server should

1. Accept a POST request of the below form,
<br>{
"input": "blablabla"
}

2. Response should be the reverse of the "input" string

"albalbalb"

The server should run on a docker container and you should be able to scale the number of containers horizontally.

## How to run

It requires docker, go library, git to run

### Steps to run without docker:
1. Clone the repository under `$GOPATH/src/github.com` directory
2.  Run `go get ./...` to install dependencies
3. Run `go run main.go -env prod`
NOTE: flags for -env are "prod", "dev" and "debug"

This will keep the server on `http://localhost:8081`

### Steps to run with docker:
1. Clone the repository
2. Run `docker build -t wobeproject .`
3. Run `docker run -p 6060:8081 wobeproject -env prod` <br>
OR<br>
Run `docker run -d -p 6060:8081 wobeproject -env prod` to run in background.

This will keep the server on `http://localhost:6060`

### Test

1. Using postman application, make "POST" request with one of above urls.
2. Select "Body" tab then "x-www-form-urlencoded" option.
3. Put key as "input" and value as any non-palindrome text (to get best results ;-) ) hit "Send"  


## About project
It has two flags, <br>

1) -env: environment of project. It pulls environment specific configuration. It supports "dev", "debug" and "prod" as values.



<table>
<tr>
<td></td>
<td>dev</td>
<td>debug</td>
<td>prod</td>
</tr>

<tr>
<td>tracelevel</td>
<td>infolevel, warnlevel,<br>errorlevel, paniclevel </td>
<td>infolevel, warnlevel,<br>errorlevel, paniclevel, debuglevel </td>
<td>infolevel, warnlevel,<br>errorlevel, paniclevel </td>
</tr>

<tr>
<td>stacktrace</td>
<td>false for all levels</td>
<td>false for all levels</td>
<td>true for only errors and panics</td>
</tr>

<tr>
<td>erroroutput (meant for writing logs in file ONLY if provided in code)</td>
<td>true for all levels</td>
<td>true for warn, panic,debug</td>
<td>false for all</td>
</tr>

<tr>
<td>caller (line number of caller and filename)</td>
<td>true for all</td>
<td>true for warn, error,panic,debug</td>
<td>false for all</td>
</tr>

<tr>
<td>caller_skip (skip line no. of caller)</td>
<td>Its not stable so please don't change</td>
<td>Its not stable so please don't change</td>
<td>Its not stable so please don't change</td>
</tr>

</table>

2) -config: To provide config file with path. Default value doesn't work in docker container, so I added working value in Dockerfile
