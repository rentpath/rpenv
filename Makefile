USER_EMAIL     := $(shell getent passwd $(LOGNAME) | cut -f5 -d:) <$(LOGNAME)@rentpath.com>

RPM_VENDOR     ?= RentPath, Inc.
RPM_PACKAGER   ?= Koji Build System
RPM_BUILDDIR   ?= $(HOME)/rpmbuild
RPM_DISTTAG    ?= rentpath
SRCDIR         := $(shell pwd)
SRPMDIR        := .
NAME           := $(shell basename $(SRCDIR))
ifneq ($(shell rpm -E '%{dist}'),%{dist})
  ORIGDIST     := $(patsubst %.$(RPM_DISTTAG), %, $(shell rpm -E '%{dist}'))
endif

RPMDEFINES     := --define "dist    $(ORIGDIST).$(RPM_DISTTAG)" \
               --define "vendor     $(RPM_VENDOR)"              \
               --define "packager   $(RPM_PACKAGER)"            \
               --define "_sourcedir $(SRCDIR)"                  \
               --define "_specdir   $(SRCDIR)"                  \
               --define "_srcrpmdir $(SRPMDIR)"                 \
               --define "_topdir    $(RPM_BUILDDIR)" 

SPECFILE   = $(shell find -maxdepth 1 -name \*.spec -exec basename {} \; )
PKGVERSION = $(shell awk '/Version:/ { print $$2 }' ${SPECFILE})
PKGNAME    = $(shell awk '/Name:/ { print $$2 }' ${SPECFILE})

# This rule must always be first; it will direct things to the generated "all"
# rule to get things going to generate binary RPMs in the default case.
all: sources

sources: 
	rpmbuild $(RPMDEFINES) -bs $(SPECFILE) 

rpm:
	rpmbuild $(RPMDEFINES) -bb $(SPECFILE)
