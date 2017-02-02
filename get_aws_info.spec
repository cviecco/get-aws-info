Name:           get_aws_info
Version:	0.1.0
Release:	1%{?dist}
Summary:	Simple utility to get aws instance tags

#Group:		
License:	ASL 2.0
URL:		https://github.com/cviecco/get_aws_info/
Source0:	get_aws_info-%{version}.tar.gz

#BuildRequires:	golang
#Requires:	

#no debug package as this is go
%define debug_package %{nil}

%description
Simple encryption using clound infrastrcture


%prep
%setup -n %{name}-%{version}

%build
go build -ldflags "-X main.Version=%{version}"  -o %{name} main.go
#go build -ldflags "-X main.Version=%{version}" get_aws_info.go 


%install
#%make_install
%{__install} -Dp -m0755 get_aws_info %{buildroot}%{_sbindir}/get_aws_info


%files
#%doc
%{_sbindir}/get_aws_info


%changelog

