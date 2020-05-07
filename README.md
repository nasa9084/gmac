GMaC - Gmail as Code
===
![Test](https://github.com/nasa9084/gmac/workflows/Test/badge.svg)

## HOW TO USE

### Installation

You can install `gmac` via `go get`:

``` shell
$ go get github.com/nasa9084/gmac
```

### Authentication/Authorization

At first, you need to do authenticate/authorize via Google OAuth. Steps are:

1. Enable Gmail API: Create a new project (or use existing project) on the [Google API Console](https://console.developers.google.com/) then enable Gmail API from `ENABLE APIS AND SERVICES`.
2. Create Credential: Create a new credential (OAuth client ID) from `CREATE CREDENTIALS` on `Credentials` page, with `Other` Application type (you can choose any name which is easy to understand for you).
3. Download Credential: Click the OAuth 2.0 Client ID name you created, then download credential file from `DOWNLOAD JSON` on the top of the page (assuming the filename is `credentials.json`).
4. Authenticate/Authorize: Run `gmac auth`, then open oauth page on your browser. You can change the credentials file name to be loaded via `-c`/`--credentials-file` option. The file will be copied into `$HOME/.gmac/credentials.yml` and you do not need to specify credentials file path anymore. You may see "This app isn't verified" screen but you can go through via "advanced" button. After Authenticate and Authorize, successful screen will be shown. Close the window/tab of your browser and go back to your terminal. OAuth token, including refresh token, will be saved in `$HOME/.gmac/token.json`.

### Filters

#### LIST Filters

``` shell
$ gmac get filters
```

This command prints gmail filter list. If you want to backup current filter list in apply-able format, you can do with `-o` option:

``` shell
$ gmac get filters -o yaml > filters.yml
```

#### APPLY Filters

``` shell
$ gmac apply -f filters.yml
```

This command applies given filters.yml to your Gmail Filters. This command removes all existing filters at first, then create filters defined in given YAML file. So note that you add a new filter via Gmail UI and do not add it into YAML file, the filter you added only via UI will be removed when the next time you run this command with your outdated config.

##### Filter Configuration

The filters definition is written in YAML format, defined by the scheme described below.

Generic placeholders are defined as follows:

* `<bool>`: a boolean value that can take the values `true` or `false`
* `<string>`: a string value
* `<integer>`: a 64-bit integer value

The whole resource file format is:

``` yaml
kind: Filter

# List of Filter Objects
filters:
  - <Filter Object>
```

###### Filter Object

``` yaml
# Configure Filter Criteria
criteria:
  <FilterCriteria Object>

# Configure Filter Action
action:
  <FilterAction Object>
```

###### FilterCriteria Object

``` yaml
# Filter by email sender's display name or email address.
from: <string>

# Filter by email recipient's display name or email address
# includes recipients in th "to", "cc", and "bcc" header fields.
to: <string>

# Filter by the message's subject (case-insensitive).
subject: <string>

# Filter by query, only return messages matching the query.
query: <string>

# Filter by query, only return messages NOT matching the query.
negated_query: <string>

# Filter by the size of the entire RFC822 message in bytes,
# including all headers and attachments.
# larger_than and smaller_than are exclusive.
larger_than: <integer>
smaller_than: <integer>

# Filter by whether the message has any attachment or not.
# Default is false.
has_attachment: <bool>

# Filter by whether the response should exclude chats.
# Default is false.
exclude_chats: <bool>
```

###### FilterAction Object

``` yaml
# Archive the messages.
archive: <bool>

# Mark the messages as read.
mark_as_read: <bool>

# Make starred the messages.
star: <bool>

# Add given label to the messages.
add_label: <string>

# Forward the message to given email address.
forward_to: <string>

# Delete the messages.
delete: <bool>

# Never mark the messages as SPAM.
never_mark_as_apam: <bool>

# Mark the messages as important or never mark the messages as important.
# Valid values are "always" or "never".
important: <string>

# Set the message's category as given one.
# Valid values are:
# - "primary"
# - "social"
# - "updates"
# - "forums"
# - "promotions".
# and there're some aliases:
# - "personal" and "main" for primary
# - "update" and "new" for updates
# - "forum" for forums
# - "promotion" for promotions
category: <string>
```
