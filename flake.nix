{
  description = "Accuknox Kubernetes-native system metrics reporter";
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-24.05";
    unstable-nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };
  outputs = { nixpkgs, flake-utils, unstable-nixpkgs, ... }@inputs:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs { inherit system; };
        unstable = import unstable-nixpkgs { inherit system; };
      in
      {
        formatter = pkgs.nixpkgs-fmt;
        devShells.default = pkgs.mkShell {
          packages = [
            pkgs.nixd
            pkgs.nixpkgs-fmt
            pkgs.go
            pkgs.go-tools
            pkgs.gopls
            pkgs.templ
            pkgs.air
            pkgs.nodejs
            pkgs.nodePackages.pnpm
            pkgs.nodePackages.vscode-langservers-extracted
            pkgs.nodePackages.typescript-language-server
            pkgs.tailwindcss-language-server
            pkgs.emmet-language-server
            pkgs.prettierd
            pkgs.kubectl
            pkgs.kind
            pkgs.kubernetes-helm
            pkgs.awscli2
            unstable.goreleaser
          ];
        };
      });
}
