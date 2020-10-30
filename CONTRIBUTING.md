# Contributing Guide

Welcome to the TKEStack family. TKEStack is an open, comprehensive, strong open source project. 

Here, you can use the release of TKEStack to build an efficient and stable container cloud platform, Or contribute to the TKEStack project, submit your ideas, questions, and code, share your work with the community, help more people to use TKEStack. Use technology to contribute to the world.

This article will help you understand how the TKEStack project is organized, guide you through how you can raise issues, write code to fix problems or implement new features, and review and merge your work.

Welcome to read this article and feedback your ideas to us. Looking forward to your first submission.

# How to Create Issue

TKEStack use [Issue](https://github.com/tkestack/tke/issues) report bugs and feature request, discuss and management, detailed issue guidelines please check: [Managing your work with issues](https://docs.github.com/en/github/managing-your-work-on-github/about-issues).

TKEStack has prepared the Issue template. Please write Issue in the format of the template. In general, follow the template to clearly describe the following situations:

- What **environment/version/situation** is the problem?

- How to **reproduce** the problem.

- What do you **expect** for?

TKEStack members or administrators will participate in discussions, respond to messages, and arrange developers to follow up the issue.

# How to Pull Request (PR)

## Create PR

Depending on the type of PR, you should select the appropriate branch. Refer to [Branch and Version Policy](#Branch and Version Policy) for the branch instructions.

- For the new feature, use the MASTER branch as the base branch, and after the development, create PR and merge it into the Master branch.

- For bug fixes, use the latest Master or Release Branch as the base branch, and after development, create the PR merge into the corresponding branch.

  Meanwhile, if other branches have the same problem, confirm with the administrator and cherry-pick to the corresponding branch.

- When creating PR, names are in line with Angular specifications, name contents should be meaningful, and organize commits before submitting the PR:

   ```<type>(<scope>): <subject>```

  - **type**: is used to indicate the category of submission and can only be identified as:
    - **feat**: New feature
    - **fix**: Fix a bug
    - **docs**: Documents
    - **test**: Add tests
    - **style**: Format (changes that do not affect code running)
    - **refactor**: Refactoring (i.e., code changes that do not add new features or fix bugs)
    - **ci**: Automated processes
    - **build**: Changes the build process or tooling
    - **perf**: Performance optimization related
    - **revert**: revert
    
  - **scope**: is used to describe the scope of influence of the submission. The recommended value is the module or function name in TKEStack, for example:
    - **installer**
    - **platform**
    - **cluster**
    - **gateway**
    - **console**
    - **auth**
    - **addon**
    - **business**
    - **registry**
    - **monitor**
    - **notify**
    - **audit**
    
  - **subject**: is a short description of the purpose of the submission, no more than 50 characters long.
    - Start with a verb and use the present tense in the first person, such as' change ', not 'changed' or 'changes'
    - Lowercase first letter
    - No period at the end
  
  For more detailed specification, see the link below. Follow the specification, PR with feat and fix type will appear in the Changelog.
  
  - [@commitlint/config-angular](https://github.com/conventional-changelog/commitlint/tree/master/%40commitlint/config-angular)
  
- After creating PR, Github will automatically start the CI process, verify the PR information and check the compiling and testing, while Github will automatically notify Reviewer of a code review.

- Adjust PR according to reviewer's comments, refer to [Reviewing changes in pull requests](https://docs.github.com/en/github/collaborating-with-issues-and-pull-requests/reviewing-changes-in-pull-requests).

- Wait for the administrator's approval, if approved, the contribution is completed.

  Note: Some abbreviations used by Reviewer

  - **PR**: Pull Request. 
  - **LGTM**: Looks Good To Me.
  - **SGTM**: Sounds Good To Me. 
  - **WIP**: Work In Progress.
  - **PTAL**: Please Take A Look. 
  - **TBR**: To Be Reviewed. 
  - **TL;DR**: Too Long; Didn't Read. 
  - **TBD**: To Be Done(or Defined/Discussed/Decided/Determined). 
  - **DNM**: Do not merge.


# Project Management

This section describes the TKEStack's version policy, explains the TKEStack branch and tag rules, 
helps users quickly select the appropriate branch and version to start their work.

## Branch and Version Policy

- The project uses the Master branch, the Release-1.x release branch, and other temporary branches

- The project uses git tag as the version number, in v1.x.y format, and supports alpha and beta tags

- The Master branch accepts developers' PR merge, and the feature and fix will be released in the next version v1.x.0

- When a new v1.x.0 version is tagged, create a new release-1.x branch from the tag 

- The release-1.x branch should provide a stable version with no obvious bugs, only the critical fix will be merged and v1.x.y will be released timely

- The administrator maintains nearly two versions of the branch, and after the critical fix is submitted on the master, the administrator cherry-pick to the maintenance branch; If there is any other fix or feat that requires cherry-pick to release branch, submit an issue or PR to notify the administrator to handle.

- The administrator is responsible for maintaining the branches and versions to ensure the long-term health, robustness, and traceability


## Release notes

TKEStack release versions as planned, for important features, the user can focus on [Roadmap](https://github.com/tkestack/tke/wiki/TKEStack-Roadmap).

To get TKEStack version scheduling, through [milestones](https://github.com/tkestack/tke/milestones) to view more detailed project development state.

At the same time, the Release TKEStack contains Release Notes and Changelog, which are synchronized on the [TKEStack documentation site](https://tkestack.github.io/docs/).

## Proposals

For changes that have a greater scope of impact, such as underlying networks, installation processes, kay data structure changes, etc., design Proposals need to be planned in advance and submitted (or comment under the Issue).

Proposals submitted documents to the [design-proposals](https://github.com/tkestack/tke/tree/master/docs/design-proposals) to the public, the user feedback through Issue comments, the administrator and TOC finally decide the Proposals results.