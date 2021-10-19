# Contribute Code

You are welcome to contribute to project BFE. To contribute to BFE, you have to agree with the 
[Contributor License Agreement](https://cla-assistant.io/bfenetworks/api-server).

We sincerely appreciate your contribution. This document explains our workflow and work style.

## Workflow

BFE uses this [Git branching model](http://nvie.com/posts/a-successful-git-branching-model/). The following steps guide usual contributions.

1. Fork

   Our development community has been growing fastly; it doesn't make sense for everyone to write into the official repo.  So, please file Pull Requests from your fork.  To make a fork, just head over to the GitHub page and click the ["Fork" button](https://help.github.com/articles/fork-a-repo/).

1. Clone

   To make a copy of your fork to your local computers, please run

   ```bash
   git clone https://github.com/your-github-account/api-server
   cd bfe
   ```

1. Create the local feature branch

   For daily works like adding a new feature or fixing a bug, please open your feature branch before coding:

   ```bash
   git checkout -b my-cool-stuff
   ```


1. Build and test

   Users can build Conf Agent natively on Linux. 

   ```bash
   make
   ```

1. Keep pulling

   An experienced Git user pulls from the official repo often -- daily or even hourly, so they notice conflicts with others work early, and it's easier to resolve smaller conflicts.

   ```bash
   git remote add upstream https://github.com/bfenetworks/api-server
   git pull upstream develop
   ```

1. Push and file a pull request

   You can "push" your local work into your forked repo:

   ```bash
   git push origin my-cool-stuff
   ```

   The push allows you to create a pull request, requesting owners of this [official repo](https://github.com/bfenetworks/api-server) to pull your change into the official one.

   To create a pull request, please follow [these steps](https://help.github.com/articles/creating-a-pull-request/).

   If your change is for fixing an issue, please write ["Fixes <issue-URL>"](https://help.github.com/articles/closing-issues-using-keywords/) in the description section of your pull request.  Github would close the issue when the owners merge your pull request.

   Please remember to specify some reviewers for your pull request. If you don't know who are the right ones, please follow Github's recommendation.


1. Delete local and remote branches

   To keep your local workspace and your fork clean, you might want to remove merged branches:

   ```bash
   git push origin :my-cool-stuff
   git checkout develop
   git pull upstream develop
   git branch -d my-cool-stuff
   ```

### Code Review

-  Please feel free to ping your reviewers by sending them the URL of your pull request via IM or email. Please do this after your pull request passes the CI.

- Please answer reviewers' every comment. If you are to follow the comment, please write "Done"; please give a reason otherwise.

- If you don't want your reviewers to get overwhelmed by email notifications, you might reply their comments by [in a batch](https://help.github.com/articles/reviewing-proposed-changes-in-a-pull-request/).

- Reduce the unnecessary commits.  Some developers commit often.  It is recommended to append a sequence of small changes into one commit by running `git commit --amend` instead of `git commit`.


## Coding Standard

### Code Style

Our Golang code follows the [Golang style guide](https://github.com/golang/go/wiki/Style).

Our build process helps to check the code style. 

Please install pre-commit, which automatically reformat the changes to Golang code whenever we run `git commit`.  

### Unit Tests

Please remember to add related unit tests.

- For Golang code, please use [Golang's standard `testing` package](https://golang.org/pkg/testing/).

