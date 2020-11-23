# First run

* Go to the site with [documentation](https://developers.google.com/sheets/api/quickstart/go) and click 
on the "Enable the Google Sheets API" button.
Then download the file credentials.json and move it to the root folder of the project. 
* Run the code / binary. The program will ask you to log in using your Google account. 
Then insert the token into the console. Arguments at the first start of the program: `google update`

# Run
Run the program with the `run` argument

# Config file

## Env
| Name | Description | Is Required? | Default | 
| ----- | ----- | ----- | ----- |
| spreadsheet_id | ID of the Google Spreadsheet | + | - |
| starting_time | Unix-time from which time to collect issue | - | 0 (collects issues of any date) |
| tmpl_sheet_id | ID of the Google sheet where the table will be copied from. This is convenient if your table has a layout | - | If omitted, the table is created |
| services | List of the control version systems (cvs) when this script will collect data | + | - |
| services.* | Service name. This is a special key that cannot be changed after the first launch | + | - |
| services.type | Type of API the script will use. Currently only gitlab is available | + | - |
| services.url | API address | + | - |
| services.available | Availability level. `external` - the script will collect all issues, `internal` - the script will collect issues with only only created by you or assigned on you | - | external|
| services.token | Token of the user who will collect issues from cvs | + | - | 
| members | List of members that the script will navigate by if availability is internal | - | - |
| members.* | Members name. This is a special key that cannot be changed after the first launch | + | - |
| name | First/Last name | + | - |
| members.services | List of cvs that users are connected to by key-value type. The key is the name of the cvs that you described above, and the value is the user's username in this system | + | - |

Example:
```yaml
spreadsheet_id: 'code'

starting_time: 0
tmpl_sheet_id: 0

services:
  test-internal:
    type: 'gitlab'
    url: 'http://gitlab.com/api/v4'
    available: 'internal'
    token: '1234456789qwerty'

  test-external:
    type: 'gitlab'
    url: 'http://self-hosted-gitlab.com/api/v4'
    available: 'external'
    token: '1234456789qwerty'

members:
  discore:
    name: 'Nikita'
    services:
      - test-external: 'discore'
      - test-internal: 'discoreme'
  alexey:
    name: 'Alexey'
    services:
      - test-internal: 'alexey'


```

# Dev:
* `go get github.com/mattn/go-sqlite3` (which requires gcc)