# Cloudflare Bulk Redirects As Redirector

### Redirector
A redirector is a domain, that does a 302 from any url of  that domain to another domain, copying url path and query parameter. In other words, it just changes the host of the website.
> Source domains are required to be proxied by cloudflare for redirectors to work.

### Usage
This script creates new redirects, which forward the requests from `given_host` to the target `host`. It replaces all the redirects in the `target_list` with newly created redirects through the cloudflare api(with token).

#### Command Line Usage
`host` - Target Host

`dfile` - data file
```
./cbrar --host "https://example.com" --dfile data.csv
```

#### Data File Syntax
```
<given_host>, <target_list> , <api_token> , <account_id>
...
```

### Required API Permissions, 
* Account: Account Rulesets : Edit
* Account: Account Filter Lists : Edit

### Setup
1. First of all, create a redirect list with inital redirects, the `target_list` will be in the url.

2. Next create crate a Redirect Rule, just a basic redirect rule attaching to the above created list,
```
http.request.full_uri in $created_list
```

3. That's it, now create a api token, get account id and you got all the things required.

### Limits
Free Cloudflare accounts allow 15 Url Redirect Rules, 5 Bulk Redirect List, and 20 redirects per list. So, One free account is capable of creating 100 Redirectors.
> Well, More Redirectors can be created, by using other kind of rule systems supported by cloudflare :)