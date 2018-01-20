# proxywall

A simple proxy that firewalls the requests that are passed through. Currently supports filtering by query parameters. 

Example config: 
```
{
  "url": "http:\/\/localhost:32769",
  "queryRestrictions": [
    {
      "key": "url",
      "valueRegex": "^https:\/\/www\\.google\\.com"
    }
  ]
}
```
Will allow only URLs that have query params where the value of `url` beings with `https://www.google.com`. 

Required ENV vars:

`PORT` and `CONFIG`, where `CONFIG` is the base64 encoding of the JSON above, like `ew0KICAidXJsIjogImh0dHA6XC9cL2xvY2FsaG9zdDozMjc2OSIsDQogICJxdWVyeVJlc3RyaWN0aW9ucyI6IFsNCiAgICB7DQogICAgICAia2V5IjogInVybCIsDQogICAgICAidmFsdWVSZWdleCI6ICJeaHR0cHM6XC9cL3d3d1xcLmdvb2dsZVxcLmNvbSINCiAgICB9DQogIF0NCn0=`
