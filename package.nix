{
  stdenv,
  lib,
  darwin,
  xorg,
  makeWrapper,
  buildGoModule,
  useGolangDesign ? false,
  ...
}:
buildGoModule {
  pname = "muscat";
  version = "2.2.2";
  vendorHash = "sha256-NyXFyuHYoRP6pXz5cmA99IzMawduPhlBs/m4BVC2h6A=";

  src = ./.;

  tags =
    if useGolangDesign
    then ["golangdesign"]
    else [];

  buildInputs =
    if stdenv.isDarwin
    then [
      darwin.apple_sdk.frameworks.Cocoa
    ]
    else if useGolangDesign
    then [
      xorg.libX11
    ]
    else [];

  nativeBuildInputs =
    if (stdenv.isLinux && useGolangDesign)
    then [makeWrapper]
    else [];

  postFixup =
    if (stdenv.isLinux && useGolangDesign)
    then ''
      wrapProgram $out/bin/muscat \
        --prefix LD_LIBRARY_PATH : ${lib.makeLibraryPath [xorg.libX11]}
    ''
    else "";

  meta = with lib; {
    description = "remote code development utils";
    homepage = "https://github.com/Warashi/muscat";
    license = licenses.mit;
  };
}
