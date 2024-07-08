1.14.2
- 2024-07-08 Release v1.14.2
- 2024-07-08 fix: Adds check for tag format to release script, also makes release script error on improper tags
- 2024-07-08 Merge pull request #251 from norwoodj/dependabot/go_modules/helm.sh/helm/v3-3.15.2
- 2024-07-08 chore(deps): bump helm.sh/helm/v3 from 3.15.1 to 3.15.2

1.14.1
- 2024-07-08 Release 1.14.1
- 2024-07-08 Merge pull request #248 from WePlan-Software/master
- 2024-07-08 Update go.mod via 'go mod tidy'
- 2024-07-06 chore: generates documentation using new v1.14.0

1.14.0
- 2024-07-06 Release 1.14.0
- 2024-07-06 Merge pull request #237 from norwoodj/dependabot/github_actions/docker/setup-qemu-action-3
- 2024-07-06 chore(deps): bump docker/setup-qemu-action from 2 to 3
- 2024-07-06 Merge pull request #236 from norwoodj/dependabot/github_actions/crazy-max/ghaction-import-gpg-6
- 2024-07-06 chore(deps): bump crazy-max/ghaction-import-gpg from 5 to 6
- 2024-07-06 Merge pull request #238 from norwoodj/dependabot/docker/alpine-3.20
- 2024-07-06 chore(deps): bump alpine from 3.19 to 3.20
- 2024-07-06 Merge pull request #239 from norwoodj/dependabot/github_actions/docker/setup-buildx-action-3
- 2024-07-06 chore(deps): bump docker/setup-buildx-action from 2 to 3
- 2024-07-06 Merge pull request #246 from norwoodj/jnorwood/fix-goreleaser
- 2024-07-06 fix: Fixes goreleaser action by using new flag
- 2024-07-06 Merge pull request #240 from norwoodj/dependabot/github_actions/goreleaser/goreleaser-action-6
- 2024-07-06 chore(deps): bump goreleaser/goreleaser-action from 4 to 6
- 2024-07-06 Merge pull request #220 from gbprz/skip-helm-version-footer
- 2024-02-27 add a flag to remove the version footer
- 2024-07-06 Merge pull request #230 from ebuildy/feat_yaml_functions
- 2024-05-11 revert go.mod
- 2024-05-11 add toYaml example
- 2024-05-11 add example
- 2024-07-06 Merge pull request #176 from jasondamour/master
- 2023-03-29 feat: add toYaml and fromYaml to functions map
- 2024-07-06 Merge pull request #243 from norwoodj/dependabot/go_modules/github.com/spf13/cobra-1.8.1
- 2024-07-06 chore(deps): bump github.com/spf13/cobra from 1.7.0 to 1.8.1
- 2024-07-06 Merge pull request #244 from norwoodj/dependabot/go_modules/github.com/stretchr/testify-1.9.0
- 2024-07-06 chore(deps): bump github.com/stretchr/testify from 1.8.3 to 1.9.0
- 2024-07-06 Merge pull request #232 from lucacome/helm3
- 2024-05-30 Use helm v3
- 2024-07-06 Merge pull request #233 from lucacome/dependabot
- 2024-05-30 Add dependabot
- 2024-05-12 Merge pull request #224 from ebuildy/fix_go_mod_version
- 2024-03-28 chore: set go version to 1.22.0
- 2024-03-02 chore: renders example chart documentation using v1.13.1

1.13.1
- 2024-03-02 Release v1.13.1
- 2024-03-02 chore: updates golang.org/x/crypto to fix vulnerability
- 2024-03-02 Merge pull request #219 from mpluhar/master
- 2024-02-26 Update documentation for default template
- 2024-02-24 chore: updates chart docs with v1.13.0

1.13.0
- 2024-02-24 Release v1.13.0
- 2024-02-24 chore: adds git-cliff configuration, generates historical changelog and adds a release script
- 2024-02-24 fix: Solves #217 where helm-docs would segfault due in charts with certain comment format
- 2024-02-24 Merge pull request #212 from chenrui333/update-license-to-use-spdx-id
- 2024-02-01 license: update to use spdx id
- 2024-02-24 Merge pull request #209 from UnsolvedCypher/master
- 2024-01-03 fix: Correct the name of the GitHub repository in the README
- 2024-02-24 Merge pull request #211 from Footur/update-alpine
- 2024-01-10 Update Alpine to v3.19
- 2024-02-24 Merge pull request #213 from chenrui333/go-1.22
- 2024-02-07 feat: update to use go1.22
- 2023-12-21 chore: updates generated chart documentation for version v1.12.0

1.12.0
- 2023-12-21 Merge pull request #194 from Haepaxlog/sections
- 2023-11-29 conform to MD022
- 2023-11-08 Change example section chart to new templates
- 2023-11-08 Add some comments about parsing of key comments
- 2023-11-08 Sectioned values are default if provied
- 2023-11-08 Factor out sorting value rows and add dedicated default section
- 2023-09-06 Put the creation and sorting of sectioned Value Rows in its own functions
- 2023-09-05 Add Tests for section feature
- 2023-09-05 Give an example of how to use sections
- 2023-09-05 Add possibility to generate subsectioned Values Tables
- 2023-11-02 Merge pull request #159 from terrycain/files_asmap
- 2022-07-18 Adds AsMap to iterate over .Files.Glob
- 2023-11-02 Merge pull request #167 from RetGal/patch-1
- 2022-11-30 Fix filename
- 2023-11-02 Merge pull request #172 from gianklug/patch-1
- 2023-02-23 fix(README): change the way helm-docs is installed
- 2023-11-02 Merge pull request #184 from sblundy/notation-on-null-values
- 2023-07-12 Copy NotationType in Nil Values
- 2023-10-21 Merge pull request #191 from Labelbox/add-docker-to-pre-commit-hooks
- 2023-10-16 Update README.md with pre-commit usage examples
- 2023-10-16 Update .pre-commit-hooks.yaml
- 2023-08-17 Pin container version to tagged release v1.11.0
- 2023-08-15 Add containerized pre-commit hook

1.11.3
- 2023-10-10 Merge pull request #201 from Nepo26/200-bug-binary-artifacts-renamed-in-v1112
- 2023-10-10 hotfix: changing back artifacts name
- 2023-09-26 Merge pull request #199 from norwoodj/feat/improving-community-standards
- 2023-09-23 fix: correct contributing link referent in pr template
- 2023-09-23 fix: changed to always get the latest version on the helm docs pre commit
- 2023-09-20 Merge pull request #198 from brettmorien/master
- 2023-09-15 Bump all available dependencies to latest.

1.11.2
- 2023-09-10 Revert "fix: String additional empty line from generated README.md"
- 2023-09-10 chore: update github actions to be able to reproduce using act
- 2023-09-10 chore: adding todo to refactor main test
- 2023-09-10 fix: removing goreleaser project env var to be able to test locally
- 2023-09-10 fix: remove var env dependency by moving tests
- 2023-07-25 chore: change readme to trigger ci

1.11.1
- 2023-07-25 chore: change readme to trigger ci
- 2023-07-25 Merge pull request #181 from edmondop/issue-169
- 2023-07-25 Update pkg/helm/chart_info.go
- 2023-07-25 fix: change error to err to not conflict with builtin interface
- 2023-07-25 fix: update goreleaser and way to get env
- 2023-07-25 fix: update actions
- 2023-06-29 Fixed GoReleaser
- 2023-06-29 Fixing build
- 2023-06-29 Fixed deprecation
- 2023-06-29 Fixing GoReleaser deprecated action
- 2023-06-29 Introducing options from the CLI and unit test to confirm strict linting of documentation comments
- 2023-07-01 Merge pull request #178 from jlec/fix-177
- 2023-04-03 fix: String additional empty line from generated README.md
- 2022-06-29 chore: updates example chart READMEs for v1.11.0

1.11.0
- 2022-06-29 fix: fixes file operations to work when not running from the chart root and fixes several tests
- 2022-06-29 Merge pull request #141 from j-buczak/ignoring_values
- 2022-05-11 Improve code according to mjpitz suggestions
- 2022-05-10 Add option for ignoring values
- 2022-06-29 Merge pull request #139 from terrycain/107-files
- 2022-05-11 Added file lazy loading
- 2022-04-26 Added Helm .Files
- 2022-06-29 Merge pull request #142 from j-buczak/rename_section_to_raw
- 2022-05-11 Rename @section to @raw
- 2022-06-29 Merge pull request #145 from j-buczak/charts_to_generate_flag
- 2022-05-12 Add an option to list charts to generate
- 2022-06-29 Merge pull request #146 from j-buczak/fix_file_sorter
- 2022-05-12 fix broken file sorting
- 2022-06-29 Merge pull request #151 from armosec/master
- 2022-05-23 adding ignore-non-descriptions flag
- 2022-05-10 chore: generate READMEs for example charts with new version

1.10.0
- 2022-05-10 Merge pull request #140 from norwoodj/fix/nil-value-types
- 2022-05-10 fix: types on nil values
- 2022-05-10 Merge pull request #87 from lucernae/notation-type
- 2021-02-02 Add support for custom notation type
- 2022-04-25 Generated example charts with new version

1.9.1
- 2022-04-25 Revert "Add angle brackets around urls in requirementsTable"
- 2022-04-25 Generated example charts with new version

1.9.0
- 2022-04-25 Merge pull request #112 from armsnyder/umbrella-values
- 2022-03-22 Fix issue where an empty global object in a child chart would be listed in the root docs
- 2021-10-27 Warn about remote dependencies without erroring; Parse local file:// repositories
- 2021-10-14 Fix documented globals prefixed with the sub-chart alias
- 2021-10-12 Tolerate dependency charts without values.yaml
- 2021-10-04 New flag --document-dependency-values
- 2022-04-25 Merge pull request #136 from norwoodj/pr-124
- 2022-04-25 Updates alpine docker image to fix issue #124
- 2022-04-25 Merge pull request #134 from vladimir-babichev/no-bare-urls-in-requirements-table
- 2022-04-06 Add angle brackets around urls in requirementsTable
- 2022-04-03 Runs newest helm-docs to update docs for example charts

1.8.1
- 2022-04-03 Don't print angle brackets for URL/email if not present

1.8.0
- 2022-04-03 Merge pull request #102 from dfarrell07/no_raw_url
- 2021-06-25 Avoid raw URLs in maintainer tables
- 2022-04-03 Updgrades sprig to v3
- 2022-04-03 Merge pull request #121 from dirien/badge-style
- 2022-04-03 Merge branch 'master' into badge-style
- 2022-04-03 Merge pull request #130 from maybolt/values-file-option-argument
- 2022-03-21 Add option for a values file named other than `values.yaml`.
- 2022-04-03 Merge pull request #132 from norwoodj/jnorwood/fix/131
- 2022-04-01 Updgrades golang/x/sys to fix #131

1.7.0
- 2022-01-19 ci: fix environment variable reference
- 2022-01-19 ci: add job for importing GPG private key (#122)
- 2022-01-19 fix: updates signing key so release builds can work again
- 2022-01-18 Merge pull request #99 from bmcustodio/bmcustodio-fix-comments
- 2021-06-08 Ignore comment nodes not containing '# --'.
- 2021-05-20 Consider only the last group of comments starting with '# --'.
- 2022-01-16 feat: make the badge style from shields.io configurable
- 2021-12-10 remove deprecated goreleaser use_buildx option
- 2021-11-29 Merge pull request #113 from armsnyder/benchmark
- 2021-10-05 Add a benchmark
- 2021-10-20 Merge pull request #114 from jrottenberg/pre-commit.com-559
- 2021-10-10 No jump required
- 2021-10-10 wip
- 2021-09-06 Merge pull request #104 from kd7lxl/fix-empty
- 2021-07-30 fix type definition when description is empty
- 2021-05-18 fix: makes description appear even if unrelated comment appears before description comment fixes #92
- 2021-05-18 Merge pull request #95 from sc250024/chore-GoReleaserM1
- 2021-04-26 ci(github): updating for docker multi-arch builds
- 2021-04-26 ci(goreleaser): updating for docker multi-arch builds
- 2021-04-26 ci(github): updating actions for apple silicon builds
- 2021-04-26 ci(goreleaser): updating linux package naming for consistency
- 2021-04-26 Merge pull request #93 from goostleek/patch-1
- 2021-04-04 Extend README Installation section with scoop alternative for Windows
- 2021-02-10 Merge pull request #79 from sagikazarmark/fix-77
- 2021-01-13 Do not stop loading templates when a file cannot be found
- 2021-01-13 Add breaking test for default template loading
- 2021-01-13 Add test for template loading
- 2021-02-10 Merge pull request #85 from jsoref/spelling
- 2021-02-01 spelling: search
- 2021-02-01 spelling: equivalently
- 2021-02-10 Merge pull request #83 from stretched/typo-fix
- 2021-01-19 Fix typo: serach -> search
- 2021-01-14 Merge pull request #80 from sagikazarmark/fix-typo
- 2021-01-13 Fix typo in readme
- 2021-01-14 Merge pull request #81 from alexrashed/bugfix/default-for-nil
- 2021-01-14 fix: fixes @default for variables without value
- 2021-01-12 chore: runs newest helm-docs version against example charts

1.5.0
- 2021-01-12 fix: fixes small formatting issue
- 2021-01-12 fix: fixes broken comment parsing on values files with dos line endings
- 2021-01-12 Merge pull request #69 from horacimacias/master
- 2020-10-30 Parse requirement's alias and display "alias(name)" if 'alias' was defined
- 2021-01-12 Merge pull request #75 from sagikazarmark/add-badges-section
- 2021-01-07 Update full-template example
- 2021-01-07 Update the default template to use badgesSection
- 2021-01-07 Add badgesSection template
- 2021-01-12 Merge pull request #76 from sagikazarmark/fix-set-env
- 2021-01-07 Upgrade go version
- 2021-01-07 Update setup-go action

1.4.0
- 2020-10-20 Merge pull request #66 from sc250024/issue-51
- 2020-10-10 fix(goreleaser): adding more os and arch types
- 2020-10-10 chore(goreleaser): fixing brew github deprecation
- 2020-10-06 chore: updates charts with new helm-docs version
- 2020-10-06 Merge pull request #65 from norwoodj/rb/use-file-order

1.3.0
- 2020-10-06 fix: fixes tests by calling correct method
- 2020-10-06 fix: small issue in sorting by file location
- 2020-10-02  change flag
- 2020-10-01 feat: add support for sorting based on presence in file
- 2020-10-06 fix: link to hashbash
- 2020-10-06 Merge pull request #64 from norwoodj/jnorwood/comprehensive-example
- 2020-10-06 feat: adds a new chart with good examples, cleans up README a bit more, and shhhh... fixes a bug
- 2020-10-04 chore: updates readmes with latest version of helm-docs

1.2.1
- 2020-10-04 fix: makes it so charts with empty values files still get helm-docs-version footers
- 2020-10-04 chore: updates readmes with latest version of helm-docs

1.2.0
- 2020-10-04 chore: updates pre-commit hook version
- 2020-10-04 Merge pull request #63 from norwoodj/add-helm-docs-version-line
- 2020-10-04 feat: updates default chart to add a footer to markdown files with the helm-docs version, if set
- 2020-10-03 Merge pull request #62 from norwoodj/fix-diry
- 2020-10-03 feat: updates to add chart search path flag and to search for template files differently based on how they're presented fixes #47
- 2020-10-03 Merge pull request #54 from DirtyCajunRice/master
- 2020-09-13 allow multiple templates
- 2020-10-03 Merge pull request #61 from norwoodj/fix-dash-versions
- 2020-10-03 fix: escapes dashes in version badges so complicated versions work fixes #56

1.1.0
- 2020-10-03 Merge pull request #60 from norwoodj/new-style-comment-parity
- 2020-10-03 feat: Largely expands all features for old comments to new comments, old-style comments effectively deprecated

1.0.0
- 2020-10-02 Merge pull request #59 from norwoodj/auto-find-comments
- 2020-10-02 feat!: Adds the capability to provide comments without the path to the documented field fixes #58

0.15.0
- 2020-08-07 Merge pull request #50 from sc250024/fix-GithubRelease
- 2020-08-07 fix: GitHub token for Homebrew tap

0.14.0
- 2020-08-06 feat: updates signature public key for action
- 2020-08-06 Merge pull request #41 from ccremer/actions
- 2020-05-04 Add FPM/DEB packaging method, sign checksums
- 2020-05-04 Add Github Actions CI/CD workflows
- 2020-08-06 Merge pull request #46 from matheusfm/feature/name-template
- 2020-06-19 add 'chart.name' template description in 'Available Templates' section
- 2020-06-17 add chart name template to examples
- 2020-06-17 add template to the chart name
- 2020-08-06 Merge pull request #49 from sc250024/fix-PreCommitHookReadmeTemplate
- 2020-07-30 fix: Updating hook config for README.md.gotmpl
- 2020-08-06 Merge pull request #43 from acim/install-from-source
- 2020-05-29 add doc how to install from source
- 2020-05-22 Merge pull request #42 from Artus-LHIND/master
- 2020-05-15 Merge pull request #4 from Artus-LHIND/develop
- 2020-05-15 chore: update default template in README.md
- 2020-05-15 Merge pull request #3 from Artus-LHIND/develop
- 2020-05-15 fix: no double whitespace using the badges
- 2020-05-15 feat: add two regex for markdown linting
- 2020-05-14 Merge pull request #2 from Artus-LHIND/develop
- 2020-05-14 refactor: badges at default template
- 2020-05-14 Merge pull request #1 from Artus-LHIND/develop
- 2020-05-14 feat: add full-template example
- 2020-05-13 feat: fixed and updated READMEs
- 2020-05-13 feat: check for line templates, update README.md
- 2020-05-13 feat: deprecation, badges, styling, more options
- 2020-05-22 Merge pull request #40 from ccremer/patch-1
- 2020-05-04 Improve Docker usage documentation

0.13.0
- 2020-04-29 Merge pull request #38 from nvtkaszpir/dockerfile-workdir
- 2020-04-14 Add WORKDIR to Dockerfile to allow mounting charts, expand docs

0.12.0
- 2020-04-08 fix: fixes docker image name
- 2020-04-08 Merge pull request #36 from eddycharly/master
- 2020-04-03 attempt to add docker to goreleaser

0.11.1
- 2020-03-31 Merge pull request #34 from norwoodj/update-custom-default-rendering
- 2020-03-31 fix: Doesn't backtick-quote custom default values, back to the way originally implemented

0.11.0
- 2020-03-30 Merge pull request #33 from norwoodj/issue-28-output
- 2020-03-30 feat: adds a --output-file cli option for specifying the file to which documentation is written (fixes #28)
- 2020-03-30 Merge pull request #32 from norwoodj/issue-24-special
- 2020-03-29 fix: renders <, >, and & characters from default values correctly (fixes #24)
- 2020-03-30 Merge pull request #31 from norwoodj/improve-custom-default-rendering
- 2020-03-29 chore: slightly changes format of @default values as they're rendered in the mardown output

0.10.0
- 2020-03-29 Merge pull request #29 from eddycharly/master
- 2020-03-25 add comments continuation and default value

0.9.0
- 2020-01-02 Merge pull request #23 from skaro13/feature/helm-3-compatibility
- 2019-12-30 Updated Build Var Name
- 2019-12-30 Added Type field & dependencies compatibilty for Helm 3 v2 API

0.8.0
- 2019-08-13 Merge pull request #19 from norwoodj/minor-fixes
- 2019-08-10 fix: cleans up some code and minor fix to ignore case

0.7.0
- 2019-08-10 Merge pull request #18 from norwoodj/improve-ignoring
- 2019-08-09 feat: improves ignore feature, accepting an ignore file at the root of the repository as well as in the directory
- 2019-08-09 Merge pull request #17 from norwoodj/git-hook-changes
- 2019-08-09 feat: changes files for git hook slightly
- 2019-08-09 Merge pull request #16 from norwoodj/fix-chart-links
- 2019-08-09 feat: updates links to source in chart files
- 2019-08-09 Feature => Pre-commit hook (#14)

0.6.0
- 2019-08-07 Merge pull request #15 from norwoodj/add-ignore-file
- 2019-08-07 feat: adds support for an ignore file to exclude charts from processing
- 2019-08-05 Merge pull request #11 from norwoodj/add-goreport
- 2019-08-05 feat: adds goreport to the readme

0.5.0
- 2019-07-29 Merge pull request #10 from norwoodj/document-lists-and-objects
- 2019-07-28 feat: updates values table generation allowing for non-empty lists/maps to be documented

0.4.0
- 2019-07-18 feat: updates documentation with homebrew installation instructions

0.3.0
- 2019-07-18 Merge pull request #7 from norwoodj/homebrew-6
- 2019-07-18 feat: adds goreleaser configuration to deploy a homebrew tap

0.2.0
- 2019-07-17 Merge pull request #5 from norwoodj/issue-3
- 2019-07-17 feat: refactor to use gotemplates to render documentation, adds new example charts
- 2019-07-10 chore: removes unnecessary log message

0.1.1
- 2019-07-10 fix: updates name of version var to get a version in the CLI flag

0.1.0
- 2019-07-10 feat: adds goreleaser to create proper releases, creates packages, adds cobra/viper

19.0614
- 2019-06-14 Merge pull request #2 from RemingtonReackhof/fix-nil-values
- 2019-06-14 fix: change rendering of <nil> in markdown to render correctly
- 2019-06-14 Merge pull request #1 from RemingtonReackhof/fix-nil-values
- 2019-06-14 chore: update readme to include simple instructions on nil values
- 2019-06-14 fix: generate docs for nil values and include types for them
- 2019-06-14 style: run go fmt

19.0110
- 2019-01-10 Updates to follow better go project conventions
- 2018-09-15 Updates to allow for helm-docs to recursively search for charts within the project and generate docs for all of them
- 2018-09-01 Adds an example chart and some documentation. Also fixes a few small issues
- 2018-09-01 Adds most of a README, fixes a small issue with list parsing
- 2018-09-01 Adds code for the initial version, and a Makefile for building it
- 2018-09-01 Initial commit


