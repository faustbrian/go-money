GOLIB ?= golib

.PHONY: check check-full ci cohesion inventory no-float repository-check source-docs-check

check: no-float source-docs-check
	$(GOLIB) check --local --all

check-full: source-docs-check
	$(GOLIB) check --all

ci: no-float source-docs-check
	$(GOLIB) repository check
	$(GOLIB) cohesion check
	$(GOLIB) check --local --all

cohesion:
	$(GOLIB) cohesion check

no-float:
	$(MAKE) -f verification/package.mk no-float

inventory repository-check:
	$(GOLIB) repository check

source-docs-check:
	$(GOLIB) source-docs check --json
