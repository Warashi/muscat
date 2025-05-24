{
  pkgs ? import <nixpkgs> { },
}:
pkgs.buildGoLatestModule {
  name = "muscat";
  version = "2.3.1";
  src = ./.;
  vendorHash = "sha256-CpfgeQ+HC53uDWBis2muee5rRuKsR7KcV8aQtygAQEA=";
}
