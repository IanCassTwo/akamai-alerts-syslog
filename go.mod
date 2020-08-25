module github.com/IanCassTwo/akamai-alerts-syslog

go 1.12

replace github.com/h2non/gock => gopkg.in/h2non/gock.v1 v1.0.14

replace github.com/akamai/AkamaiOPEN-edgegrid-golang => github.com/IanCassTwo/AkamaiOPEN-edgegrid-golang v0.9.6-0.20200825114129-a621aacf77d8

require (
	github.com/akamai/AkamaiOPEN-edgegrid-golang v0.9.5
	github.com/pcktdmp/cef v0.1.1
)
