# Regarding a cyberattack resulting in data-exposure 
## Happened on 30/10/2026 - This document written on 31/10/2026
Yesterday at around 11:21, we were notified of a security incident by its perpetrator and began investigating immediately. At 11:30, we identified the misconfiguration through a missed correction suggested by an automated code review tool on the pull request which containerized the application. At 11:31, we changed the single line of configuration which caused the vulnerability and reran the application container at which point the vulnerability was corrected. 

The perpetrator gained access to our customer database through a misconfiguration of the Docker container running the application server, which exposed the database to the public on port 3306 inadvertently. We have seen a subset of the data they allege to possess and verified its authenticity. For clarity, the attacker alleges to have *all* customer records. The attacker states that they are giving us a "grace period" to fix our security hole, but not whether they are releasing your records or not. As such, we assume, and correct for, the worst.
### What you must do
As is custom, your passwords are stored in hashed form with salt, but it is only a matter of time before a rainbow table reverses some of these, especially if they're weak. As such, we recommend that you **change your password immediately**. If you reused that same password on other platforms, change it on those platforms too. You ought to assume that you will fall victim to credential stuffing attacks in the near future, should you fail to change your password.

Finally, **assume your email is exposed**, and known to other hackers too. Be on the lookout for suspicious sender domains and scam attempts.
### What we have done
- We have built a patch which invalidates all passwords created before the release of the patch and reworked the login page to reflect this change. This patch will be released before any new featues.
- We have corrected the misconfiguration in our application container.
- We have integrated automated, state of the art security scanners to continuously run security tests on our code, infrastructure and dependencies.
  - We run these tests in different segments of our pipeline to catch errors before they ever get to affect you.
- We have contacted the relevant authorities to bring the perpetrator to justice.
### What we will do
- We will continuously look for ways to integrate security into our daily work, at all steps of the development process. 