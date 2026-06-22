.PHONY: install typespec openapi check

install:
	npm install

typespec:
	npm run typespec

openapi:
	npm run openapi

check:
	npm run check
