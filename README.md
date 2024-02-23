 # NomadInspect

## Description
NomadInspect is a project that allows users to query their nomad cluster searching the deployed jobs to match them against certain strings. With this functionality you can quickly search through all the namespaces which jobs contain info you are looking for.

## Installation
You only have to compile the go binary with go build.

## Usage
To launch nomadinspect you only need to invoke the binary with the Nomad environment variables loaded
``` 
nomadinspect [-namespace <namespace>] -match <match> [-match <match>] [-and]
```
Nomadinspect will query all the jobs deployed in the nomad namespaces included in the command and check if said jobs match against certain parameters. The options of **nomadinspect** are:

* match - to select the string to look for in the jobs definition. 
* namespace - to select which namespace we want to look into
* and - this option is just in case we want to look for jobs that match all the parameters configured when invoking the program.

```bash
# Example that would look through all the jobs in namespaces [admin, test, utilities] and return the jobs that contains the string admin and the string notallowed
nomadinspect -namespace admin -namespace test -namespace utilities -match admin -match notallowed -and
```

![Spinner Screen](https://github.com/smorenodp/nomadinspect/blob/master/images/spinner_screen.png)

The first screen is the loading one, in this screen it shows feedback about the namespaces it's looking and the jobs contained in them. After obtaining all the information, it transitions into the list screen.

![List Screen](https://github.com/smorenodp/nomadinspect/blob/master/images/list_screen.png)

In the list screen you can review all the jobs that matched with the query sent, with the job name as the title and the namespace as the subtitle. You can press `enter` to review the job definition with the matches highlighted.

Reviewing the job definition you can either move with the arrow keys, mouse wheel or press the key `m` to move between all the matches in the JSON.

To return to the list screen you just press `enter` again.