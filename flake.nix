{
  description = "Moony Server Environment";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (
      system:
        let
            pkgs = nixpkgs.legacyPackages.${system};

            moonyShell = pkgs.writeTextFile {
                name = "shell.sh";
                destination = "/shell.sh";
                text = builtins.readFile ./shell.sh;
                executable = true;
            };

        in
        {
          devShells.default = pkgs.mkShell {
            pure = true;

            buildInputs = with pkgs; [
              # go deps
              go
              gopls
              gotools
              golangci-lint
              air

              # moony
              moonyShell

              # sys utils
              bash
              coreutils
              findutils     # find command
              which         # which
              gnused        # sed
              gnugrep       # grep
              nano          # nano
              ncurses       # TUI, tput
            ];

            shellHook = ''
              tput clear

              export PS1="\[\033[01;32m\][Moony]\[\033[00m\] \w $ "
              export PATH=""

              for p in $buildInputs; do
                export PATH="$p/bin:$PATH"
              done

              export GOPATH="$PWD/.nix-go"
              export GOBIN="$GOPATH/bin"
              export PATH="$GOBIN:$PATH"

              echo ""
              echo "🌙 Moony development environment loaded"
              echo ""
              echo "Go version: $(go version)"
              echo "GOPATH set to: $GOPATH"
              echo ""

              source ${moonyShell}/shell.sh
            '';
          };
        }
    );
}
