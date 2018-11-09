Name: gobetween
Version: 0.7.0
Release: 1
Summary: Modern & minimalistic load balancer for the Сloud era (forked version)

Group: Unspecified
License: MIT
URL: https://github.com/yousong/gobetween

%description
This is a forked version of gobetween from https://github.com/yyyar/gobetween

%prep

%build

%install
install -D -m 0755 %{gobetween_bin} %{buildroot}/usr/sbin/gobetween

%files
/usr/sbin/gobetween
