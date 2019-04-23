package util

// version information populated at compile time with govvv

//Version contents of ./VERSION file, if exists, or the value passed via the -version option
var Version string

//BuildDate RFC3339 formatted UTC date
var BuildDate string

//GitCommit short commit hash of source tree
var GitCommit string

//GitBranch current branch name the code is built off
var GitBranch string

//GitState whether there are uncommitted changes (clean or dirty)
var GitState string

//GitSummary output of git describe --tags --dirty --always
var GitSummary string
