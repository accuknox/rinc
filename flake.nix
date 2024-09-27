{
  description = "Accuknox Kubernetes-native system metrics reporter";
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-24.05";
    flake-utils.url = "github:numtide/flake-utils";
  };
  outputs = { nixpkgs, flake-utils, ... }@inputs:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs { inherit system; };
      in
      {
        formatter = pkgs.nixpkgs-fmt;
        devShells.default = pkgs.mkShell {
          packages = with pkgs; [
            nixd
            nixpkgs-fmt
            go
            go-tools
            gopls
            templ
            air
            nodejs
            nodePackages.pnpm
            nodePackages.vscode-langservers-extracted
            nodePackages.typescript-language-server
            tailwindcss-language-server
            emmet-language-server
            prettierd
            kubectl
            kind
            kubernetes-helm
            awscli2
          ];
        };
      });
}
