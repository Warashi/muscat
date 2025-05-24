{
  pkgs ? import <nixpkgs> { },
  useGolangDesign ? false,
  stdenv ? pkgs.stdenv,
  makeWrapper ? pkgs.makeWrapper,
  lib ? pkgs.lib,
  xorg ? pkgs.xorg,
  darwin ? pkgs.darwin,
}:
pkgs.buildGoLatestModule {
  name = "muscat";
  version = "2.3.1";
  src = ./.;
  vendorHash = "sha256-CpfgeQ+HC53uDWBis2muee5rRuKsR7KcV8aQtygAQEA=";

  tags = if useGolangDesign then [ "golangdesign" ] else [ ];

  buildInputs =
    if stdenv.isDarwin then
      [ darwin.apple_sdk.frameworks.Cocoa ]
    else if useGolangDesign then
      [ xorg.libX11 ]
    else
      [ ];

  nativeBuildInputs = if (stdenv.isLinux && useGolangDesign) then [ makeWrapper ] else [ ];

  postFixup =
    if (stdenv.isLinux && useGolangDesign) then
      ''
        wrapProgram $out/bin/muscat \
          --prefix LD_LIBRARY_PATH : ${lib.makeLibraryPath [ xorg.libX11 ]}
      ''
    else
      "";

  meta = with lib; {
    description = "remote code development utils";
    homepage = "https://github.com/Warashi/muscat";
    license = licenses.mit;
  };
}
