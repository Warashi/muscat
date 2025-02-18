{
  pkgs ? (
    let
      inherit (builtins) fetchTree fromJSON readFile;
      inherit ((fromJSON (readFile ./flake.lock)).nodes) nixpkgs gomod2nix;
    in
    import (fetchTree nixpkgs.locked) {
      overlays = [
        (import "${fetchTree gomod2nix.locked}/overlay.nix")
      ];
    }
  ),
  stdenv ? pkgs.stdenv,
  lib ? pkgs.lib,
  darwin ? pkgs.darwin,
  xorg ? pkgs.xorg,
  makeWrapper ? pkgs.makeWrapper,
  buildGoApplication ? pkgs.buildGoApplication,
  useGolangDesign ? false,
  ...
}:
buildGoApplication {
  pname = "muscat";
  version = "2.3.1";
  pwd = ./.;
  src = ./.;
  modules = ./gomod2nix.toml;

  tags = if useGolangDesign then [ "golangdesign" ] else [ ];

  buildInputs =
    if stdenv.isDarwin then
      [
        darwin.apple_sdk.frameworks.Cocoa
      ]
    else if useGolangDesign then
      [
        xorg.libX11
      ]
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
