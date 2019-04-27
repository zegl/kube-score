# Contributing to kube-score

👋 Hey, thanks for taking the time to contribute! Your help is appreciated.

## How can I contribute?

### Reporting bugs

Bug reports are always welcome, and should be reported as a [GitHub issue](https://github.com/zegl/kube-score/issues/new/choose).

It's easy to open an issue, all you need to do is to answer the following questions:

1. Which version of kube-score are you using?
2. What did you do?
3. What did you expect to see?
4. What did you see instead?

### Feature requests

Feature requests are always welcome, this should also be done as a [GitHub issue](https://github.com/zegl/kube-score/issues/new/choose).

There is currently no set template for feature requests, but largely the same template as the issues can be used.

Describe the feature that you would like to see as clearly as possible.

If the feature request is related to a new "check", include example objects (in YAML format) that are OK, and that should trigger a failure. 

### Contributing code

Code contributions are welcome as GitHub Pull Requests.

#### Good commit messages

kube-score tries to use the same commit message format as [the Go programming language](https://golang.org/doc/contribute.html#commit_messages).

Example of a good commit message:

```
score/container: always add a comment if Container Image Pull Policy fails

Fixes #79
```


The first line of the commit message should contain a short description of the change, prefixed by the primary affected package.

Additional lines can be used if a longer explanation of the change is needed.

Issues should be referenced with the syntax `Fixes #123` or `Updates #123` to track that this change is related to an issue.

#### After submit

After a PR has been opened, the complete set of tests will be run automatically. All tests needs to pass before the PR can be merged.

kube-score is using [bors-ng](https://bors.tech/documentation/) for merging Pull Requests, after the PR has been approved by a maintainer,
the `bors r+` command will be run, this will rebase the PR on master, and run the tests again. If the tests run, the PR will be merged!
