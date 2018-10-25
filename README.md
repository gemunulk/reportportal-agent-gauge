## Build from Source

### Requirements
* [Golang](http://golang.org/)

### Compiling

Dependencies
```
go get ./…
```
Compilation
```
go run build/make.go
```
```
go run build/make.go --distro
```
### Creating distributable


For distributable across platforms: Windows and Linux for both x86 and x86_64

```
go run build/make.go --all-platforms
```
```
go run build/make.go --distro --all-platforms
```
#### Offline installation
* Download the plugin from [Releases](https://github.com/gemunulk)
```
gauge uninstall reportportal
```
```
cd /go/src/reportportal-agent-gauge/deploy
```
```
gauge install reportportal --file reportportal-0.0.1-darwin.x86_64.zip
```

### Configuration

* To set Report Portal properties add `env/default/reportportal.properties` file in to the gauge project.
```
REPORTPORTAL_SERVER  = http://35.189.175.106:8080
```
```
REPORTPORTAL_UUID    = UUID_FROM_REPORT_PORTAL
```
```
REPORTPORTAL_PROJECT_NAME = demo
```
```
REPORTPORTAL_LAUNCH_NAME = regression
```
```
REPORTPORTAL_TAGS = gui,api,v5.4.0
```

* Update the `plugin.json` file.
```
E.g.
{
  "Language": "java",
  "Plugins": [
    "html-report",
    "reportportal"
  ]
}

```
