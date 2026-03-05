{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs?ref=nixos-unstable";
  };

  outputs =
    { nixpkgs, ... }:
    let
      lib = nixpkgs.lib;
      forEachSystem =
        f:
        builtins.foldl' lib.recursiveUpdate { } (
          map f [
            "aarch64-linux"
            "aarch64-darwin"
            "x86_64-darwin"
            "x86_64-linux"
          ]
        );
    in
    forEachSystem (
      system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in
      {
        packages.${system} = rec {
          stl = pkgs.callPackage ./nix/package.nix { };
          default = stl;
        };
      }
    );
}
