{
  description = "A flake for muscat development / build";

  inputs = {
    nixpkgs = {
      url = "github:nixos/nixpkgs/nixpkgs-unstable";
    };
    flake-parts = {
      url = "github:hercules-ci/flake-parts";
    };
    treefmt-nix = {
      url = "github:numtide/treefmt-nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };
    git-hooks = {
      url = "github:cachix/git-hooks.nix";
      inputs = {
        nixpkgs.follows = "nixpkgs";
      };
    };
  };

  outputs =
    { flake-parts, ... }@inputs:
    flake-parts.lib.mkFlake { inherit inputs; } {
      systems = [
        "x86_64-linux"
        "aarch64-linux"
        "x86_64-darwin"
        "aarch64-darwin"
      ];

      imports = [
        inputs.git-hooks.flakeModule
        inputs.treefmt-nix.flakeModule
      ];

      perSystem =
        { config, pkgs, ... }:
        {
          pre-commit = {
            check.enable = true;
            settings = {
              src = ./.;
              hooks = {
                actionlint.enable = true;
                treefmt.enable = true;
              };
            };
          };

          treefmt = {
            projectRootFile = "flake.nix";
            programs = {
              nixfmt = {
                enable = true;
                strict = true;
              };
              # keep-sorted start
              gofmt.enable = true;
              gofumpt.enable = true;
              goimports.enable = true;
              golines.enable = true;
              keep-sorted.enable = true;
              pinact.enable = true;
              typos.enable = true;
              # keep-sorted end
            };
          };

          packages.default = pkgs.callPackage ./. { };

          devShells.default = pkgs.mkShell {
            packages = with pkgs; [
              # keep-sorted start
              buf
              go
              gotools
              # keep-sorted end
            ];
            shellHook = config.pre-commit.installationScript;
          };
        };
    };
}
