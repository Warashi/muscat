{
  description = "A flake for muscat";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = {
    nixpkgs,
    flake-utils,
    ...
  }:
    flake-utils.lib.eachDefaultSystem (system: let
      pkgs = nixpkgs.legacyPackages.${system};
    in
      with pkgs; {
        packages.default = (callPackage ./package.nix {});

        devShell = mkShell {
          GOTOOLCHAIN = "local";
          nativeBuildInputs = [
            go
            buf
          ];
        };
      });
}
