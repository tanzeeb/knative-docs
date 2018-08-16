# PR Bot

PR Bot uses Knative Serving, Build, and Eventing to validate PRs to a
Github repository.

## Usage

1. Install Knative Serving, Build and Eventing (w/ Stub bus, Github
   event source)
1. Create a Github repository.
1. Edit `config.yaml`:
  a. Replace GITHUB_REPOSITORY with your Github Repository
  a. Replace BUILD_TASK with the command to build your image
  a. Replace TEST_TASK with the command to run your tests
1. Create a PR to your repository

PR Bot will:

1. Update the PR with a message that it's processing the PR
1. Run the tests against the PR
1. Post the test results back to the PR
