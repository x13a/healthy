NAME        := healthy

prefix      ?= /usr/local
exec_prefix ?= $(prefix)
bindir      ?= $(exec_prefix)/bin
srcdir      ?= ./src

targetdir   := ./target
target      := $(targetdir)/$(NAME)
bindestdir  := $(DESTDIR)$(bindir)

all: build

build:
	go build -o $(target) $(srcdir)/

installdirs:
	install -d $(bindestdir)/

install: installdirs
	install $(target) $(bindestdir)/

uninstall:
	rm -f $(bindestdir)/$(NAME)

clean:
	rm -rf $(targetdir)/
