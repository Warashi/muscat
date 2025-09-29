{
  pkgs ? import <nixpkgs> { },
  useGolangDesign ? false,
  stdenv ? pkgs.stdenv,
  makeWrapper ? pkgs.makeWrapper,
  lib ? pkgs.lib,
  xorg ? pkgs.xorg,
  apple-sdk ? pkgs.apple-sdk,
}:
pkgs.buildGoLatestModule {
  pname = "muscat";
  version = "2.3.12";
  src = ./.;
  vendorHash = "sha256-iwot2qeOiWZMVAEm0ij1UE7xXSG1Q7NN8aj4cKgtOCw=";

  tags = if useGolangDesign then [ "golangdesign" ] else [ ];

  buildInputs =
    if stdenv.isDarwin then
      [ apple-sdk ]
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
