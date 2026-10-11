GOLIB ?= golib

.PHONY: check check-full ci cohesion inventory no-float repository-check

check: no-float
	$(GOLIB) check --local --all

check-full:
	$(GOLIB) check --all

ci: no-float
	$(GOLIB) repository check
	$(GOLIB) cohesion check
	$(GOLIB) check --local --all

cohesion:
	$(GOLIB) cohesion check

no-float:
	$(MAKE) -f verification/package.mk no-float

inventory repository-check:
	$(GOLIB) repository check
