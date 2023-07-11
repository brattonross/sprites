VERSION = $(shell cat version.txt)

build:
	go build -o npm/bin/sprites -v

publish:
	# Commit and tag
	git commit -am "release: v$(VERSION)"
	git tag "v$(VERSION)"
	@test -z "`git status --porcelain`" || (echo "Working directory is not clean" && false)

	# Clean up the npm directory
	rm -rf npm && git checkout npm

	# Build and publish
	$(MAKE) build
	@echo Enter one-time password:
	@read OTP && test -n "$$OTP" && cd npm && npm publish --otp "$$OTP" --access public

	git push origin main "v$(VERSION)"
