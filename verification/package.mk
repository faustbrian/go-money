.PHONY: docs no-float

docs:
	./verification/check-docs.sh

no-float:
	./verification/check-no-float.sh
