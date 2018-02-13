Name:   rpenv
Version:  3.0.4
Release:  1%{?dist}
Summary: displays env vars set from existing environment.
Source0: rpenv.go
URL: https://github.com/rentpath/rpenv
BuildRoot:  %(mktemp -ud %{_tmppath}/%{name}-%{version}-%{release}-XXXXXX)
License:  Copyright

BuildRequires: golang >= 1.9.3

%description
displays env vars set from existing environment and loaded from config file in
specified environment (ci, qa, or prod) or executes command in the context of
the existing environment variables and ones loaded from a config file.

%prep
mkdir rpenv
cp %{SOURCE0} rpenv/rpenv.go

%build
cd rpenv
/usr/bin/go build


%install
rm -rf %{buildroot}
mkdir -p %{buildroot}%{appdir}/
%{__install} -D -m 0655 rpenv/rpenv %{buildroot}%{_bindir}/rpenv


%clean
rm -rf %{buildroot}


%files
%defattr(-,root,root,-)
%{_bindir}/%{name}


%changelog
* Mon Feb 12 2018 Alan Voss <avoss@rentpath.com> - 3.0.4
- Fixing README and various version errors in repo

* Mon Feb 12 2018 Alan Voss <avoss@rentpath.com> - 3.0.3
- remove external dependency on jimlawless/cfg

* Tue Jan 30 2018 Alan Voss <avoss@rentpath.com> - 3.0.2
- Ensures config files exist with correct keys

* Tue Jan 30 2018 Alan Voss <avoss@rentpath.com> - 3.0.1
- Fixing README and various version errors in repo

* Wed Aug 28 2018 Dan McGuire <dmcguire@rentpath.com> - 3.0.0
- adds the jimlawless/cfg vendor dependency
- moves urls out of the repo and into configs

* Wed Aug 03 2016 Brad Anderson <banderson@rentpath.com> - 2.0.1
- Updating version

* Fri Jul 29 2016 Steve Doyle <sdoyle@rentpath.com> - 1.0.1
- Fix up the spec file

* Fri Feb 06 2015 Andrew Ward <award at rentpath dot com> 0.0.1
- Initial
