GaC - Gmail as Code
===
![Pushed](https://github.com/nasa9084/gac/workflows/Pushed/badge.svg)

## HOW TO USE

You can install `gac` via `go get`:

``` shell
$ go get github.com/nasa9084/gac
```

At first, you need to do authenticate/authorize via Google OAuth. Steps are:

1. Enable Gmail API: Create a new project (or use existing project) on the [Google API Console](https://console.developers.google.com/) then enable Gmail API from `ENABLE APIS AND SERVICES`.
2. Create Credential: Create a new credential (OAuth client ID) from `CREATE CREDENTIALS` on `Credentials` page, with `Other` Application type (you can choose any name which is easy to understand for you).
3. Download Credential: Click the OAuth 2.0 Client ID name you created, then download credential file from `DOWNLOAD JSON` on the top of the page (assuming the filename is `credentials.json`).
4. Authenticate/Authorize: Run `gac auth`, then open oauth page on your browser. You can change the credentials file name to be loaded via `-c`/`--credentials-file` option. You may see "This app isn't verified" screen but you can go through via "advanced" button. After Authenticate and Authorize, successful screen will be shown. Close the window/tab of your browser and go back to your terminal. You can see a new file `token.json`.

Now you can use gac, with `credentials.json` and `token.json` in current dir.

### Filters

#### LIST Filters

with `credentials.json` and `token.json` in same dir:

``` shell
$ gac filter list
```

This command prints gmail filter list as YAML format.

You can change the file name of `credentials.json` to be loaded by `-c`/`--credentials-file` option. Can also be used `-` as the argument of the this option, meaning "read from stdin.
Although you cannot change the file name of `token.json` to be loaded, if you want to pass only refresh token not all content of `token.json`, in environment like CI, you can pass refresh token via `-t`/`--refresh-token` option.
