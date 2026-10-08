# Branching strategy and practical rules
See choices and challenges 05/10/2026 for reflections and decision process.  
This file is for looking up the rules during the daily work.  
  
## Rules  

### How to decide how many reviews to request

Generally just add everybody unless it is specifically about a collab on a function or a change of code that another team member previously coded.

### Everyone can merge a pull request when the criterias are met

### How many approvals needed before merge

- Typos and UI just need one approval.  
- Everything that writes to the database needs all team members approval.  
- If two out of three team members has approved it can be merged, unless another team member has reviewed with a comment.  
- A pull request can not be merged after less than 15 minutes, unless every team member has approved.

### Latest changes from main must always be merged into feature branch before merging PR

### Main branch is protected

### Feature branches / short-lived branches are deleted on GitHub right after PR is merged

### We have WIP rule for max 2 active feature branches per team member

### We choose to merge instead of rebase
We prioritize to document the development as it happened over having a more linear history.


## Description of choice
We've "chosen" GitHub flow as a branching strategy. The development workflow we employ can be summarized by create branch -> make changes -> create a PR with the changes -> pass review and automated checks -> delete branch -> repeat. Our commit strategy is to commit after the smallest possible unit of complete work to maximize the chances of successfully rolling back changes.

