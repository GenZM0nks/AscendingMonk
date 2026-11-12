ROLLBACK_SCRIPT ?= /opt/ascendingmonk-deploy/rollback.sh

.PHONY: list rollback deploy comment

list:
	sudo -n "$(ROLLBACK_SCRIPT)" list

rollback:
	sudo -n "$(ROLLBACK_SCRIPT)" rollback

deploy:
	@test -n "$(DEPLOY_ID)" || (echo "Usage: make deploy DEPLOY_ID=<DEPLOY_ID-or-sha>" >&2; exit 2)
	sudo -n "$(ROLLBACK_SCRIPT)" deploy "$(DEPLOY_ID)"

comment:
	@test -n "$(DEPLOY_ID)" -a -n "$(TEXT)" || (echo "Usage: make comment DEPLOY_ID=<DEPLOY_ID> TEXT=\"comment\"" >&2; exit 2)
	sudo -n "$(ROLLBACK_SCRIPT)" comment add "$(DEPLOY_ID)" "$(TEXT)"
