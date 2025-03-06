{
  description = "Moony Server Environment";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (
      system:
        let pkgs = nixpkgs.legacyPackages.${system};
        in
        {
          devShells.default = pkgs.mkShell {
            buildInputs = with pkgs; [
              go
              gopls
              gotools
              golangci-lint
              delve
              air
            ];

            shellHook = ''
              echo "🌙 Moony Server Go development environment loaded"
              echo "Go version: $(go version)"

              # Add go bin to PATH
              # export PATH=$PATH:$(go env GOPATH)/bin

              # Set any environment variables your project needs
              # export DATABASE_URL="..."
            '';
          };
        }
    );
}
