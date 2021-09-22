# This is a shared Makefile for building menuet apps. 
# To use it, create a Makefile in your applications directory, set the name of the app, and include this file. 
# For example:
#
#   APP=Hello World
#   include $(GOPATH)/src/github.com/caseymrm/menuet/menuet.mk
#
# Optional features:
# 
# To set your CFBundleIdentifier, set IDENTIFIER:
#   IDENTIFIER=whyawake.caseymrm.github.com
#
# To monitor other directories for source changes, set LIBDIRS:
#   LIBDIR= ../go-pmset ../go-caffeinate
#
# To sign your app, set IDENTITY:
#   IDENTITY=Developer ID Application: Hello World LLC (AP2AFA8XAX)
#
# To release your app, put a Github access token in GITHUB_ACCESS_TOKEN (https://github.com/settings/tokens/new).
# You may wish to use an environment variable to avoid checking it in:
#   REPO=caseymrm/menuet
#   export GITHUB_ACCESS_TOKEN="asdfasdfasdf..."

ifndef APP
  $(error APP variable must be defined, e.g. APP=Hello World)
endif

space :=
space +=
ESCAPED_APP = $(subst $(space),\$(space),$(APP))
EXECUTABLE := $(shell echo $(subst $(space),,$(APP)) | tr '[:upper:]' '[:lower:]')
BINARY = $(ESCAPED_APP).app/Contents/MacOS/$(EXECUTABLE)
PLIST = $(ESCAPED_APP).app/Contents/Info.plist

run: $(BINARY) $(PLIST)
	./$(BINARY)

SOURCEDIRS = $(abspath $(dir $(MAKEFILE_LIST)))
SOURCES := $(shell find $(SOURCEDIRS) $(LIBDIRS) -name '*.go' -o -name '*.m' -o -name '*.h' -o -name '*.c' -o -name '*.mk' -o -name Makefile)

$(BINARY): $(SOURCES)
	go build -o $(BINARY)

ZIPFILE = $(ESCAPED_APP).zip

$(ZIPFILE): sign $(BINARY) $(PLIST)
	zip -r $(ZIPFILE) $(ESCAPED_APP).app

clean:
	rm -f $(BINARY) $(PLIST) $(ZIPFILE)

.PHONY: zip
zip: $(ZIPFILE)

.PHONY: releases
releases:
	@curl -s -H "Authorization: token $(GITHUB_ACCESS_TOKEN)" https://api.github.com/repos/$(REPO)/releases

.PHONY: release
release: $(ZIPFILE)
	curl -s -H "Authorization: token $(GITHUB_ACCESS_TOKEN)" https://api.github.com/repos/$(REPO)/releases | grep name\"
	@read -p "Version (tag_name)? " VERSION; \
	echo version $${VERSION}; \
	read -p "Name (name)? " NAME; \
	echo name $${NAME}; \
	read -p "Description (body)? " BODY; \
	echo body $${BODY}; \
	# TODO: WIP
	echo curl -H "Authorization: token $(GITHUB_ACCESS_TOKEN)" --data '{"tag_name":"$${VERSION}"}' https://api.github.com/repos/$(REPO)/releases

IDENTIFIER ?= $(EXECUTABLE).menuet.caseymrm.github.com

$(PLIST):
	echo "Generating plist..."
	@echo '<?xml version="1.0" encoding="UTF-8"?>' > $(PLIST)
	@echo '<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">' >> $(PLIST)
	@echo '<plist version="1.0">' >> $(PLIST)
	@echo '<dict>' >> $(PLIST)
	@echo '  <key>CFBundleExecutable</key>' >> $(PLIST)
	@echo '  <string>$(EXECUTABLE)</string>' >> $(PLIST)
	@echo '  <key>CFBundleIconFile</key>' >> $(PLIST)
	@echo '  <string>icon</string>' >> $(PLIST)
	@echo '  <key>CFBundleGetInfoString</key>' >> $(PLIST)
	@echo '  <string>$(APP)</string>' >> $(PLIST)
	@echo '  <key>CFBundleIdentifier</key>' >> $(PLIST)
	@echo '  <string>$(IDENTIFIER)</string>' >> $(PLIST)
	@echo '  <key>CFBundleName</key>' >> $(PLIST)
	@echo '  <string>$(APP)</string>' >> $(PLIST)
	@echo '  <key>CFBundleShortVersionString</key>' >> $(PLIST)
	@echo '  <string>0.1</string>' >> $(PLIST)
	@echo '  <key>CFBundleInfoDictionaryVersion</key>' >> $(PLIST)
	@echo '  <string>6.0</string>' >> $(PLIST)
	@echo '  <key>CFBundlePackageType</key>' >> $(PLIST)
	@echo '  <string>APPL</string>' >> $(PLIST)
	@echo '  <key>IFMajorVersion</key>' >> $(PLIST)
	@echo '  <integer>0</integer>' >> $(PLIST)
	@echo '  <key>IFMinorVersion</key>' >> $(PLIST)
	@echo '  <integer>1</integer>' >> $(PLIST)
	@echo '  <key>NSHighResolutionCapable</key><true/>' >> $(PLIST)
	@echo '  <key>NSSupportsAutomaticGraphicsSwitching</key><true/>' >> $(PLIST)
	@echo '</dict>' >> $(PLIST)
	@echo '</plist>' >> $(PLIST)

.PHONY: sign
sign: $(BINARY) $(PLIST)
	codesign -f -s "$(IDENTITY)" $(ESCAPED_APP).app --deep
