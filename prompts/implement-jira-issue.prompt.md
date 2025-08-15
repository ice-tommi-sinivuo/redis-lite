Implement a Jira issue.

If the Jira issues key is not provided, please ask for it.

Please follow the workflow below carefully to implement the Jira issue successfully.

Always be aware of which step in the workflow we currently are in, so that you can proceed to the next step after
for example iterating on a step.

## Workflow

1. Fetch the Jira issue details using the provided key.
2. Read the issue description to understand the requirements of the task.
3. Check `README.md`, `docs/architecture.md` and `docs/ai-learnings.md` (if exists) for any relevant information for successfully completing the task
4. Check reference implementations in the codebase, if available. ALWAYS follow the existing patterns and conventions in the codebase, unless they contradict with `docs/architecture.md` or given intructions
5. If the issue and it's requirements are not clear at this point, ask for clarification
6. Outline the steps required to complete the task
7. Update the Jira issue's status to "In Progress"
8. Make sure you are on the "main" or "master" branch (whichever name is used in this repo) and pull the latest changes
9. Checkout a new branch, named after the Jira issue key
10. Implement the task
11. Make sure you have written unit tests for the implementation
12. Update (or create) `README.md` and `docs/architecture.md` if necessary to reflect the changes made
13. Ask the human developer to review the implementation. If applicable, provide the human developer with steps to test the implementation themselves.
14. Once the human developer approves the implementation, please suggest a commit message for the human developer
15. Once the human developer approves the commit message, please commit the changes
16. Push the changes to the remote repository
17. Create a pull request for the changes using GitHub CLI
18. The task is now complete. Good job! Remember to provide the link to the pull request to the human developer.
